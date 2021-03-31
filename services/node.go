package services

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"os"
	"path"
)

type NodeServiceInterface interface {
	Init() (err error)
	GetCurrentNode() (n *models.Node, err error)
	GetAllNodeIds() (ids []primitive.ObjectID, err error)
}

type NodeServiceOptions struct {
	Master       bool
	RegisterType string
	DataPath     string
	Ip           string
	Hostname     string
	MacAddress   string
}

func NewNodeService(opts *NodeServiceOptions) (svc *nodeService, err error) {
	svc = &nodeService{
		opts: opts,
	}
	if opts.DataPath == "" {
		homePath := os.Getenv("HOME")
		if homePath == "" {
			homePath = "/"
		}
		svc.dataPath = path.Join(homePath, ".crawlab")
	}
	if err := svc.Init(); err != nil {
		return nil, err
	}
	return svc, nil
}

func InitNodeService() (err error) {
	NodeService, err = NewNodeService(&NodeServiceOptions{
		Master:       viper.GetBool("server.master"),
		DataPath:     viper.GetString("server.register.dataPath"),
		RegisterType: viper.GetString("server.register.type"),
		Ip:           viper.GetString("server.register.ip"),
		Hostname:     viper.GetString("server.register.hostname"),
		MacAddress:   viper.GetString("server.register.mac"),
	})
	return err
}

type nodeService struct {
	opts     *NodeServiceOptions
	dataPath string
	data     entity.NodeData
}

func (svc *nodeService) Init() (err error) {
	// create data directory if not exists
	if _, err := os.Stat(svc.dataPath); err != nil {
		if err := os.MkdirAll(svc.dataPath, os.ModePerm); err != nil {
			return err
		}
	}

	// create node data file if not exists
	var nodeData entity.NodeData
	nodeFilePath := path.Join(svc.dataPath, "node.json")
	if _, err := os.Stat(nodeFilePath); err != nil {
		nodeData = entity.NodeData{
			Ip:         svc.opts.Ip,
			Hostname:   svc.opts.Hostname,
			MacAddress: svc.opts.MacAddress,
			UUID:       uuid.New().String(),
		}
		nodeData.Hash = svc.getHash(nodeData.UUID)
		nodeData.Key = svc.getKey(nodeData)
		data, err := json.Marshal(&nodeData)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(nodeFilePath, data, os.ModePerm); err != nil {
			return err
		}
	} else {
		data, err := ioutil.ReadFile(nodeFilePath)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &nodeData); err != nil {
			return err
		}
	}
	svc.data = nodeData

	return nil
}

func (svc *nodeService) GetCurrentNode() (n *models.Node, err error) {
	node, err := models.NodeService.Get(bson.M{"key": svc.data.Key}, nil)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (svc *nodeService) GetAllNodeIds() (ids []primitive.ObjectID, err error) {
	nodes, err := models.NodeService.GetList(bson.M{"enabled": true, "active": true}, nil)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		ids = append(ids, n.Id)
	}
	return ids, nil
}

func (svc *nodeService) getHash(sum string) (md5sum string) {
	h := md5.New()
	if svc.opts.Ip != "" {
		h.Write([]byte(svc.opts.Ip))
	}
	if svc.opts.Hostname != "" {
		h.Write([]byte(svc.opts.Hostname))
	}
	if svc.opts.MacAddress != "" {
		h.Write([]byte(svc.opts.MacAddress))
	}
	if sum != "" {
		h.Write([]byte(sum))
	}
	md5sum = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return md5sum
}

func (svc *nodeService) getKey(nodeData entity.NodeData) (key string) {
	switch svc.opts.RegisterType {
	case constants.RegisterTypeIp:
		return svc.opts.Ip
	case constants.RegisterTypeHostname:
		return svc.opts.Hostname
	case constants.RegisterTypeMac:
		return svc.opts.MacAddress
	default:
		return nodeData.Hash
	}
}

var NodeService *nodeService
