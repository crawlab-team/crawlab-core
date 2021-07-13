module github.com/crawlab-team/crawlab-core

go 1.15

replace (
	github.com/crawlab-team/crawlab-db => /Users/marvzhang/projects/crawlab-team/crawlab-db
	github.com/crawlab-team/crawlab-fs => /Users/marvzhang/projects/crawlab-team/crawlab-fs
	github.com/crawlab-team/crawlab-grpc => /Users/marvzhang/projects/crawlab-team/crawlab-grpc/dist/go
	github.com/crawlab-team/crawlab-log => /Users/marvzhang/projects/crawlab-team/crawlab-log
	github.com/crawlab-team/crawlab-vcs => /Users/marvzhang/projects/crawlab-team/crawlab-vcs
	github.com/crawlab-team/go-trace => /Users/marvzhang/projects/crawlab-team/go-trace
	github.com/crawlab-team/goseaweedfs => /Users/marvzhang/projects/crawlab-team/goseaweedfs
)

require (
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/crawlab-db v0.0.0
	github.com/crawlab-team/crawlab-fs v0.0.0
	github.com/crawlab-team/crawlab-grpc v0.0.0
	github.com/crawlab-team/crawlab-log v0.0.0
	github.com/crawlab-team/crawlab-vcs v0.0.0
	github.com/crawlab-team/go-trace v0.1.0
	github.com/crawlab-team/goseaweedfs v0.1.6
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emirpasic/gods v1.12.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gamexg/proxyclient v0.0.0-20210207161252-499908056324 // indirect
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/hashicorp/go-sockaddr v1.0.0
	github.com/imroc/req v0.3.0
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901
	github.com/matcornic/hermes v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olivere/elastic/v7 v7.0.15
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.4.0
	github.com/robfig/cron v1.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/shadowsocks/shadowsocks-go v0.0.0-20200409064450-3e585ff90601 // indirect
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/ztrue/tracerr v0.3.0
	go.mongodb.org/mongo-driver v1.4.5
	go.uber.org/atomic v1.6.0
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
