package models

type ModelMap struct {
	Node     Node
	Project  Project
	Spider   Spider
	Task     Task
	Job      Job
	Schedule Schedule
	User     User
	Setting  Setting
	Token    Token
	Variable Variable
}

type ModelListMap struct {
	Nodes     []Node
	Projects  []Project
	Spiders   []Spider
	Tasks     []Task
	Jobs      []Job
	Schedules []Schedule
	Users     []User
	Settings  []Setting
	Tokens    []Token
	Variables []Variable
}

func NewModelMap() (m *ModelMap) {
	return &ModelMap{}
}

func NewModelListMap() (m *ModelListMap) {
	return &ModelListMap{}
}
