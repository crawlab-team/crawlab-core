package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Env struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Spider struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`                 // spider id
	Name        string               `json:"name" bson:"name"`               // spider name
	Type        string               `json:"type" bson:"type"`               // spider type
	ColId       primitive.ObjectID   `json:"col_id" bson:"col_id"`           // data collection id
	ColName     string               `json:"col_name,omitempty" bson:"-"`    // data collection name
	Description string               `json:"description" bson:"description"` // description
	ProjectId   primitive.ObjectID   `json:"project_id" bson:"project_id"`   // Project.Id
	Mode        string               `json:"mode" bson:"mode"`               // default Task.Mode
	NodeIds     []primitive.ObjectID `json:"node_ids" bson:"node_ids"`       // default Task.NodeIds
	NodeTags    []string             `json:"node_tags" bson:"node_tags"`     // default Task.NodeTags
	Tags        []Tag                `json:"tags" bson:"-"`                  // tags
	Stat        *SpiderStat          `json:"stat,omitempty" bson:"-"`

	IsPublic bool  `json:"is_public" bson:"is_public"` // 是否公开
	Envs     []Env `json:"envs" bson:"envs"`           // 环境变量

	// 自定义爬虫
	Cmd      string `json:"cmd" bson:"cmd"`     // 执行命令
	Param    string `json:"param" bson:"param"` // default task param
	Priority int    `json:"priority" bson:"priority"`

	// Scrapy 爬虫（属于自定义爬虫）
	IsScrapy    bool     `json:"is_scrapy" bson:"is_scrapy"`       // 是否为 Scrapy 爬虫
	SpiderNames []string `json:"spider_names" bson:"spider_names"` // 爬虫名称列表

	// 可配置爬虫
	Template string `json:"template" bson:"template"` // Spiderfile模版

	// Git 设置
	IsGit            bool   `json:"is_git" bson:"is_git"`                         // 是否为 Git
	GitUrl           string `json:"git_url" bson:"git_url"`                       // Git URL
	GitBranch        string `json:"git_branch" bson:"git_branch"`                 // Git 分支
	GitHasCredential bool   `json:"git_has_credential" bson:"git_has_credential"` // Git 是否加密
	GitUsername      string `json:"git_username" bson:"git_username"`             // Git 用户名
	GitPassword      string `json:"git_password" bson:"git_password"`             // Git 密码
	GitAutoSync      bool   `json:"git_auto_sync" bson:"git_auto_sync"`           // Git 是否自动同步
	GitSyncFrequency string `json:"git_sync_frequency" bson:"git_sync_frequency"` // Git 同步频率
	GitSyncError     string `json:"git_sync_error" bson:"git_sync_error"`         // Git 同步错误

	// 长任务
	IsLongTask bool `json:"is_long_task" bson:"is_long_task"` // 是否为长任务

	// 去重
	IsDedup     bool   `json:"is_dedup" bson:"is_dedup"`         // 是否去重
	DedupField  string `json:"dedup_field" bson:"dedup_field"`   // 去重字段
	DedupMethod string `json:"dedup_method" bson:"dedup_method"` // 去重方式

	// Web Hook
	IsWebHook  bool   `json:"is_web_hook" bson:"is_web_hook"`   // 是否开启 Web Hook
	WebHookUrl string `json:"web_hook_url" bson:"web_hook_url"` // Web Hook URL
}

func (s *Spider) GetId() (id primitive.ObjectID) {
	return s.Id
}

func (s *Spider) SetId(id primitive.ObjectID) {
	s.Id = id
}

func (s *Spider) GetTags() (tags []interfaces.Tag) {
	return convertTagsToInterfaces(s.Tags)
}

func (s *Spider) SetTags(tags []interfaces.Tag) {
	s.Tags = convertInterfacesToTags(tags)
}

func (s *Spider) GetName() (n string) {
	return s.Name
}

func (s *Spider) GetType() (ty string) {
	return s.Type
}

func (s *Spider) GetMode() (mode string) {
	return s.Mode
}

func (s *Spider) SetMode(mode string) {
	s.Mode = mode
}

func (s *Spider) GetNodeIds() (ids []primitive.ObjectID) {
	return s.NodeIds
}

func (s *Spider) SetNodeIds(ids []primitive.ObjectID) {
	s.NodeIds = ids
}

func (s *Spider) GetNodeTags() (tags []string) {
	return s.NodeTags
}

func (s *Spider) SetNodeTags(tags []string) {
	s.NodeTags = tags
}

func (s *Spider) GetCmd() (cmd string) {
	return s.Cmd
}

func (s *Spider) SetCmd(cmd string) {
	s.Cmd = cmd
}

func (s *Spider) GetParam() (param string) {
	return s.Param
}

func (s *Spider) SetParam(param string) {
	s.Param = param
}

func (s *Spider) GetPriority() (p int) {
	return s.Priority
}

func (s *Spider) SetPriority(p int) {
	s.Priority = p
}
