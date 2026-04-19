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
	userRoleMap       = map[string]int{"normal": 0, "admin": 1}
	userAuthStatusMap = map[string]int{"unauth": 0, "pending": 2, "verified": 3}
	staffRoleTypeMap  = map[string]int{"super": 1, "ops": 2}
	orderStatusMap    = map[string]int{"unpaid": 1, "unconfirmed": 2, "completed": 3, "cancelled": 4, "closed": 5}
)

func (a *App) runUser(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("user")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runUserList(args[1:])
	case "show":
		return a.runUserShow(args[1:])
	case "enable":
		return a.runUserStatus(args[1:], true)
	case "disable":
		return a.runUserStatus(args[1:], false)
	default:
		return fmt.Errorf("unknown user action %q", args[0])
	}
}

func (a *App) runUserList(args []string) error {
	fs := flag.NewFlagSet("user list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var id, phone, nickname, roleRaw, statusRaw, authStatusRaw, studentID, realName, college, createdFromRaw, createdToRaw string
	fs.StringVar(&id, "id", "", "User ID")
	fs.StringVar(&phone, "phone", "", "Phone")
	fs.StringVar(&nickname, "nickname", "", "Nickname")
	fs.StringVar(&roleRaw, "role", "", "normal or admin")
	fs.StringVar(&statusRaw, "status", "", "enabled or disabled")
	fs.StringVar(&authStatusRaw, "auth-status", "", "unauth, pending, or verified")
	fs.StringVar(&studentID, "student-id", "", "Student ID")
	fs.StringVar(&realName, "real-name", "", "Real name")
	fs.StringVar(&college, "college", "", "College")
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
	for key, value := range map[string]string{
		"id":        id,
		"phone":     phone,
		"nickname":  nickname,
		"studentId": studentID,
		"realName":  realName,
		"college":   college,
	} {
		if value != "" {
			values.Set(key, value)
		}
	}
	if roleRaw != "" {
		role, err := parseEnumValue(roleRaw, userRoleMap, "role")
		if err != nil {
			return err
		}
		values.Set("role", strconv.Itoa(role))
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, statusMap, "status")
		if err != nil {
			return err
		}
		values.Set("status", strconv.Itoa(status))
	}
	if authStatusRaw != "" {
		authStatus, err := parseEnumValue(authStatusRaw, userAuthStatusMap, "auth-status")
		if err != nil {
			return err
		}
		values.Set("authStatus", strconv.Itoa(authStatus))
	}
	createdFrom, createdTo, err := parseDateRange(createdFromRaw, createdToRaw)
	if err != nil {
		return err
	}
	maybeAddRange(values, "createdAtRange[]", createdFrom, createdTo)
	var result pageResponse[campus.CampusUser]
	headers, err := client.doJSON(http.MethodGet, "campusUser/getCampusUserList", values, nil, &result)
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
			item.Phone,
			stringValue(item.Nickname),
			item.RoleText,
			item.StatusText,
			item.AuthStatusText,
			stringValue(item.StudentID),
			stringValue(item.RealName),
			stringValue(item.College),
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "PHONE", "NICKNAME", "ROLE", "STATUS", "AUTH_STATUS", "STUDENT_ID", "REAL_NAME", "COLLEGE", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runUserShow(args []string) error {
	fs := flag.NewFlagSet("user show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "User ID")
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
	var item campus.CampusUser
	headers, err := client.doJSON(http.MethodGet, "campusUser/findCampusUser", values, nil, &item)
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
		{"Phone", item.Phone},
		{"Nickname", stringValue(item.Nickname)},
		{"Username", stringValue(item.Username)},
		{"Role", item.RoleText},
		{"Status", item.StatusText},
		{"Auth Status", item.AuthStatusText},
		{"Student ID", stringValue(item.StudentID)},
		{"Real Name", stringValue(item.RealName)},
		{"College", stringValue(item.College)},
		{"Grade", stringValue(item.Grade)},
		{"Dormitory", stringValue(item.Dormitory)},
		{"Wechat ID", stringValue(item.WechatID)},
		{"Review Remark", stringValue(item.ReviewRemark)},
		{"Reviewed By", stringValue(item.ReviewedByName)},
		{"Reviewed At", formatTimePtr(item.ReviewedAt)},
		{"Created At", formatTimeValue(item.CreatedAt)},
		{"Updated At", formatTimeValue(item.UpdatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runUserStatus(args []string, enabled bool) error {
	action := "disable"
	status := 1
	if enabled {
		action = "enable"
		status = 0
	}
	fs := flag.NewFlagSet("user "+action, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var reason string
	fs.UintVar(&id, "id", 0, "User ID")
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
	headers, err := client.doJSON(http.MethodPost, "campusUser/updateCampusUserStatus", nil, map[string]any{
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
	return a.printActionResult(flags.json, fmt.Sprintf("%sd user %d", action, id), map[string]any{"id": id})
}

func (a *App) runStaff(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("staff")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runStaffList(args[1:])
	case "show":
		return a.runStaffShow(args[1:])
	case "enable":
		return a.runStaffStatus(args[1:], true)
	case "disable":
		return a.runStaffStatus(args[1:], false)
	default:
		return fmt.Errorf("unknown staff action %q", args[0])
	}
}

func (a *App) runStaffList(args []string) error {
	fs := flag.NewFlagSet("staff list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var id, username, displayName, roleTypeRaw, statusRaw, createdFromRaw, createdToRaw string
	fs.StringVar(&id, "id", "", "Staff ID")
	fs.StringVar(&username, "username", "", "Username")
	fs.StringVar(&displayName, "display-name", "", "Display name")
	fs.StringVar(&roleTypeRaw, "role-type", "", "super or ops")
	fs.StringVar(&statusRaw, "status", "", "enabled or disabled")
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
	for key, value := range map[string]string{"id": id, "username": username, "displayName": displayName} {
		if value != "" {
			values.Set(key, value)
		}
	}
	if roleTypeRaw != "" {
		roleType, err := parseEnumValue(roleTypeRaw, staffRoleTypeMap, "role-type")
		if err != nil {
			return err
		}
		values.Set("roleType", strconv.Itoa(roleType))
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, statusMap, "status")
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
	var result pageResponse[campus.CampusAdminStaff]
	headers, err := client.doJSON(http.MethodGet, "campusAdminStaff/getCampusAdminStaffList", values, nil, &result)
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
			item.Username,
			item.DisplayName,
			item.RoleTypeText,
			item.StatusText,
			formatTimePtr(item.LastLoginAt),
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "USERNAME", "DISPLAY_NAME", "ROLE_TYPE", "STATUS", "LAST_LOGIN_AT", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runStaffShow(args []string) error {
	fs := flag.NewFlagSet("staff show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Staff ID")
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
	var item campus.CampusAdminStaff
	headers, err := client.doJSON(http.MethodGet, "campusAdminStaff/findCampusAdminStaff", values, nil, &item)
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
		{"Username", item.Username},
		{"Display Name", item.DisplayName},
		{"Role Type", item.RoleTypeText},
		{"Status", item.StatusText},
		{"Last Login At", formatTimePtr(item.LastLoginAt)},
		{"Created At", formatTimeValue(item.CreatedAt)},
		{"Updated At", formatTimeValue(item.UpdatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}

func (a *App) runStaffStatus(args []string, enabled bool) error {
	action := "disable"
	status := 1
	if enabled {
		action = "enable"
		status = 0
	}
	fs := flag.NewFlagSet("staff "+action, flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	var reason string
	fs.UintVar(&id, "id", 0, "Staff ID")
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
	headers, err := client.doJSON(http.MethodPost, "campusAdminStaff/updateCampusAdminStaffStatus", nil, map[string]any{
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
	return a.printActionResult(flags.json, fmt.Sprintf("%sd staff %d", action, id), map[string]any{"id": id})
}

func (a *App) runOrder(args []string) error {
	if len(args) == 0 || isHelpToken(args[0]) {
		a.printCommandUsage("order")
		return nil
	}
	switch args[0] {
	case "list":
		return a.runOrderList(args[1:])
	case "show":
		return a.runOrderShow(args[1:])
	default:
		return fmt.Errorf("unknown order action %q", args[0])
	}
}

func (a *App) runOrderList(args []string) error {
	fs := flag.NewFlagSet("order list", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	var page pageFlags
	addCommonFlags(fs, &flags)
	addPageFlags(fs, &page)
	var orderNo, buyerID, sellerID, productID, statusRaw, createdFromRaw, createdToRaw string
	fs.StringVar(&orderNo, "order-no", "", "Order number")
	fs.StringVar(&buyerID, "buyer-id", "", "Buyer ID")
	fs.StringVar(&sellerID, "seller-id", "", "Seller ID")
	fs.StringVar(&productID, "product-id", "", "Product ID")
	fs.StringVar(&statusRaw, "status", "", "unpaid, unconfirmed, completed, cancelled, or closed")
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
	for key, value := range map[string]string{"orderNo": orderNo, "buyerId": buyerID, "sellerId": sellerID, "productId": productID} {
		if value != "" {
			values.Set(key, value)
		}
	}
	if statusRaw != "" {
		status, err := parseEnumValue(statusRaw, orderStatusMap, "status")
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
	var result pageResponse[campus.CampusOrder]
	headers, err := client.doJSON(http.MethodGet, "campusOrder/getCampusOrderList", values, nil, &result)
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
			item.OrderNo,
			item.ProductTitle,
			stringValue(item.BuyerNickname),
			stringValue(item.SellerNickname),
			fmt.Sprintf("%.2f", item.Price),
			item.StatusText,
			formatTimeValue(item.CreatedAt),
		})
	}
	printTable(a.stdout, []string{"ID", "ORDER_NO", "PRODUCT_TITLE", "BUYER", "SELLER", "PRICE", "STATUS", "CREATED_AT"}, rows)
	fmt.Fprintf(a.stdout, "\nTotal: %d  Page: %d  PageSize: %d\n", result.Total, result.Page, result.PageSize)
	return nil
}

func (a *App) runOrderShow(args []string) error {
	fs := flag.NewFlagSet("order show", flag.ContinueOnError)
	fs.SetOutput(a.stderr)
	var flags commonFlags
	addCommonFlags(fs, &flags)
	var id uint
	fs.UintVar(&id, "id", 0, "Order ID")
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
	var item campus.CampusOrder
	headers, err := client.doJSON(http.MethodGet, "campusOrder/findCampusOrder", values, nil, &item)
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
		{"Order No", item.OrderNo},
		{"Product ID", strconv.FormatUint(uint64(item.ProductID), 10)},
		{"Product Title", item.ProductTitle},
		{"Product Image", stringValue(item.ProductImage)},
		{"Buyer ID", strconv.FormatUint(uint64(item.BuyerID), 10)},
		{"Buyer", stringValue(item.BuyerNickname)},
		{"Buyer Phone", stringValue(item.BuyerPhone)},
		{"Seller ID", strconv.FormatUint(uint64(item.SellerID), 10)},
		{"Seller", stringValue(item.SellerNickname)},
		{"Seller Phone", stringValue(item.SellerPhone)},
		{"Price", fmt.Sprintf("%.2f", item.Price)},
		{"Status", item.StatusText},
		{"Remark", stringValue(item.Remark)},
		{"Close Reason", stringValue(item.CloseReason)},
		{"Close By", item.CloseByText},
		{"Close Confirmed", item.CloseConfirmedText},
		{"Confirmed At", formatTimePtr(item.ConfirmedAt)},
		{"Completed At", formatTimePtr(item.CompletedAt)},
		{"Cancelled At", formatTimePtr(item.CancelledAt)},
		{"Created At", formatTimeValue(item.CreatedAt)},
		{"Updated At", formatTimeValue(item.UpdatedAt)},
	}
	printTable(a.stdout, []string{"FIELD", "VALUE"}, rows)
	return nil
}
