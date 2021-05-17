package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task interface {
	Model
	GetNodeId() (id primitive.ObjectID)
	GetNodeIds() (ids []primitive.ObjectID)
	GetNodeTags() (tags []string)
	GetStatus() (status string)
	SetStatus(status string)
	GetSpiderId() (id primitive.ObjectID)
	GetType() (ty string)
	GetCmd() (cmd string)
	GetParam() (param string)
}
