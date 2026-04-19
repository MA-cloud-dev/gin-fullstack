package unitrade

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/campus"
)

var announcementStatusMap = map[string]int{"offline": 0, "online": 1}

type pageResponse[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

func (a *App) runAuth(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("auth")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runAuthList(args[1:])
	case "show":
		return a.runAuthShow(args[1:])
	case "approve":
		return a.runAuthApprove(args[1:])
	case "reject":
		return a.runAuthReject(args[1:])
	case "revoke":
		return a.runAuthRevoke(args[1:])
	default:
		return fmt.Errorf("unknown auth action %q", args[0])
	}
}

func (a *App) runAuthList(args []string) error {
	fs := flag.NewFlagSet("auth list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var studentID, realName, college, reviewed, reviewStatus string
	var createdFromRaw, createdToRaw, reviewedFromRaw, reviewedToRaw string
	fs.StringVar(&studentID, "student-id", "", "Student ID")
	fs.StringVar(&realName, "real-name", "", "Real name")
	fs.StringVar(&college, "college", "", "College")
	fs.StringVar(&reviewed, "reviewed", "", "pending or reviewed")
	fs.StringVar(&reviewStatus, "review-status", "", "processing, approved or rejected")
	fs.StringVar(&createdFromRaw, "created-from", "", "Created from")
	fs.StringVar(&createdToRaw, "created-to", "", "Created to")
	fs.StringVar(&reviewedFromRaw, "reviewed-from", "", "Reviewed from")
	fs.StringVar(&reviewedToRaw, "reviewed-to", "", "Reviewed to")
	if err := fs.Parse(args); err != nil {
		return err
	}

	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	values := url.Values{}
	values.Set("page", strconv.Itoa(page.page))
	values.Set("pageSize", strconv.Itoa(page.pageSize))
	if studentID != "" {
		values.Set("studentId", studentID)
	}
	if realName != "" {
		values.Set("realName", realName)
	}
	if college != "" {
		values.Set("college", college)
	}
	if reviewStatus != "" {
		values.Set("reviewStatus", reviewStatus)
	} else if reviewed != "" {
		switch reviewed {
		case "pending":
			values.Set("reviewStatus", "processing")
		case "reviewed":
			// Backward compatibility: fetch all reviewed records without forcing a single final state.
		default:
			return fmt.Errorf("invalid reviewed value %q", reviewed)
		}
	}
	createdFrom, createdTo, err := parseDateRange(createdFromRaw, createdToRaw)
	if err != nil {
		return err
	}
	maybeAddRange(values, "createdAtRange[]", createdFrom, createdTo)
	reviewedFrom, reviewedTo, err := parseDateRange(reviewedFromRaw, reviewedToRaw)
	if err != nil {
		return err
	}
	maybeAddRange(values, "reviewedAtRange[]", reviewedFrom, reviewedTo)

	var result pageResponse[campus.CampusAuth]
	headers, err := client.doJSON(http.MethodGet, "campusAuth/getCampusAuthList", values, nil, &result)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(result))
	}
	rows := make([][]string, 0, len(result.List))
	for _, item := range result.List {
		status := strings.TrimSpace(item.ReviewStatus)
		if status == "" {
			if item.ReviewedAt != nil {
				status = "approved"
			} else {
				status = "processing"
			}
		}
		if reviewed == "reviewed" && status == "processing" {
			continue
		}
		rows = append(rows, []string{
			strconv.FormatUint(uint64(item.ID), 10),
			strconv.FormatUint(uint64(item.UserID), 10),
			item.StudentID,
			item.RealName,
			item.College,
			status,
			stringValue(item.ReviewRemark),
			stringValue(item.ReviewedByName),
			formatTimePtr(item.ReviewedAt),
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "USER_ID", "STUDENT_ID", "REAL_NAME", "COLLEGE", "STATUS", "REMARK", "REVIEWED_BY", "REVIEWED_AT", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runAuthShow(args []string) error {
	fs := flag.NewFlagSet("auth show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Auth ID")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	values := url.Values{"id": []string{strconv.FormatUint(uint64(id), 10)}}
	var item campus.CampusAuth
	headers, err := client.doJSON(http.MethodGet, "campusAuth/findCampusAuth", values, nil, &item)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(item))
	}
	rows := [][]string{
		{"ID", strconv.FormatUint(uint64(item.ID), 10)},
		{"User ID", strconv.FormatUint(uint64(item.UserID), 10)},
		{"Student ID", item.StudentID},
		{"Real Name", item.RealName},
		{"College", item.College},
		{"Review Remark", stringValue(item.ReviewRemark)},
		{"Reviewed By", stringValue(item.ReviewedByName)},
		{"Reviewed At", formatTimePtr(item.ReviewedAt)},
		{"Created At", formatTimeValue(item.CreatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runAuthApprove(args []string) error {
	fs := flag.NewFlagSet("auth approve", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var reason string
	fs.UintVar(&id, "id", 0, "Auth ID")
	fs.StringVar(&reason, "reason", "", "Audit reason")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("--reason is required")
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusAuth/reviewCampusAuth", nil, map[string]any{
		"id":           id,
		"reviewRemark": reason,
		"auditReason":  reason,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("approved auth record %d", id), map[string]any{"id": id})
}

func (a *App) runAuthReject(args []string) error {
	fs := flag.NewFlagSet("auth reject", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var reason string
	fs.UintVar(&id, "id", 0, "Auth ID")
	fs.StringVar(&reason, "reason", "", "Audit reason")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("--reason is required")
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusAuth/rejectCampusAuth", nil, map[string]any{
		"id":           id,
		"reviewRemark": reason,
		"auditReason":  reason,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("rejected auth record %d", id), map[string]any{"id": id})
}

func (a *App) runAuthRevoke(args []string) error {
	fs := flag.NewFlagSet("auth revoke", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var reason string
	fs.UintVar(&id, "id", 0, "Auth ID")
	fs.StringVar(&reason, "reason", "", "Audit reason")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("--reason is required")
	}
	if err := a.confirm(fmt.Sprintf("Revoke auth review %d?", id), flags.yes); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusAuth/revokeCampusAuth", nil, map[string]any{
		"id":          id,
		"auditReason": reason,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("revoked auth review %d", id), map[string]any{"id": id})
}

func (a *App) runAnnouncement(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("announcement")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runAnnouncementList(args[1:])
	case "show":
		return a.runAnnouncementShow(args[1:])
	case "create":
		return a.runAnnouncementCreate(args[1:])
	case "update":
		return a.runAnnouncementUpdate(args[1:])
	case "publish":
		return a.runAnnouncementStatus(args[1:], true)
	case "unpublish":
		return a.runAnnouncementStatus(args[1:], false)
	case "delete":
		return a.runAnnouncementDelete(args[1:])
	default:
		return fmt.Errorf("unknown announcement action %q", args[0])
	}
}

func (a *App) runAnnouncementList(args []string) error {
	fs := flag.NewFlagSet("announcement list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var title, publisherID, publisherKeyword, statusRaw, createdFromRaw, createdToRaw string
	fs.StringVar(&title, "title", "", "Title")
	fs.StringVar(&publisherID, "publisher-id", "", "Publisher user ID")
	fs.StringVar(&publisherKeyword, "publisher-keyword", "", "Publisher keyword")
	fs.StringVar(&statusRaw, "status", "", "online or offline")
	fs.StringVar(&createdFromRaw, "created-from", "", "Created from")
	fs.StringVar(&createdToRaw, "created-to", "", "Created to")
	if err := fs.Parse(args); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	values := url.Values{
		"page":     []string{strconv.Itoa(page.page)},
		"pageSize": []string{strconv.Itoa(page.pageSize)},
	}
	if title != "" {
		values.Set("title", title)
	}
	if publisherID != "" {
		values.Set("publisherId", publisherID)
	}
	if publisherKeyword != "" {
		values.Set("publisherKeyword", publisherKeyword)
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, announcementStatusMap, "status")
		if err != nil {
			return err
		}
		values.Set("status", strconv.Itoa(status))
	}
	createdFrom, createdTo, err := parseDateRange(createdFromRaw, createdToRaw)
	if err != nil {
		return err
	}
	maybeAddRange(values, "createdAtRange[]", createdFrom, createdTo)
	var result pageResponse[campus.CampusAnnouncement]
	headers, err := client.doJSON(http.MethodGet, "campusAnnouncement/getCampusAnnouncementList", values, nil, &result)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(result))
	}
	rows := make([][]string, 0, len(result.List))
	for _, item := range result.List {
		rows = append(rows, []string{
			strconv.FormatUint(uint64(item.ID), 10),
			item.Title,
			strconv.FormatUint(uint64(item.PublisherID), 10),
			stringValue(item.PublisherNickname),
			stringValue(item.PublisherPhone),
			item.StatusText,
			formatTimeValue(item.CreatedAt),
			formatTimeValue(item.UpdatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "TITLE", "PUBLISHER_ID", "PUBLISHER", "PHONE", "STATUS", "CREATED_AT", "UPDATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runAnnouncementShow(args []string) error {
	fs := flag.NewFlagSet("announcement show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Announcement ID")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	values := url.Values{"id": []string{strconv.FormatUint(uint64(id), 10)}}
	var item campus.CampusAnnouncement
	headers, err := client.doJSON(http.MethodGet, "campusAnnouncement/findCampusAnnouncement", values, nil, &item)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(item))
	}
	rows := [][]string{
		{"ID", strconv.FormatUint(uint64(item.ID), 10)},
		{"Title", item.Title},
		{"Publisher ID", strconv.FormatUint(uint64(item.PublisherID), 10)},
		{"Publisher", stringValue(item.PublisherNickname)},
		{"Publisher Phone", stringValue(item.PublisherPhone)},
		{"Status", item.StatusText},
		{"Created At", formatTimeValue(item.CreatedAt)},
		{"Updated At", formatTimeValue(item.UpdatedAt)},
		{"Content", item.Content},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runAnnouncementCreate(args []string) error {
	fs := flag.NewFlagSet("announcement create", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var title, statusRaw, contentInline, contentFile string
	var publisherID uint
	var contentSTDIN bool
	fs.StringVar(&title, "title", "", "Title")
	fs.UintVar(&publisherID, "publisher-id", 0, "Publisher user ID")
	fs.StringVar(&statusRaw, "status", "online", "online or offline")
	fs.StringVar(&contentInline, "content", "", "Inline content")
	fs.StringVar(&contentFile, "content-file", "", "Content file path")
	fs.BoolVar(&contentSTDIN, "content-stdin", false, "Read content from stdin")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(title) == "" {
		return errors.New("--title is required")
	}
	if publisherID == 0 {
		return errors.New("--publisher-id is required")
	}
	status, err := parseEnumValue(statusRaw, announcementStatusMap, "status")
	if err != nil {
		return err
	}
	content, err := readContentInput(contentInline, contentFile, contentSTDIN, a.stdin)
	if err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusAnnouncement/createCampusAnnouncement", nil, map[string]any{
		"title":       title,
		"content":     content,
		"publisherId": publisherID,
		"status":      status,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("created announcement %q", title), map[string]any{"title": title, "publisherId": publisherID})
}

func (a *App) runAnnouncementUpdate(args []string) error {
	fs := flag.NewFlagSet("announcement update", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id, publisherID uint
	var title, statusRaw, contentInline, contentFile string
	var contentSTDIN bool
	fs.UintVar(&id, "id", 0, "Announcement ID")
	fs.StringVar(&title, "title", "", "Title")
	fs.UintVar(&publisherID, "publisher-id", 0, "Publisher user ID")
	fs.StringVar(&statusRaw, "status", "", "online or offline")
	fs.StringVar(&contentInline, "content", "", "Inline content")
	fs.StringVar(&contentFile, "content-file", "", "Content file path")
	fs.BoolVar(&contentSTDIN, "content-stdin", false, "Read content from stdin")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	current, headerToken, err := a.findAnnouncement(client, id)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headerToken); err != nil {
		return err
	}
	if title != "" {
		current.Title = title
	}
	if publisherID != 0 {
		current.PublisherID = publisherID
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, announcementStatusMap, "status")
		if err != nil {
			return err
		}
		current.Status = status
	}
	if contentInline != "" || contentFile != "" || contentSTDIN {
		content, err := readContentInput(contentInline, contentFile, contentSTDIN, a.stdin)
		if err != nil {
			return err
		}
		current.Content = content
	}

	headers, err := client.doJSON(http.MethodPut, "campusAnnouncement/updateCampusAnnouncement", nil, map[string]any{
		"id":          current.ID,
		"title":       current.Title,
		"content":     current.Content,
		"publisherId": current.PublisherID,
		"status":      current.Status,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("updated announcement %d", id), map[string]any{"id": id})
}

func (a *App) runAnnouncementStatus(args []string, online bool) error {
	action := "publish"
	status := 1
	if !online {
		action = "unpublish"
		status = 0
	}
	fs := flag.NewFlagSet("announcement "+action, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Announcement ID")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusAnnouncement/updateCampusAnnouncementStatus", nil, map[string]any{
		"id":     id,
		"status": status,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("%sed announcement %d", action, id), map[string]any{"id": id, "status": status})
}

func (a *App) runAnnouncementDelete(args []string) error {
	fs := flag.NewFlagSet("announcement delete", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Announcement ID")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	if err := a.confirm(fmt.Sprintf("Delete announcement %d?", id), flags.yes); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodDelete, "campusAnnouncement/deleteCampusAnnouncement", nil, map[string]any{"id": id}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("deleted announcement %d", id), map[string]any{"id": id})
}

func (a *App) findAnnouncement(client *apiClient, id uint) (campus.CampusAnnouncement, http.Header, error) {
	values := url.Values{"id": []string{strconv.FormatUint(uint64(id), 10)}}
	var item campus.CampusAnnouncement
	headers, err := client.doJSON(http.MethodGet, "campusAnnouncement/findCampusAnnouncement", values, nil, &item)
	return item, headers, err
}
