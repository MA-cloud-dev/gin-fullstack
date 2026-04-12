package request

type CampusCategorySearch struct {
	Name   string `json:"name" form:"name"`
	Status *int   `json:"status" form:"status"`
}

type CreateCampusCategoryReq struct {
	Name      string `json:"name" binding:"required"`
	ParentID  *uint  `json:"parentId"`
	SortOrder int    `json:"sortOrder"`
	Icon      string `json:"icon"`
	Status    int    `json:"status"`
}

type UpdateCampusCategoryReq struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	ParentID  *uint  `json:"parentId"`
	SortOrder int    `json:"sortOrder"`
	Icon      string `json:"icon"`
	Status    int    `json:"status"`
}

type UpdateCampusCategoryStatusReq struct {
	ID     uint `json:"id" binding:"required"`
	Status *int `json:"status" binding:"required"`
}
