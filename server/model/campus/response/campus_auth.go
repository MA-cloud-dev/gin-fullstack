package response

type SubmitCampusAuthResp struct {
	AuthRecordID uint64 `json:"authRecordId"`
	TaskID       string `json:"taskId"`
	ReviewStatus string `json:"reviewStatus"`
	AuthStatus   int    `json:"authStatus"`
}
