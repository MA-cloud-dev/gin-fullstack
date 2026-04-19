package campus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	campusModel "github.com/flipped-aurora/gin-vue-admin/server/model/campus"
	campusResp "github.com/flipped-aurora/gin-vue-admin/server/model/campus/response"
)

type CampusOverviewService struct{}

type campusOverviewCounts struct {
	UserTotal                int64
	UserEnabledTotal         int64
	UserDisabledTotal        int64
	VerifiedUserTotal        int64
	PendingAuthTotal         int64
	AuthTotal                int64
	AuthReviewedTotal        int64
	ProductTotal             int64
	ProductOnSaleTotal       int64
	ProductPendingTotal      int64
	ProductSoldTotal         int64
	OrderTotal               int64
	OrderPendingPayTotal     int64
	OrderPendingConfirmTotal int64
	OrderCompletedTotal      int64
	ReportTotal              int64
	PendingReportTotal       int64
	HandledReportTotal       int64
	AnnouncementTotal        int64
	AnnouncementOnlineTotal  int64
	AnnouncementOfflineTotal int64
	CategoryTotal            int64
	CategoryEnabledTotal     int64
	CategoryDisabledTotal    int64
	CategoryCreatedLast7     int64
	StaffTotal               int64
	StaffEnabledTotal        int64
	StaffDisabledTotal       int64
	StaffCreatedLast7        int64
	OperationLogTotal        int64
}

type campusOverviewTrendSums struct {
	NewUsers         int64
	AuthApplications int64
	AuthApproved     int64
	NewProducts      int64
	NewOrders        int64
	NewReports       int64
	NewAnnouncements int64
	OperationLogs    int64
}

type campusOverviewUserAuthSource struct {
	Role       int `gorm:"column:role"`
	AuthStatus int `gorm:"column:auth_status"`
}

type campusOverviewStatusCount struct {
	Status int   `gorm:"column:status"`
	Count  int64 `gorm:"column:count"`
}

func (s *CampusOverviewService) GetCampusOverview(ctx context.Context) (campusResp.CampusOverviewResponse, error) {
	location := campusOverviewLocation()
	now := time.Now().In(location)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location)
	rangeStart := dayStart.AddDate(0, 0, -6)
	rangeEnd := dayStart.AddDate(0, 0, 1)

	counts, err := s.collectOverviewCounts(ctx, rangeStart, rangeEnd)
	if err != nil {
		return campusResp.CampusOverviewResponse{}, err
	}

	dailyTrend, err := s.buildDailyTrend(ctx, location, rangeStart, rangeEnd)
	if err != nil {
		return campusResp.CampusOverviewResponse{}, err
	}

	distributions, err := s.buildDistributions(ctx)
	if err != nil {
		return campusResp.CampusOverviewResponse{}, err
	}

	trendSum := sumDailyTrend(dailyTrend)

	return campusResp.CampusOverviewResponse{
		Kpis: campusResp.CampusOverviewKpis{
			UserTotal:               counts.UserTotal,
			VerifiedUserTotal:       counts.VerifiedUserTotal,
			PendingAuthTotal:        counts.PendingAuthTotal,
			ProductOnSaleTotal:      counts.ProductOnSaleTotal,
			OrderTotal:              counts.OrderTotal,
			OrderCompletedTotal:     counts.OrderCompletedTotal,
			PendingReportTotal:      counts.PendingReportTotal,
			AnnouncementOnlineTotal: counts.AnnouncementOnlineTotal,
		},
		ModuleTotals: []campusResp.CampusOverviewModuleRow{
			{
				Module:           "用户",
				Total:            counts.UserTotal,
				CoreStatus1Label: "启用",
				CoreStatus1Count: counts.UserEnabledTotal,
				CoreStatus2Label: "禁用",
				CoreStatus2Count: counts.UserDisabledTotal,
				CoreStatus3Label: "已认证",
				CoreStatus3Count: counts.VerifiedUserTotal,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.NewUsers,
			},
			{
				Module:           "认证",
				Total:            counts.AuthTotal,
				CoreStatus1Label: "审核中",
				CoreStatus1Count: counts.PendingAuthTotal,
				CoreStatus2Label: "已审核",
				CoreStatus2Count: counts.AuthReviewedTotal,
				CoreStatus3Label: "近7天通过",
				CoreStatus3Count: trendSum.AuthApproved,
				Last7DaysLabel:   "近7天申请",
				Last7DaysCount:   trendSum.AuthApplications,
			},
			{
				Module:           "商品",
				Total:            counts.ProductTotal,
				CoreStatus1Label: "在售",
				CoreStatus1Count: counts.ProductOnSaleTotal,
				CoreStatus2Label: "待审核",
				CoreStatus2Count: counts.ProductPendingTotal,
				CoreStatus3Label: "已售出",
				CoreStatus3Count: counts.ProductSoldTotal,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.NewProducts,
			},
			{
				Module:           "订单",
				Total:            counts.OrderTotal,
				CoreStatus1Label: "待付款",
				CoreStatus1Count: counts.OrderPendingPayTotal,
				CoreStatus2Label: "待确认",
				CoreStatus2Count: counts.OrderPendingConfirmTotal,
				CoreStatus3Label: "已完成",
				CoreStatus3Count: counts.OrderCompletedTotal,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.NewOrders,
			},
			{
				Module:           "举报",
				Total:            counts.ReportTotal,
				CoreStatus1Label: "待处理",
				CoreStatus1Count: counts.PendingReportTotal,
				CoreStatus2Label: "已处理",
				CoreStatus2Count: counts.HandledReportTotal,
				CoreStatus3Label: "-",
				CoreStatus3Count: 0,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.NewReports,
			},
			{
				Module:           "公告",
				Total:            counts.AnnouncementTotal,
				CoreStatus1Label: "上线",
				CoreStatus1Count: counts.AnnouncementOnlineTotal,
				CoreStatus2Label: "下线",
				CoreStatus2Count: counts.AnnouncementOfflineTotal,
				CoreStatus3Label: "-",
				CoreStatus3Count: 0,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.NewAnnouncements,
			},
			{
				Module:           "分类",
				Total:            counts.CategoryTotal,
				CoreStatus1Label: "启用",
				CoreStatus1Count: counts.CategoryEnabledTotal,
				CoreStatus2Label: "禁用",
				CoreStatus2Count: counts.CategoryDisabledTotal,
				CoreStatus3Label: "-",
				CoreStatus3Count: 0,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   counts.CategoryCreatedLast7,
			},
			{
				Module:           "B端管理员",
				Total:            counts.StaffTotal,
				CoreStatus1Label: "启用",
				CoreStatus1Count: counts.StaffEnabledTotal,
				CoreStatus2Label: "禁用",
				CoreStatus2Count: counts.StaffDisabledTotal,
				CoreStatus3Label: "-",
				CoreStatus3Count: 0,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   counts.StaffCreatedLast7,
			},
			{
				Module:           "操作记录",
				Total:            counts.OperationLogTotal,
				CoreStatus1Label: "近7天总数",
				CoreStatus1Count: trendSum.OperationLogs,
				CoreStatus2Label: "-",
				CoreStatus2Count: 0,
				CoreStatus3Label: "-",
				CoreStatus3Count: 0,
				Last7DaysLabel:   "近7天新增",
				Last7DaysCount:   trendSum.OperationLogs,
			},
		},
		DailyTrend:    dailyTrend,
		Distributions: distributions,
		GeneratedAt:   now,
	}, nil
}

func (s *CampusOverviewService) collectOverviewCounts(ctx context.Context, rangeStart time.Time, rangeEnd time.Time) (campusOverviewCounts, error) {
	db := global.GVA_DB.WithContext(ctx)
	counts := campusOverviewCounts{}

	if err := db.Table("t_user").Count(&counts.UserTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_user").Where("status = ?", 0).Count(&counts.UserEnabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_user").Where("status = ?", 1).Count(&counts.UserDisabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_user").Where("auth_status = ?", campusModel.CampusUserAuthStatusApproved).Count(&counts.VerifiedUserTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_campus_auth").Where("review_status = ?", campusModel.CampusAuthReviewStatusProcessing).Count(&counts.PendingAuthTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_campus_auth").Count(&counts.AuthTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_campus_auth").Where("review_status IN ?", []string{campusModel.CampusAuthReviewStatusApproved, campusModel.CampusAuthReviewStatusRejected}).Count(&counts.AuthReviewedTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_product").Count(&counts.ProductTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_product").Where("status = ?", 0).Count(&counts.ProductOnSaleTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_product").Where("status = ?", 1).Count(&counts.ProductPendingTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_product").Where("status = ?", 4).Count(&counts.ProductSoldTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_order").Count(&counts.OrderTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_order").Where("status = ?", 1).Count(&counts.OrderPendingPayTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_order").Where("status = ?", 2).Count(&counts.OrderPendingConfirmTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_order").Where("status = ?", 3).Count(&counts.OrderCompletedTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_report").Count(&counts.ReportTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_report").Where("status = ?", 0).Count(&counts.PendingReportTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_report").Where("status = ?", 1).Count(&counts.HandledReportTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_announcement").Count(&counts.AnnouncementTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_announcement").Where("status = ?", 1).Count(&counts.AnnouncementOnlineTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_announcement").Where("status = ?", 0).Count(&counts.AnnouncementOfflineTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_category").Count(&counts.CategoryTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_category").Where("status = ?", 0).Count(&counts.CategoryEnabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_category").Where("status = ?", 1).Count(&counts.CategoryDisabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_category").Where("created_at >= ? AND created_at < ?", rangeStart, rangeEnd).Count(&counts.CategoryCreatedLast7).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_admin_staff").Count(&counts.StaffTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_admin_staff").Where("status = ?", 0).Count(&counts.StaffEnabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_admin_staff").Where("status = ?", 1).Count(&counts.StaffDisabledTotal).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_admin_staff").Where("created_at >= ? AND created_at < ?", rangeStart, rangeEnd).Count(&counts.StaffCreatedLast7).Error; err != nil {
		return counts, err
	}
	if err := db.Table("t_campus_operation_log").Count(&counts.OperationLogTotal).Error; err != nil {
		return counts, err
	}

	return counts, nil
}

func (s *CampusOverviewService) buildDailyTrend(ctx context.Context, location *time.Location, start, end time.Time) ([]campusResp.CampusOverviewDailyTrendItem, error) {
	base := make(map[string]*campusResp.CampusOverviewDailyTrendItem, 7)
	items := make([]campusResp.CampusOverviewDailyTrendItem, 0, 7)
	for i := 0; i < 7; i++ {
		day := start.AddDate(0, 0, i)
		key := day.Format("2006-01-02")
		item := campusResp.CampusOverviewDailyTrendItem{Date: key}
		items = append(items, item)
		base[key] = &items[len(items)-1]
	}

	if err := s.accumulateTimeColumn(ctx, "t_user", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.NewUsers++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_campus_auth", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.AuthApplications++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumnWhere(ctx, "t_campus_auth", "reviewed_at", start, end, location, "review_status = ?", []interface{}{campusModel.CampusAuthReviewStatusApproved}, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.AuthApproved++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_product", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.NewProducts++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_order", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.NewOrders++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_report", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.NewReports++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_announcement", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.NewAnnouncements++
	}, base); err != nil {
		return nil, err
	}
	if err := s.accumulateTimeColumn(ctx, "t_campus_operation_log", "created_at", start, end, location, func(item *campusResp.CampusOverviewDailyTrendItem) {
		item.OperationLogs++
	}, base); err != nil {
		return nil, err
	}

	return items, nil
}

func (s *CampusOverviewService) accumulateTimeColumn(
	ctx context.Context,
	table string,
	column string,
	start time.Time,
	end time.Time,
	location *time.Location,
	apply func(item *campusResp.CampusOverviewDailyTrendItem),
	base map[string]*campusResp.CampusOverviewDailyTrendItem,
) error {
	var values []time.Time
	if err := global.GVA_DB.WithContext(ctx).
		Table(table).
		Where(fmt.Sprintf("%s >= ? AND %s < ?", column, column), start, end).
		Pluck(column, &values).Error; err != nil {
		return err
	}

	for _, value := range values {
		key := value.In(location).Format("2006-01-02")
		if item, ok := base[key]; ok {
			apply(item)
		}
	}
	return nil
}

func (s *CampusOverviewService) buildDistributions(ctx context.Context) (campusResp.CampusOverviewDistributions, error) {
	userAuthStatus, err := s.buildUserAuthStatusDistribution(ctx)
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}
	productStatus, err := s.buildStatusDistribution(ctx, "t_product", buildProductStatusBuckets())
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}
	orderStatus, err := s.buildStatusDistribution(ctx, "t_order", buildOrderStatusBuckets())
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}
	reportStatus, err := s.buildStatusDistribution(ctx, "t_report", buildReportStatusBuckets())
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}
	announcementStatus, err := s.buildStatusDistribution(ctx, "t_announcement", buildAnnouncementStatusBuckets())
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}
	staffStatus, err := s.buildStatusDistribution(ctx, "t_admin_staff", buildEnabledStatusBuckets())
	if err != nil {
		return campusResp.CampusOverviewDistributions{}, err
	}

	return campusResp.CampusOverviewDistributions{
		UserAuthStatus:     userAuthStatus,
		ProductStatus:      productStatus,
		OrderStatus:        orderStatus,
		ReportStatus:       reportStatus,
		AnnouncementStatus: announcementStatus,
		StaffStatus:        staffStatus,
	}, nil
}

func (s *CampusOverviewService) buildUserAuthStatusDistribution(ctx context.Context) ([]campusResp.CampusOverviewBucketItem, error) {
	var sources []campusOverviewUserAuthSource
	if err := global.GVA_DB.WithContext(ctx).
		Table("t_user AS u").
		Select("u.role, u.auth_status").
		Find(&sources).Error; err != nil {
		return nil, err
	}

	countMap := map[string]int64{
		"unverified": 0,
		"pending":    0,
		"verified":   0,
		"admin":      0,
	}
	for _, source := range sources {
		user := campusModel.CampusUser{
			Role:       source.Role,
			AuthStatus: source.AuthStatus,
		}
		switch buildAuthStatusText(user) {
		case "审核中":
			countMap["pending"]++
		case "已认证":
			countMap["verified"]++
		case "管理员":
			countMap["admin"]++
		default:
			countMap["unverified"]++
		}
	}

	return []campusResp.CampusOverviewBucketItem{
		{Key: "unverified", Label: "未认证", Count: countMap["unverified"]},
		{Key: "pending", Label: "审核中", Count: countMap["pending"]},
		{Key: "verified", Label: "已认证", Count: countMap["verified"]},
		{Key: "admin", Label: "管理员", Count: countMap["admin"]},
	}, nil
}

func (s *CampusOverviewService) accumulateTimeColumnWhere(
	ctx context.Context,
	table string,
	column string,
	start time.Time,
	end time.Time,
	location *time.Location,
	where string,
	args []interface{},
	apply func(item *campusResp.CampusOverviewDailyTrendItem),
	base map[string]*campusResp.CampusOverviewDailyTrendItem,
) error {
	var values []time.Time
	query := global.GVA_DB.WithContext(ctx).
		Table(table).
		Where(fmt.Sprintf("%s >= ? AND %s < ?", column, column), start, end)
	if strings.TrimSpace(where) != "" {
		query = query.Where(where, args...)
	}
	if err := query.Pluck(column, &values).Error; err != nil {
		return err
	}

	for _, value := range values {
		key := value.In(location).Format("2006-01-02")
		if item, ok := base[key]; ok {
			apply(item)
		}
	}
	return nil
}

func (s *CampusOverviewService) buildStatusDistribution(ctx context.Context, table string, buckets []campusResp.CampusOverviewBucketItem) ([]campusResp.CampusOverviewBucketItem, error) {
	var rows []campusOverviewStatusCount
	if err := global.GVA_DB.WithContext(ctx).
		Table(table).
		Select("status, COUNT(*) AS count").
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int]int64, len(rows))
	for _, row := range rows {
		countMap[row.Status] = row.Count
	}
	for i := range buckets {
		var status int
		if _, err := fmt.Sscanf(buckets[i].Key, "%d", &status); err == nil {
			buckets[i].Count = countMap[status]
		}
	}
	return buckets, nil
}

func sumDailyTrend(items []campusResp.CampusOverviewDailyTrendItem) campusOverviewTrendSums {
	sums := campusOverviewTrendSums{}
	for _, item := range items {
		sums.NewUsers += item.NewUsers
		sums.AuthApplications += item.AuthApplications
		sums.AuthApproved += item.AuthApproved
		sums.NewProducts += item.NewProducts
		sums.NewOrders += item.NewOrders
		sums.NewReports += item.NewReports
		sums.NewAnnouncements += item.NewAnnouncements
		sums.OperationLogs += item.OperationLogs
	}
	return sums
}

func buildEnabledStatusBuckets() []campusResp.CampusOverviewBucketItem {
	return []campusResp.CampusOverviewBucketItem{
		{Key: "0", Label: "启用"},
		{Key: "1", Label: "禁用"},
	}
}

func buildAnnouncementStatusBuckets() []campusResp.CampusOverviewBucketItem {
	return []campusResp.CampusOverviewBucketItem{
		{Key: "1", Label: "上线"},
		{Key: "0", Label: "下线"},
	}
}

func buildReportStatusBuckets() []campusResp.CampusOverviewBucketItem {
	return []campusResp.CampusOverviewBucketItem{
		{Key: "0", Label: "待处理"},
		{Key: "1", Label: "已处理"},
	}
}

func buildOrderStatusBuckets() []campusResp.CampusOverviewBucketItem {
	return []campusResp.CampusOverviewBucketItem{
		{Key: "1", Label: "待付款"},
		{Key: "2", Label: "待确认"},
		{Key: "3", Label: "已完成"},
		{Key: "4", Label: "已取消"},
		{Key: "5", Label: "已关闭"},
	}
}

func buildProductStatusBuckets() []campusResp.CampusOverviewBucketItem {
	return []campusResp.CampusOverviewBucketItem{
		{Key: "0", Label: "在售"},
		{Key: "1", Label: "待审核"},
		{Key: "2", Label: "交易中"},
		{Key: "3", Label: "已下架"},
		{Key: "4", Label: "已售出"},
	}
}

func campusOverviewLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		return location
	}
	return time.FixedZone("Asia/Shanghai", 8*60*60)
}
