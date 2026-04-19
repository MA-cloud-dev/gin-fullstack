package response

import "time"

type CampusOverviewResponse struct {
	Kpis          CampusOverviewKpis             `json:"kpis"`
	ModuleTotals  []CampusOverviewModuleRow      `json:"moduleTotals"`
	DailyTrend    []CampusOverviewDailyTrendItem `json:"dailyTrend"`
	Distributions CampusOverviewDistributions    `json:"distributions"`
	GeneratedAt   time.Time                      `json:"generatedAt"`
}

type CampusOverviewKpis struct {
	UserTotal               int64 `json:"userTotal"`
	VerifiedUserTotal       int64 `json:"verifiedUserTotal"`
	PendingAuthTotal        int64 `json:"pendingAuthTotal"`
	ProductOnSaleTotal      int64 `json:"productOnSaleTotal"`
	OrderTotal              int64 `json:"orderTotal"`
	OrderCompletedTotal     int64 `json:"orderCompletedTotal"`
	PendingReportTotal      int64 `json:"pendingReportTotal"`
	AnnouncementOnlineTotal int64 `json:"announcementOnlineTotal"`
}

type CampusOverviewModuleRow struct {
	Module           string `json:"module"`
	Total            int64  `json:"total"`
	CoreStatus1Label string `json:"coreStatus1Label"`
	CoreStatus1Count int64  `json:"coreStatus1Count"`
	CoreStatus2Label string `json:"coreStatus2Label"`
	CoreStatus2Count int64  `json:"coreStatus2Count"`
	CoreStatus3Label string `json:"coreStatus3Label"`
	CoreStatus3Count int64  `json:"coreStatus3Count"`
	Last7DaysLabel   string `json:"last7DaysLabel"`
	Last7DaysCount   int64  `json:"last7DaysCount"`
}

type CampusOverviewDailyTrendItem struct {
	Date             string `json:"date"`
	NewUsers         int64  `json:"newUsers"`
	AuthApplications int64  `json:"authApplications"`
	AuthApproved     int64  `json:"authApproved"`
	NewProducts      int64  `json:"newProducts"`
	NewOrders        int64  `json:"newOrders"`
	NewReports       int64  `json:"newReports"`
	NewAnnouncements int64  `json:"newAnnouncements"`
	OperationLogs    int64  `json:"operationLogs"`
}

type CampusOverviewBucketItem struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type CampusOverviewDistributions struct {
	UserAuthStatus     []CampusOverviewBucketItem `json:"userAuthStatus"`
	ProductStatus      []CampusOverviewBucketItem `json:"productStatus"`
	OrderStatus        []CampusOverviewBucketItem `json:"orderStatus"`
	ReportStatus       []CampusOverviewBucketItem `json:"reportStatus"`
	AnnouncementStatus []CampusOverviewBucketItem `json:"announcementStatus"`
	StaffStatus        []CampusOverviewBucketItem `json:"staffStatus"`
}
