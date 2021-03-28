package entity

type Condition struct {
	Key   string      `json:"key"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

type Filter struct {
	IsOr       bool        `form:"is_or" url:"is_or"`
	Conditions []Condition `json:"conditions"`
}
