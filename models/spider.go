package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Env struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Spider struct {
	Id          primitive.ObjectID   `json:"_id" bson:"_id"`                   // 爬虫ID
	Name        string               `json:"name" bson:"name"`                 // 爬虫名称（唯一）
	DisplayName string               `json:"display_name" bson:"display_name"` // 爬虫显示名称
	Type        string               `json:"type" bson:"type"`                 // 爬虫类别
	Col         string               `json:"col" bson:"col"`                   // 结果储存位置
	Envs        []Env                `json:"envs" bson:"envs"`                 // 环境变量
	Description string               `json:"description" bson:"description"`   // 备注
	ProjectId   primitive.ObjectID   `json:"project_id" bson:"project_id"`     // Project.Id
	IsPublic    bool                 `json:"is_public" bson:"is_public"`       // 是否公开
	Mode        string               `json:"mode" bson:"mode"`                 // default Task.Mode
	NodeIds     []primitive.ObjectID `json:"node_ids" bson:"node_ids"`         // default Task.NodeIds
	NodeTags    []string             `json:"node_tags" bson:"node_tags"`       // default Task.NodeTags

	// 自定义爬虫
	Cmd   string `json:"cmd" bson:"cmd"`     // 执行命令
	Param string `json:"param" bson:"param"` // default task param

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

func (s *Spider) Add() (err error) {
	if s.Id.IsZero() {
		s.Id = primitive.NewObjectID()
	}
	m := NewDelegate(ModelColNameSpider, s)
	return m.Add()
}

func (s *Spider) Save() (err error) {
	m := NewDelegate(ModelColNameSpider, s)
	return m.Save()
}

func (s *Spider) Delete() (err error) {
	m := NewDelegate(ModelColNameSpider, s)
	return m.Delete()
}

func (s *Spider) GetArtifact() (a Artifact, err error) {
	d := NewDelegate(ModelColNameSpider, s)
	return d.GetArtifact()
}

func (s *Spider) GetId() (id primitive.ObjectID) {
	return s.Id
}
