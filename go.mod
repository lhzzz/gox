module singer.com

go 1.16

require (
	github.com/armon/go-metrics v0.4.1
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/gin-gonic/gin v1.8.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/grpc-opentracing v0.0.0-20180507213350-8e809c8a8645
	github.com/juju/ratelimit v1.0.2
	github.com/nsqio/go-nsq v1.1.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.13.0
	github.com/sirupsen/logrus v1.9.0
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/cast v1.5.0
	github.com/spf13/viper v1.13.0
	github.com/stretchr/testify v1.8.0
	github.com/uber/jaeger-client-go v2.30.0+incompatible
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
	gorm.io/driver/mysql v1.4.3
	gorm.io/gorm v1.24.0
	gorm.io/plugin/opentracing v0.0.0-20211220013347-7d2b2af23560
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	gorm.io/plugin/dbresolver v1.3.0
)
