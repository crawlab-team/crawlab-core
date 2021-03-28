package entity

type Pagination struct {
	Page int `form:"page" url:"page"`
	Size int `form:"size" url:"size"`
}
