package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Tags        []Tag              `json:"tags" bson:"-"`
	Spiders     int                `json:"spiders" bson:"-"`
}

func (p *Project) GetId() (id primitive.ObjectID) {
	return p.Id
}

func (p *Project) SetId(id primitive.ObjectID) {
	p.Id = id
}

func (p *Project) GetTags() (tags []interfaces.Tag) {
	return convertTagsToInterfaces(p.Tags)
}

func (p *Project) SetTags(tags []interfaces.Tag) {
	p.Tags = convertInterfacesToTags(tags)
}
