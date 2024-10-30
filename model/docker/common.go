package model

type Pagination struct {
	PageSize int `json:"pageSize" form:"pageSize"`
	Page     int `json:"page" form:"page"`
	Total    int `json:"total" form:"total"`
}
