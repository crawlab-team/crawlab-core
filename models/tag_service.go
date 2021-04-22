package models

import (
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagServiceInterface interface {
	getTagIds(tagNames []string) (tagIds []primitive.ObjectID, err error)
	GetModelById(id primitive.ObjectID) (res Tag, err error)
	GetModel(query bson.M, opts *mongo.FindOptions) (res Tag, err error)
	GetModelList(query bson.M, opts *mongo.FindOptions) (res []Tag, err error)
	UpdateTagsById(id primitive.ObjectID, tagNames []string) (err error)
	UpdateTags(query bson.M, tagNames []string) (err error)
}

type tagService struct {
	*CommonService
}

func (svc *tagService) getTagIds(tagNames []string) (tagIds []primitive.ObjectID, err error) {
	// iterate tag names
	for _, tagName := range tagNames {
		count, err := TagService.Count(bson.M{"name": tagName, "mid": svc.id})
		if err != nil {
			return tagIds, err
		}

		// declare tag
		var tag Tag

		// add new tag if not exists
		if count == 0 {
			color, _ := ColorService.GetRandom()
			tag = Tag{
				Name:  tagName,
				Color: color.Hex,
				Col:   getModelColName(svc.id),
			}
			if err := tag.Add(); err != nil {
				return tagIds, err
			}
		} else {
			tag, err = TagService.GetModel(bson.M{"name": tagName}, nil)
			if err != nil {
				return tagIds, err
			}
		}

		// add to tag ids
		tagIds = append(tagIds, tag.Id)
	}

	return tagIds, nil
}

func (svc *tagService) GetModelById(id primitive.ObjectID) (res Tag, err error) {
	err = svc.findId(id).One(&res)
	return res, err
}

func (svc *tagService) GetModel(query bson.M, opts *mongo.FindOptions) (res Tag, err error) {
	err = svc.find(query, opts).One(&res)
	return res, err
}

func (svc *tagService) GetModelList(query bson.M, opts *mongo.FindOptions) (res []Tag, err error) {
	err = svc.find(query, opts).All(&res)
	return res, err
}

func (svc *tagService) UpdateTagsById(id primitive.ObjectID, tagNames []string) (err error) {
	// tag ids to update
	tagIds, err := svc.getTagIds(tagNames)
	if err != nil {
		return err
	}

	// update in db
	if err := ArtifactService.UpdateById(id, bson.M{"_tid": tagIds}); err != nil {
		return err
	}

	return nil
}

func (svc *tagService) UpdateTags(query bson.M, tagNames []string) (err error) {
	// tag ids to update
	tagIds, err := svc.getTagIds(tagNames)
	if err != nil {
		return err
	}

	// update
	update := bson.M{"_tid": tagIds}

	// fields
	fields := []string{"_tid"}

	// update in db
	if err := ArtifactService.Update(query, update, fields); err != nil {
		return err
	}

	return nil
}

func NewTagService() (svc *tagService) {
	return &tagService{NewCommonService(ModelIdTag)}
}

var TagService *tagService
