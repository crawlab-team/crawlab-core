package models

import (
	"github.com/crawlab-team/crawlab-core/errors"
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/store"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/go-trace"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func convertTypeTag(d interface{}, err error) (res *Tag, err2 error) {
	if err != nil {
		return nil, err
	}
	res, ok := d.(*Tag)
	if !ok {
		return nil, trace.TraceError(errors.ErrorModelInvalidType)
	}
	return res, nil
}

func (svc *Service) GetTagById(id primitive.ObjectID) (res *Tag, err error) {
	d, err := MustGetService(interfaces.ModelIdTag).GetById(id)
	return convertTypeTag(d, err)
}

func (svc *Service) GetTag(query bson.M, opts *mongo.FindOptions) (res *Tag, err error) {
	d, err := MustGetService(interfaces.ModelIdTag).Get(query, opts)
	return convertTypeTag(d, err)
}

func (svc *Service) GetTagList(query bson.M, opts *mongo.FindOptions) (res []Tag, err error) {
	err = getListSerializeTarget(interfaces.ModelIdTag, query, opts, &res)
	return res, trace.TraceError(err)
}

func (svc *Service) getTagIds(colName string, tags []Tag) (tagIds []primitive.ObjectID, err error) {
	// iterate tag names
	for _, tag := range tags {
		// count of tags with the name
		tagDb, err := svc.GetTag(bson.M{"name": tag.Name, "col": colName}, nil)
		if err == nil {
			// tag exists
			tag = *tagDb
		} else if err == mongo2.ErrNoDocuments {
			// add new tag if not exists
			colorHex := tag.Color
			if colorHex == "" {
				color, _ := store.ColorService.GetRandom()
				colorHex = color.GetHex()
			}
			tag = Tag{
				Name:  tag.Name,
				Color: colorHex,
				Col:   colName,
			}
			if err := tag.Add(); err != nil {
				return tagIds, trace.TraceError(err)
			}
		}

		// add to tag ids
		tagIds = append(tagIds, tag.Id)
	}

	return tagIds, nil
}

func (svc *Service) UpdateTagsById(colName string, id primitive.ObjectID, tags []Tag) (tagIds []primitive.ObjectID, err error) {
	// get tag ids to update
	tagIds, err = svc.getTagIds(colName, tags)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}

	// update in db
	a, err := svc.GetArtifactById(id)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}
	a.TagIds = tagIds
	if err := mongo.GetMongoCol(interfaces.ModelColNameArtifact).ReplaceId(id, a); err != nil {
		return tagIds, err
	}
	return tagIds, nil
}

func (svc *Service) UpdateTags(colName string, query bson.M, tags []Tag) (tagIds []primitive.ObjectID, err error) {
	// tag ids to update
	tagIds, err = svc.getTagIds(colName, tags)
	if err != nil {
		return tagIds, trace.TraceError(err)
	}

	// update
	update := bson.M{
		"_tid": tagIds,
	}

	// fields
	fields := []string{"_tid"}

	// update in db
	if err := svc.Update(query, update, fields); err != nil {
		return tagIds, trace.TraceError(err)
	}

	return tagIds, nil
}
