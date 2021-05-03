package entity

type NodeInfo struct {
	Key         string `json:"key"`
	IsMaster    bool   `json:"is_master"`
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	Mac         string `json:"mac"`
	Hostname    string `json:"hostname"`
	Description string `json:"description"`
}

func (n NodeInfo) Value() interface{} {
	return n
}
