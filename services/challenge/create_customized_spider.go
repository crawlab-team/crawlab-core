package challenge

import (
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/globalsign/mgo/bson"
)

type CreateCustomizedSpiderService struct {
	UserId bson.ObjectId
}

func (s *CreateCustomizedSpiderService) Check() (bool, error) {
	query := bson.M{
		"user_id": s.UserId,
		"type":    constants.Customized,
	}
	_, count, err := model.GetSpiderList(query, 0, 1, "-_id")
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
