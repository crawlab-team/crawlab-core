package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crawlab-team/crawlab-core/constants"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/crawlab-team/crawlab-core/model"
	"github.com/crawlab-team/crawlab-core/utils"
	db "github.com/crawlab-team/crawlab-db"
	pb "github.com/crawlab-team/crawlab-grpc"
	"github.com/globalsign/mgo/bson"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"testing"
	"time"
)

type TaskTestObject struct {
	host     string
	port     string
	s        *CrawlabGrpcService
	c        pb.TaskServiceClient
	conn     *grpc.ClientConn
	spiderId bson.ObjectId
	spider   model.Spider
	taskId   string
	task     model.Task
	nodeId   bson.ObjectId
	node     model.Node
}

func setupTask() (to *TaskTestObject, err error) {
	// grpc address
	host := "0.0.0.0"
	port := "9999"
	address := fmt.Sprintf("%s:%s", host, port)

	// grpc service
	s, err := NewCrawlabGrpcService()
	if err != nil {
		return nil, err
	}

	// test object
	to = &TaskTestObject{
		host:     host,
		port:     port,
		s:        s,
		spiderId: bson.NewObjectId(),
		taskId:   uuid.New().String(),
		nodeId:   bson.NewObjectId(),
	}

	// set debug
	viper.Set("debug", true)

	// set grpc address
	viper.Set("grpc.host", to.host)
	viper.Set("grpc.port", to.port)

	// init
	go to.s.Init()

	// dial
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	to.conn = conn
	to.c = pb.NewTaskServiceClient(conn)

	// mongo
	viper.Set("mongo.host", "localhost")
	viper.Set("mongo.port", "27017")
	viper.Set("mongo.db", "test")
	if err := db.InitMongo(); err != nil {
		return nil, err
	}

	// spider
	to.spider = model.Spider{
		Id:   to.spiderId,
		Name: "test_spider",
		Type: constants.Customized,
		Cmd:  "python main.py",
		Envs: []model.Env{
			{Name: "Env1", Value: "Value1"},
			{Name: "Env2", Value: "Value2"},
		},
		Col:       "results_test_spider",
		FileId:    bson.ObjectIdHex(constants.ObjectIdNull),
		ProjectId: bson.ObjectIdHex(constants.ObjectIdNull),
		UserId:    bson.ObjectIdHex(constants.ObjectIdNull),
	}
	if err := to.spider.Add(); err != nil {
		return nil, err
	}

	// task
	to.task = model.Task{
		Id:         to.taskId,
		SpiderId:   to.spiderId,
		NodeId:     to.nodeId,
		Type:       constants.TaskTypeSpider,
		ScheduleId: bson.ObjectIdHex(constants.ObjectIdNull),
		UserId:     bson.ObjectIdHex(constants.ObjectIdNull),
	}
	if err := model.AddTask(to.task); err != nil {
		return nil, err
	}

	// node
	to.node = model.Node{
		Id:       to.nodeId,
		Name:     "node name",
		IsMaster: true,
	}
	if err := to.node.Add(); err != nil {
		return nil, err
	}

	return to, nil
}

func cleanupTask(to *TaskTestObject) {
	_ = model.RemoveSpider(to.spiderId)
	_ = model.RemoveTask(to.taskId)
	_ = to.node.Delete()
	to.s.server.Stop()
	_ = to.conn.Close()
	s, col := db.GetCol(to.spider.Col)
	defer s.Close()
	_ = col.DropCollection()
}

func TestTaskService_GetTaskInfo(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// test GetTaskInfo
	req := &pb.TaskServiceRequest{
		TaskId: to.taskId,
	}
	res, err := to.c.GetTaskInfo(context.Background(), req)
	require.Nil(t, err)
	require.NotNil(t, res.Task)
	require.NotNil(t, res.Spider)
	require.NotNil(t, res.Node)
	require.Equal(t, to.taskId, res.Task.XId)
	require.Equal(t, to.spiderId.Hex(), res.Spider.XId)
	require.Equal(t, to.nodeId.Hex(), res.Node.XId)

	cleanupTask(to)
}

func TestTaskService_SaveItem(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// test SaveItem
	key1 := "f1"
	key2 := "f2"
	value1 := "v1"
	value2 := "v2"
	item := map[string]string{
		key1: value1,
		key2: value2,
	}
	data, err := json.Marshal(&item)
	require.Nil(t, err)
	req := &pb.TaskServiceRequest{
		TaskId: to.taskId,
		Data:   data,
	}
	_, err = to.c.SaveItem(context.Background(), req)
	require.Nil(t, err)

	// wait for a period
	time.Sleep(2 * time.Second)

	// test saved item
	s, col := db.GetCol(to.spider.Col)
	defer s.Close()
	var dbItem entity.ResultItem
	err = col.Find(nil).One(&dbItem)
	require.Nil(t, err)
	require.Equal(t, item[key1], dbItem[key1])
	require.Equal(t, item[key2], dbItem[key2])

	cleanupTask(to)
}

func TestTaskService_SaveItems(t *testing.T) {
	to, err := setupTask()
	require.Nil(t, err)

	// test SaveItems
	batch := 3
	size := 100
	var items []entity.ResultItem
	for i := 0; i < batch; i++ {
		for j := 0; j < size; j++ {
			items = append(items, entity.ResultItem{
				"num": i*size + j,
			})
		}
		data, err := json.Marshal(&items)
		require.Nil(t, err)
		req := &pb.TaskServiceRequest{
			TaskId: to.taskId,
			Data:   data,
		}
		_, err = to.c.SaveItems(context.Background(), req)
		require.Nil(t, err)

		// reset
		items = []entity.ResultItem{}

		// wait
		time.Sleep(1 * time.Second)
	}

	// wait for a period
	time.Sleep(3 * time.Second)

	// test saved item
	s, col := db.GetCol(to.spider.Col)
	defer s.Close()
	count, err := col.Count()
	utils.LogDebug(fmt.Sprintf("count: %d", count))
	require.Nil(t, err)
	require.Equal(t, batch*size, count)
	var dbItems []entity.ResultItem
	err = col.Find(nil).Sort("+num").All(&dbItems)
	require.Equal(t, batch*size, len(dbItems))
	for i, item := range dbItems {
		require.IsType(t, float64(i), item["num"])
		require.Equal(t, i, int(item["num"].(float64)))
	}

	cleanupTask(to)
}
