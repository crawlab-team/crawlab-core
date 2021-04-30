package models

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-core/utils/binders"
)

func getModelColName(id interfaces.ModelId) (colName string) {
	return binders.NewColNameBinder(id).MustBindString()
}
