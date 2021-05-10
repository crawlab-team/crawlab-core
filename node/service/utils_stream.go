package service

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-grpc"
)

func GetStreamMessageWithError(code grpc.StreamMessageCode, err error) (res *grpc.StreamMessage) {
	return &grpc.StreamMessage{
		Code:  code,
		Error: err.Error(),
	}
}

func GetStreamMessage(code grpc.StreamMessageCode) (res *grpc.StreamMessage) {
	return &grpc.StreamMessage{
		Code: code,
	}
}

func GetStreamMessageWithData(code grpc.StreamMessageCode, data interface{}) (res *grpc.StreamMessage) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return GetStreamMessageWithError(code, err)
	}
	return &grpc.StreamMessage{
		Data: bytes,
	}
}
