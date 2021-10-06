package models

type ModelMap struct {
	Artifact       Artifact
	Tag            Tag
	Node           Node
	Project        Project
	Spider         Spider
	Task           Task
	Job            Job
	Schedule       Schedule
	User           User
	Setting        Setting
	Token          Token
	Variable       Variable
	TaskQueueItem  TaskQueueItem
	TaskStat       TaskStat
	Plugin         Plugin
	SpiderStat     SpiderStat
	DataSource     DataSource
	DataCollection DataCollection
	Result         Result
	Password       Password
	ExtraValue     ExtraValue
	PluginStatus   PluginStatus
}

type ModelListMap struct {
	Artifacts       []Artifact
	Tags            []Tag
	Nodes           []Node
	Projects        []Project
	Spiders         []Spider
	Tasks           []Task
	Jobs            []Job
	Schedules       []Schedule
	Users           []User
	Settings        []Setting
	Tokens          []Token
	Variables       []Variable
	TaskQueueItems  []TaskQueueItem
	TaskStats       []TaskStat
	Plugins         []Plugin
	SpiderStats     []SpiderStat
	DataSources     []DataSource
	DataCollections []DataCollection
	Results         []Result
	Passwords       []Password
	ExtraValues     []ExtraValue
	PluginStatus    []PluginStatus
}

func NewModelMap() (m *ModelMap) {
	return &ModelMap{}
}

func NewModelListMap() (m *ModelListMap) {
	return &ModelListMap{
		Artifacts:       []Artifact{},
		Tags:            []Tag{},
		Nodes:           []Node{},
		Projects:        []Project{},
		Spiders:         []Spider{},
		Tasks:           []Task{},
		Jobs:            []Job{},
		Schedules:       []Schedule{},
		Users:           []User{},
		Settings:        []Setting{},
		Tokens:          []Token{},
		Variables:       []Variable{},
		TaskQueueItems:  []TaskQueueItem{},
		TaskStats:       []TaskStat{},
		Plugins:         []Plugin{},
		SpiderStats:     []SpiderStat{},
		DataSources:     []DataSource{},
		DataCollections: []DataCollection{},
		Results:         []Result{},
		Passwords:       []Password{},
		ExtraValues:     []ExtraValue{},
		PluginStatus:    []PluginStatus{},
	}
}
