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

var (
	productStatusMap    = map[string]int{"on-sale": 0, "trading": 2, "off-shelf": 3}
	tradeModeMap        = map[string]int{"pickup": 0, "delivery": 1}
	reportReasonMap     = map[string]int{"fake": 1, "violation": 2, "fraud": 3}
	reportStatusMap     = map[string]int{"pending": 0, "handled": 1}
	reportTargetTypeMap = map[string]int{"product": 1}
	statusMap           = map[string]int{"enabled": 0, "disabled": 1}
)

func (a *App) runProduct(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("product")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runProductList(args[1:])
	case "show":
		return a.runProductShow(args[1:])
	case "set-status":
		return a.runProductSetStatus(args[1:])
	default:
		return fmt.Errorf("unknown product action %q", args[0])
	}
}

func (a *App) runProductList(args []string) error {
	fs := flag.NewFlagSet("product list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var id, title, userID, categoryID, statusRaw, tradeModeRaw, createdFromRaw, createdToRaw string
	fs.StringVar(&id, "id", "", "Product ID")
	fs.StringVar(&title, "title", "", "Title")
	fs.StringVar(&userID, "user-id", "", "Publisher user ID")
	fs.StringVar(&categoryID, "category-id", "", "Category ID")
	fs.StringVar(&statusRaw, "status", "", "on-sale, trading, or off-shelf")
	fs.StringVar(&tradeModeRaw, "trade-mode", "", "pickup or delivery")
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
	for key, value := range map[string]string{"id": id, "title": title, "userId": userID, "categoryId": categoryID} {
		if value != "" {
			values.Set(key, value)
		}
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, productStatusMap, "status")
		if err != nil {
			return err
		}
		values.Set("status", strconv.Itoa(status))
	}
	if tradeModeRaw != "" {
		mode, err := parseEnumValue(tradeModeRaw, tradeModeMap, "trade-mode")
		if err != nil {
			return err
		}
		values.Set("tradeMode", strconv.Itoa(mode))
	}
	createdFrom, createdTo, err := parseDateRange(createdFromRaw, createdToRaw)
	if err != nil {
		return err
	}
	maybeAddRange(values, "createdAtRange[]", createdFrom, createdTo)
	var result pageResponse[campus.CampusProduct]
	headers, err := client.doJSON(http.MethodGet, "campusProduct/getCampusProductList", values, nil, &result)
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
			strconv.FormatUint(uint64(item.UserID), 10),
			stringValue(item.CategoryName),
			item.StatusText,
			item.TradeModeText,
			fmt.Sprintf("%.2f", item.Price),
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "TITLE", "USER_ID", "CATEGORY", "STATUS", "TRADE_MODE", "PRICE", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runProductShow(args []string) error {
	fs := flag.NewFlagSet("product show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Product ID")
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
	var item campus.CampusProduct
	headers, err := client.doJSON(http.MethodGet, "campusProduct/findCampusProduct", values, nil, &item)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(item))
	}
	imageURLs := make([]string, 0, len(item.Images))
	for _, image := range item.Images {
		imageURLs = append(imageURLs, image.ImageURL)
	}
	rows := [][]string{
		{"ID", strconv.FormatUint(uint64(item.ID), 10)},
		{"Title", item.Title},
		{"Category", stringValue(item.CategoryName)},
		{"Publisher", stringValue(item.PublisherNickname)},
		{"Publisher Phone", stringValue(item.PublisherPhone)},
		{"Status", item.StatusText},
		{"Trade Mode", item.TradeModeText},
		{"Price", fmt.Sprintf("%.2f", item.Price)},
		{"Original Price", floatPtrString(item.OriginalPrice)},
		{"Contact", item.ContactInfo},
		{"Description", stringValue(item.Description)},
		{"Cover URL", stringValue(item.CoverURL)},
		{"Images", strings.Join(imageURLs, ", ")},
		{"Created At", formatTimeValue(item.CreatedAt)},
		{"Updated At", formatTimeValue(item.UpdatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runProductSetStatus(args []string) error {
	fs := flag.NewFlagSet("product set-status", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var statusRaw, reason string
	fs.UintVar(&id, "id", 0, "Product ID")
	fs.StringVar(&statusRaw, "status", "", "on-sale, trading, or off-shelf")
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
	status, err := parseEnumValue(statusRaw, productStatusMap, "status")
	if err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusProduct/updateCampusProductStatus", nil, map[string]any{
		"id":          id,
		"status":      status,
		"auditReason": reason,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("updated product %d status to %s", id, statusRaw), map[string]any{"id": id, "status": statusRaw})
}

func (a *App) runReport(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("report")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runReportList(args[1:])
	case "show":
		return a.runReportShow(args[1:])
	case "handle":
		return a.runReportHandle(args[1:])
	default:
		return fmt.Errorf("unknown report action %q", args[0])
	}
}

func (a *App) runReportList(args []string) error {
	fs := flag.NewFlagSet("report list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var id, reporterID, targetTypeRaw, targetID, reasonRaw, statusRaw, createdFromRaw, createdToRaw string
	fs.StringVar(&id, "id", "", "Report ID")
	fs.StringVar(&reporterID, "reporter-id", "", "Reporter ID")
	fs.StringVar(&targetTypeRaw, "target-type", "", "product")
	fs.StringVar(&targetID, "target-id", "", "Target ID")
	fs.StringVar(&reasonRaw, "reason", "", "fake, violation, or fraud")
	fs.StringVar(&statusRaw, "status", "", "pending or handled")
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
	for key, value := range map[string]string{"id": id, "reporterId": reporterID, "targetId": targetID} {
		if value != "" {
			values.Set(key, value)
		}
	}
	if targetTypeRaw != "" {
		targetType, err := parseEnumValue(targetTypeRaw, reportTargetTypeMap, "target-type")
		if err != nil {
			return err
		}
		values.Set("targetType", strconv.Itoa(targetType))
	}
	if reasonRaw != "" {
		reason, err := parseEnumValue(reasonRaw, reportReasonMap, "reason")
		if err != nil {
			return err
		}
		values.Set("reason", strconv.Itoa(reason))
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, reportStatusMap, "status")
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
	var result pageResponse[campus.CampusReport]
	headers, err := client.doJSON(http.MethodGet, "campusReport/getCampusReportList", values, nil, &result)
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
			strconv.FormatUint(uint64(item.ReporterID), 10),
			item.TargetTypeText,
			strconv.FormatUint(uint64(item.TargetID), 10),
			item.ReasonText,
			item.StatusText,
			stringValue(item.ReporterNickname),
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "REPORTER_ID", "TARGET_TYPE", "TARGET_ID", "REASON", "STATUS", "REPORTER", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runReportShow(args []string) error {
	fs := flag.NewFlagSet("report show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Report ID")
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
	var item campus.CampusReport
	headers, err := client.doJSON(http.MethodGet, "campusReport/findCampusReport", values, nil, &item)
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
		{"Reporter ID", strconv.FormatUint(uint64(item.ReporterID), 10)},
		{"Reporter", stringValue(item.ReporterNickname)},
		{"Reporter Phone", stringValue(item.ReporterPhone)},
		{"Target Type", item.TargetTypeText},
		{"Target ID", strconv.FormatUint(uint64(item.TargetID), 10)},
		{"Target Product", stringValue(item.TargetProductTitle)},
		{"Reason", item.ReasonText},
		{"Description", stringValue(item.Description)},
		{"Status", item.StatusText},
		{"Handled By", stringValue(item.HandledByNickname)},
		{"Handled By Phone", stringValue(item.HandledByPhone)},
		{"Handle Result", stringValue(item.HandleResult)},
		{"Created At", formatTimeValue(item.CreatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runReportHandle(args []string) error {
	fs := flag.NewFlagSet("report handle", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var result, statusRaw, reason string
	fs.UintVar(&id, "id", 0, "Report ID")
	fs.StringVar(&result, "result", "", "Handle result")
	fs.StringVar(&statusRaw, "status", "handled", "pending or handled")
	fs.StringVar(&reason, "reason", "", "Audit reason")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if err := requireID(id, "--id"); err != nil {
		return err
	}
	if strings.TrimSpace(result) == "" {
		return errors.New("--result is required")
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("--reason is required")
	}
	status, err := parseEnumValue(statusRaw, reportStatusMap, "status")
	if err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusReport/handleCampusReport", nil, map[string]any{
		"id":           id,
		"status":       status,
		"handleResult": result,
		"auditReason":  reason,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("updated report %d status to %s", id, statusRaw), map[string]any{"id": id, "status": statusRaw})
}

func (a *App) runCategory(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("category")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runCategoryList(args[1:])
	case "show":
		return a.runCategoryShow(args[1:])
	case "create":
		return a.runCategoryCreate(args[1:])
	case "update":
		return a.runCategoryUpdate(args[1:])
	case "enable":
		return a.runCategoryStatus(args[1:], true)
	case "disable":
		return a.runCategoryStatus(args[1:], false)
	default:
		return fmt.Errorf("unknown category action %q", args[0])
	}
}

func (a *App) runCategoryList(args []string) error {
	fs := flag.NewFlagSet("category list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var name, statusRaw string
	fs.StringVar(&name, "name", "", "Category name")
	fs.StringVar(&statusRaw, "status", "", "enabled or disabled")
	if err := fs.Parse(args); err != nil {
		return err
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	values := url.Values{}
	if name != "" {
		values.Set("name", name)
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, statusMap, "status")
		if err != nil {
			return err
		}
		values.Set("status", strconv.Itoa(status))
	}
	var items []campus.CampusCategory
	headers, err := client.doJSON(http.MethodGet, "campusCategory/getCampusCategoryList", values, nil, &items)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	if flags.json {
		return printJSON(a.stdout, marshalForJSON(items))
	}
	rows := make([][]string, 0)
	var walk func([]campus.CampusCategory, int)
	walk = func(nodes []campus.CampusCategory, depth int) {
		prefix := strings.Repeat("  ", depth)
		for _, item := range nodes {
			rows = append(rows, []string{
				strconv.FormatUint(uint64(item.ID), 10),
				prefix + item.Name,
				uintPtrValue(item.ParentID),
				stringValue(item.Icon),
				strconv.Itoa(item.SortOrder),
				item.StatusText,
				formatTimeValue(item.CreatedAt),
			})
			walk(item.Children, depth+1)
		}
	}
	walk(items, 0)
	printTable(a.stdout, []string{"ID", "NAME", "PARENT_ID", "ICON", "SORT_ORDER", "STATUS", "CREATED_AT"}, rows)
	return nil
}

func (a *App) runCategoryShow(args []string) error {
	fs := flag.NewFlagSet("category show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Category ID")
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
	var item campus.CampusCategory
	headers, err := client.doJSON(http.MethodGet, "campusCategory/findCampusCategory", values, nil, &item)
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
		{"Name", item.Name},
		{"Parent ID", uintPtrValue(item.ParentID)},
		{"Parent Name", stringValue(item.ParentName)},
		{"Sort Order", strconv.Itoa(item.SortOrder)},
		{"Icon", stringValue(item.Icon)},
		{"Status", item.StatusText},
		{"Created At", formatTimeValue(item.CreatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runCategoryCreate(args []string) error {
	fs := flag.NewFlagSet("category create", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var name, icon, statusRaw string
	var parentID uint
	var sortOrder int
	fs.StringVar(&name, "name", "", "Category name")
	fs.UintVar(&parentID, "parent-id", 0, "Parent category ID")
	fs.IntVar(&sortOrder, "sort-order", 0, "Sort order")
	fs.StringVar(&icon, "icon", "", "Icon")
	fs.StringVar(&statusRaw, "status", "enabled", "enabled or disabled")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("--name is required")
	}
	status, err := parseEnumValue(statusRaw, statusMap, "status")
	if err != nil {
		return err
	}
	body := map[string]any{
		"name":      name,
		"sortOrder": sortOrder,
		"icon":      icon,
		"status":    status,
	}
	if parentID != 0 {
		body["parentId"] = parentID
	} else {
		body["parentId"] = nil
	}
	client, profile, err := a.newAuthedClient(flags)
	if err != nil {
		return err
	}
	headers, err := client.doJSON(http.MethodPost, "campusCategory/createCampusCategory", nil, body, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("created category %q", name), map[string]any{"name": name})
}

func (a *App) runCategoryUpdate(args []string) error {
	fs := flag.NewFlagSet("category update", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id, parentID uint
	var name, icon, statusRaw string
	var sortOrder int
	var parentSet, sortSet bool
	fs.UintVar(&id, "id", 0, "Category ID")
	fs.StringVar(&name, "name", "", "Category name")
	fs.StringVar(&icon, "icon", "", "Icon")
	fs.StringVar(&statusRaw, "status", "", "enabled or disabled")
	fs.Func("parent-id", "Parent category ID, use 0 for top-level", func(value string) error {
		parsed, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		parentID = uint(parsed)
		parentSet = true
		return nil
	})
	fs.Func("sort-order", "Sort order", func(value string) error {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		sortOrder = parsed
		sortSet = true
		return nil
	})
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
	current, headerToken, err := a.findCategory(client, id)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headerToken); err != nil {
		return err
	}
	if name != "" {
		current.Name = name
	}
	if parentSet {
		if parentID == 0 {
			current.ParentID = nil
		} else {
			current.ParentID = &parentID
		}
	}
	if sortSet {
		current.SortOrder = sortOrder
	}
	if icon != "" {
		current.Icon = &icon
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, statusMap, "status")
		if err != nil {
			return err
		}
		current.Status = status
	}
	body := map[string]any{
		"id":        current.ID,
		"name":      current.Name,
		"parentId":  current.ParentID,
		"sortOrder": current.SortOrder,
		"icon":      chooseOptionalString(current.Icon),
		"status":    current.Status,
	}
	headers, err := client.doJSON(http.MethodPut, "campusCategory/updateCampusCategory", nil, body, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("updated category %d", id), map[string]any{"id": id})
}

func (a *App) runCategoryStatus(args []string, enabled bool) error {
	action := "disable"
	status := 1
	if enabled {
		action = "enable"
		status = 0
	}
	fs := flag.NewFlagSet("category "+action, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Category ID")
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
	headers, err := client.doJSON(http.MethodPost, "campusCategory/updateCampusCategoryStatus", nil, map[string]any{
		"id":     id,
		"status": status,
	}, nil)
	if err != nil {
		return err
	}
	if err := a.saveUpdatedProfile(profile, headers); err != nil {
		return err
	}
	return a.printActionResult(flags.json, fmt.Sprintf("%sd category %d", action, id), map[string]any{"id": id})
}

func (a *App) findCategory(client *apiClient, id uint) (campus.CampusCategory, http.Header, error) {
	values := url.Values{"id": []string{strconv.FormatUint(uint64(id), 10)}}
	var item campus.CampusCategory
	headers, err := client.doJSON(http.MethodGet, "campusCategory/findCampusCategory", values, nil, &item)
	return item, headers, err
}
