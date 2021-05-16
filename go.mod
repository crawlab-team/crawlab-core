module github.com/crawlab-team/crawlab-core

go 1.15

replace (
	github.com/crawlab-team/crawlab-db => /Users/marvzhang/projects/crawlab-team/crawlab-db
	github.com/crawlab-team/crawlab-fs => /Users/marvzhang/projects/crawlab-team/crawlab-fs
	github.com/crawlab-team/crawlab-grpc => /Users/marvzhang/projects/crawlab-team/crawlab-grpc/dist/go
	github.com/crawlab-team/crawlab-log => /Users/marvzhang/projects/crawlab-team/crawlab-log
	github.com/crawlab-team/crawlab-vcs => /Users/marvzhang/projects/crawlab-team/crawlab-vcs
	github.com/crawlab-team/go-trace => /Users/marvzhang/projects/crawlab-team/go-trace
)

require (
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/crawlab-db v0.0.0
	github.com/crawlab-team/crawlab-fs v0.0.0
	github.com/crawlab-team/crawlab-grpc v0.0.0
	github.com/crawlab-team/crawlab-log v0.0.0
	github.com/crawlab-team/crawlab-vcs v0.0.0
	github.com/crawlab-team/go-trace v0.1.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emirpasic/gods v1.12.0
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gavv/httpexpect/v2 v2.2.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.3.0
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/uuid v1.1.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/hashicorp/go-sockaddr v1.0.0
	github.com/imroc/req v0.3.0
	github.com/matcornic/hermes v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/olivere/elastic/v7 v7.0.15
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v3.20.11+incompatible
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/ztrue/tracerr v0.3.0
	go.mongodb.org/mongo-driver v1.4.5
	go.uber.org/atomic v1.6.0
	go.uber.org/dig v1.10.0
	google.golang.org/grpc v1.34.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
