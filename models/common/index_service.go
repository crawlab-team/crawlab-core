package common

import (
	"github.com/crawlab-team/crawlab-core/interfaces"
	"github.com/crawlab-team/crawlab-db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func CreateIndexes() {
	// artifacts
	mongo.GetMongoCol(interfaces.ModelColNameArtifact).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"_col": 1}},
		{Keys: bson.M{"_del": 1}},
		{Keys: bson.M{"_tid": 1}},
	})

	// tags
	mongo.GetMongoCol(interfaces.ModelColNameTag).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"col": 1}},
		{Keys: bson.M{"name": 1}},
	})

	// nodes
	mongo.GetMongoCol(interfaces.ModelColNameNode).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"k": 1}},  // key
		{Keys: bson.M{"n": 1}},  // name
		{Keys: bson.M{"im": 1}}, // is_master
		{Keys: bson.M{"s": 1}},  // status
		{Keys: bson.M{"en": 1}}, // enabled
		{Keys: bson.M{"a": 1}},  // active
	})

	// projects
	mongo.GetMongoCol(interfaces.ModelColNameNode).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// spiders
	mongo.GetMongoCol(interfaces.ModelColNameSpider).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"col_id": 1}},
		{Keys: bson.M{"project_id": 1}},
	})

	// tasks
	mongo.GetMongoCol(interfaces.ModelColNameTask).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"status": 1}},
		{Keys: bson.M{"node_id": 1}},
		{Keys: bson.M{"schedule_id": 1}},
		{Keys: bson.M{"type": 1}},
		{Keys: bson.M{"mode": 1}},
		{Keys: bson.M{"priority": 1}},
		{Keys: bson.M{"parent_id": 1}},
		{Keys: bson.M{"has_sub": 1}},
	})

	// schedules
	mongo.GetMongoCol(interfaces.ModelColNameSchedule).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
		{Keys: bson.M{"spider_id": 1}},
		{Keys: bson.M{"enabled": 1}},
	})

	// users
	mongo.GetMongoCol(interfaces.ModelColNameUser).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"username": 1}},
		{Keys: bson.M{"role": 1}},
		{Keys: bson.M{"email": 1}},
	})

	// settings
	mongo.GetMongoCol(interfaces.ModelColNameSetting).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},
	})

	// tokens
	mongo.GetMongoCol(interfaces.ModelColNameToken).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// variables
	mongo.GetMongoCol(interfaces.ModelColNameVariable).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"key": 1}},
	})

	// plugins
	mongo.GetMongoCol(interfaces.ModelColNamePlugin).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// data sources
	mongo.GetMongoCol(interfaces.ModelColNameDataSource).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// data collections
	mongo.GetMongoCol(interfaces.ModelColNameDataCollection).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"name": 1}},
	})

	// extra values
	mongo.GetMongoCol(interfaces.ModelColNameExtraValues).MustCreateIndexes([]mongo2.IndexModel{
		{Keys: bson.M{"oid": 1}},
		{Keys: bson.M{"m": 1}},
		{Keys: bson.M{"t": 1}},
		{Keys: bson.M{"m": 1, "t": 1}},
		{Keys: bson.M{"oid": 1, "m": 1, "t": 1}},
	})
}
