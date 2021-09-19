package models

import (
	"github.com/crawlab-team/crawlab-core/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Plugin struct {
	Id            primitive.ObjectID         `json:"_id" bson:"_id"`
	Name          string                     `json:"name" bson:"name"`
	Description   string                     `json:"description" bson:"description"`
	Type          string                     `json:"type" bson:"type"`
	Proto         string                     `json:"proto" bson:"proto"`
	Active        bool                       `json:"active" bson:"active"`
	Endpoint      string                     `json:"endpoint" bson:"endpoint"`
	Cmd           string                     `json:"cmd" bson:"cmd"`
	EventKey      entity.PluginEventKey      `json:"event_key" bson:"event_key"`
	UIComponents  []entity.PluginUIComponent `json:"ui_components" bson:"ui_components"`
	UISidebarNavs []entity.PluginUINav       `json:"ui_sidebar_navs" bson:"ui_sidebar_navs"`
	UIAssets      []entity.PluginUIAsset     `json:"ui_assets" bson:"ui_assets"`
}

func (p *Plugin) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Plugin) SetId(id primitive.ObjectID) {
	p.Id = id
}

func (p *Plugin) GetName() (name string) {
	return p.Name
}

func (p *Plugin) SetName(name string) {
	p.Name = name
}
