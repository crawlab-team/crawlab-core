package entity

type NodeInfo struct {
	Key      string `json:"key"`
	IsMaster bool   `json:"is_master"`
}

func (n NodeInfo) Value() interface{} {
	return n
}
