package interfaces

import "go.mongodb.org/mongo-driver/bson/primitive"

type Spider interface {
	Model
	GetName() (n string)
	GetType() (ty string)
	GetMode() (mode string)
	SetMode(mode string)
	GetNodeIds() (ids []primitive.ObjectID)
	SetNodeIds(ids []primitive.ObjectID)
	GetNodeTags() (tags []string)
	SetNodeTags(tags []string)
	GetCmd() (cmd string)
	SetCmd(cmd string)
	GetParam() (param string)
	SetParam(param string)
	GetPriority() (p int)
	SetPriority(p int)
}
