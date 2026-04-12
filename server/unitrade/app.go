package unitrade

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

type App struct {
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
	storePath string
	store     *ProfileStore
}

type commonFlags struct {
	host    string
	profile string
	json    bool
	yes     bool
}

type pageFlags struct {
	page     int
	pageSize int
}

type httpHeaderAdapter interface {
	Get(string) string
}

func NewApp(stdin io.Reader, stdout, stderr io.Writer) (*App, error) {
	store, path, err := loadProfileStore()
	if err != nil {
		return nil, err
	}
	return &App{
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
		storePath: path,
		store:     store,
	}, nil
}

func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.printRootUsage()
		return nil
	}

	switch args[0] {
	case "login":
		return a.runLogin(args[1:])
	case "logout":
		return a.runLogout(args[1:])
	case "whoami":
		return a.runWhoami(args[1:])
	case "auth":
		return a.runAuth(args[1:])
	case "announcement":
		return a.runAnnouncement(args[1:])
	case "product":
		return a.runProduct(args[1:])
	case "report":
		return a.runReport(args[1:])
	case "category":
		return a.runCategory(args[1:])
	case "user":
		return a.runUser(args[1:])
	case "staff":
		return a.runStaff(args[1:])
	case "order":
		return a.runOrder(args[1:])
	case "help", "-h", "--help":
		a.printRootUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *App) printRootUsage() {
	fmt.Fprintln(a.stdout, "Unitrade campus CLI")
	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Usage:")
	fmt.Fprintln(a.stdout, "  unitrade <command> <action> [flags]")
	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Commands:")
	fmt.Fprintln(a.stdout, "  login          Save a PAT-backed session")
	fmt.Fprintln(a.stdout, "  logout         Remove a saved profile")
	fmt.Fprintln(a.stdout, "  whoami         Show the active session")
	fmt.Fprintln(a.stdout, "  auth           Campus authentication workflows")
	fmt.Fprintln(a.stdout, "  announcement   Campus announcement workflows")
	fmt.Fprintln(a.stdout, "  product        Campus product workflows")
	fmt.Fprintln(a.stdout, "  report         Campus report workflows")
	fmt.Fprintln(a.stdout, "  category       Campus category workflows")
	fmt.Fprintln(a.stdout, "  user           Campus user workflows")
	fmt.Fprintln(a.stdout, "  staff          Campus staff workflows")
	fmt.Fprintln(a.stdout, "  order          Campus order queries")
}

func addCommonFlags(fs *flag.FlagSet, flags *commonFlags) {
	fs.StringVar(&flags.host, "host", "", "API host, for example http://127.0.0.1:8888")
	fs.StringVar(&flags.profile, "profile", "", "Profile name, defaults to the active profile")
	fs.BoolVar(&flags.json, "json", false, "Print JSON output")
	fs.BoolVar(&flags.yes, "yes", false, "Skip confirmation prompts")
}

func addPageFlags(fs *flag.FlagSet, flags *pageFlags) {
	fs.IntVar(&flags.page, "page", 1, "Page number")
	fs.IntVar(&flags.pageSize, "page-size", 10, "Page size")
}

func parseDateTime(raw string) (time.Time, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return time.Time{}, nil
	}
	layouts := []string{
		time.DateTime,
		time.RFC3339,
		"2006-01-02",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime %q, expected YYYY-MM-DD or YYYY-MM-DDTHH:mm:ss", raw)
}

func (a *App) writeProfile(profile Profile) error {
	a.store.setProfile(profile)
	return saveProfileStore(a.storePath, a.store)
}

func (a *App) getProfile(name string) (Profile, error) {
	profileName := a.store.activeProfileName(name)
	profile, ok := a.store.getProfile(profileName)
	if !ok {
		return Profile{}, fmt.Errorf("profile %q not found, run `unitrade login --token ...` first", profileName)
	}
	profile.Name = profileName
	return profile, nil
}

func (a *App) saveUpdatedProfile(profile Profile, headers httpHeaderAdapter) error {
	newToken := headers.Get("new-token")
	if newToken == "" {
		return nil
	}
	profile.Token = newToken
	if expRaw := headers.Get("new-expires-at"); expRaw != "" {
		if unixSeconds, err := strconv.ParseInt(expRaw, 10, 64); err == nil {
			profile.ExpiresAt = time.Unix(unixSeconds, 0).Local()
		}
	}
	return a.writeProfile(profile)
}

func (a *App) newAuthedClient(flags commonFlags) (*apiClient, Profile, error) {
	profile, err := a.getProfile(flags.profile)
	if err != nil {
		return nil, Profile{}, err
	}
	host := profile.Host
	if strings.TrimSpace(flags.host) != "" {
		host = flags.host
	}
	client, err := newAPIClient(host, profile.Token, profile.UserID)
	if err != nil {
		return nil, Profile{}, err
	}
	return client, profile, nil
}

func (a *App) printActionResult(asJSON bool, message string, extra map[string]any) error {
	if asJSON {
		payload := map[string]any{"ok": true, "message": message}
		for key, value := range extra {
			payload[key] = value
		}
		return printJSON(a.stdout, payload)
	}
	_, err := fmt.Fprintln(a.stdout, message)
	return err
}

func (a *App) confirm(prompt string, autoYes bool) error {
	if autoYes {
		return nil
	}
	fmt.Fprintf(a.stdout, "%s [y/N]: ", prompt)
	reader := bufio.NewReader(a.stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	if answer == "y" || answer == "yes" {
		return nil
	}
	return errors.New("operation cancelled")
}

func requireID(id uint, label string) error {
	if id == 0 {
		return fmt.Errorf("%s is required", label)
	}
	return nil
}

func parseEnumValue[T ~int](raw string, mapping map[string]T, label string) (T, error) {
	var zero T
	value := strings.TrimSpace(raw)
	if value == "" {
		return zero, fmt.Errorf("%s is required", label)
	}
	mapped, ok := mapping[value]
	if !ok {
		return zero, fmt.Errorf("invalid %s %q", label, raw)
	}
	return mapped, nil
}

func readContentInput(inline, filePath string, fromSTDIN bool, stdin io.Reader) (string, error) {
	count := 0
	if strings.TrimSpace(inline) != "" {
		count++
	}
	if strings.TrimSpace(filePath) != "" {
		count++
	}
	if fromSTDIN {
		count++
	}
	if count != 1 {
		return "", errors.New("choose exactly one of --content, --content-file, or --content-stdin")
	}

	if strings.TrimSpace(inline) != "" {
		return inline, nil
	}
	if strings.TrimSpace(filePath) != "" {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func maybeAddRange(values url.Values, key string, from, to time.Time) {
	if !from.IsZero() {
		values.Add(key, from.Format(time.RFC3339))
	}
	if !to.IsZero() {
		values.Add(key, to.Format(time.RFC3339))
	}
}

func chooseNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func formatAnyTime(value any) string {
	if t, ok := value.(time.Time); ok {
		return formatTimeValue(t)
	}
	return fmt.Sprintf("%v", value)
}

func marshalForJSON(value any) any {
	raw, err := json.Marshal(value)
	if err != nil {
		return value
	}
	var out any
	if err := json.Unmarshal(raw, &out); err != nil {
		return value
	}
	return out
}

func (a *App) runLogin(args []string) error {
	fs := flag.NewFlagSet("login", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var token string
	fs.StringVar(&token, "token", "", "PAT or API token")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(token) == "" {
		return errors.New("--token is required")
	}

	client, err := newAPIClient(flags.host, token, 0)
	if err != nil {
		return err
	}

	var userResp struct {
		UserInfo system.SysUser `json:"userInfo"`
	}
	headers, err := client.doJSON(http.MethodGet, "user/getUserInfo", nil, nil, &userResp)
	if err != nil {
		return err
	}

	claims, err := parseTokenClaims(token)
	if err != nil {
		return fmt.Errorf("token is not a valid JWT: %w", err)
	}

	profile := Profile{
		Name:        a.store.activeProfileName(flags.profile),
		Host:        chooseNonEmpty(flags.host, defaultHost),
		Token:       token,
		UserID:      userResp.UserInfo.ID,
		AuthorityID: userResp.UserInfo.AuthorityId,
		Username:    claims.Username,
		NickName:    claims.NickName,
	}
	if claims.ExpiresAt != nil {
		profile.ExpiresAt = claims.ExpiresAt.Time.Local()
	}
	if profile.Username == "" {
		profile.Username = userResp.UserInfo.Username
	}
	if profile.NickName == "" {
		profile.NickName = userResp.UserInfo.NickName
	}
	if headerToken := headers.Get("new-token"); headerToken != "" {
		profile.Token = headerToken
		if parsed, parseErr := parseTokenClaims(headerToken); parseErr == nil && parsed.ExpiresAt != nil {
			profile.ExpiresAt = parsed.ExpiresAt.Time.Local()
		}
	}

	if err := a.writeProfile(profile); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("saved profile %q for %s", profile.Name, profile.Username), map[string]any{
		"profile":     profile.Name,
		"host":        profile.Host,
		"userId":      profile.UserID,
		"authorityId": profile.AuthorityID,
		"expiresAt":   profile.ExpiresAt,
	})
}

func (a *App) runLogout(args []string) error {
	fs := flag.NewFlagSet("logout", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	if err := fs.Parse(args); err != nil {
		return err
	}
	name := a.store.activeProfileName(flags.profile)
	if _, ok := a.store.getProfile(name); !ok {
		return fmt.Errorf("profile %q not found", name)
	}
	a.store.deleteProfile(name)
	if err := saveProfileStore(a.storePath, a.store); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("removed profile %q", name), map[string]any{"profile": name})
}

func (a *App) runWhoami(args []string) error {
	fs := flag.NewFlagSet("whoami", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	var userResp struct {
		UserInfo system.SysUser `json:"userInfo"`
	}
	headers, err := client.doJSON(http.MethodGet, "user/getUserInfo", nil, nil, &userResp)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	payload := map[string]any{
		"profile":     profile.Name,
		"host":        chooseNonEmpty(flags.host, profile.Host),
		"userId":      userResp.UserInfo.ID,
		"username":    userResp.UserInfo.Username,
		"nickName":    userResp.UserInfo.NickName,
		"authorityId": userResp.UserInfo.AuthorityId,
		"expiresAt":   profile.ExpiresAt,
	}
	if flags.json {
		return printJSON(a.stdout, payload)
	}
	printTable(a.stdout, []string{"PROFILE", "HOST", "USER_ID", "USERNAME", "NICKNAME", "AUTHORITY_ID", "EXPIRES_AT"}, [][]string{{
		fmt.Sprintf("%v", payload["profile"]),
		fmt.Sprintf("%v", payload["host"]),
		fmt.Sprintf("%v", payload["userId"]),
		fmt.Sprintf("%v", payload["username"]),
		fmt.Sprintf("%v", payload["nickName"]),
		fmt.Sprintf("%v", payload["authorityId"]),
		formatAnyTime(payload["expiresAt"]),
	}})
	return nil
}
