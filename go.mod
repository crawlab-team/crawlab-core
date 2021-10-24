module github.com/crawlab-team/crawlab-core

go 1.15

replace (
	github.com/crawlab-team/crawlab-vcs => /Users/marvzhang/projects/crawlab-team/crawlab-vcs
)

require (
	github.com/StackExchange/wmi v1.2.0 // indirect
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/crawlab-db v0.1.1
	github.com/crawlab-team/crawlab-fs v0.1.0
	github.com/crawlab-team/crawlab-grpc v0.6.0-beta.20211009.1455
	github.com/crawlab-team/crawlab-log v0.1.0
	github.com/crawlab-team/crawlab-vcs v0.1.0
	github.com/crawlab-team/go-trace v0.1.0
	github.com/crawlab-team/goseaweedfs v0.6.0-beta.20210725.1917
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emirpasic/gods v1.12.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/imroc/req v0.3.0
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olivere/elastic/v7 v7.0.15
	github.com/prometheus/common v0.4.0
	github.com/robfig/cron/v3 v3.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/thoas/go-funk v0.9.1
	github.com/ztrue/tracerr v0.3.0
	go.mongodb.org/mongo-driver v1.4.5
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/net v0.0.0-20210928044308-7d9f5e0b762b // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/grpc v1.34.0
)
