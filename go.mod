module github.com/factorysh/traefik-log-multiplexer

go 1.16

require (
	github.com/docker/docker v20.10.11+incompatible
	github.com/factorysh/docker-visitor v1.7.4
	github.com/fluent/fluent-logger-golang v1.7.0
	github.com/getsentry/sentry-go v0.11.0
	github.com/grafana/loki v1.6.2-0.20211021114919-0ae0d4da122d
	github.com/imdario/mergo v0.3.12
	github.com/influxdata/tail v1.0.1-0.20200707181643-03a791b270e4
	github.com/mitchellh/mapstructure v1.4.2
	github.com/onrik/logrus v0.9.0
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/client_model v0.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tinylib/msgp v1.1.6
	github.com/valyala/fastjson v1.6.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

// Replace directives from Loki
replace (
	github.com/Azure/azure-sdk-for-go => github.com/Azure/azure-sdk-for-go v36.2.0+incompatible
	github.com/hashicorp/consul => github.com/hashicorp/consul v1.10.1
	github.com/hpcloud/tail => github.com/grafana/tail v0.0.0-20201004203643-7aa4e4a91f03
	github.com/prometheus/prometheus => github.com/grafana/prometheus v1.8.2-0.20211103031328-89bb32ee4ae7
	gopkg.in/yaml.v2 => github.com/rfratto/go-yaml v0.0.0-20200521142311-984fc90c8a04
	k8s.io/api => k8s.io/api v0.21.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.0
	k8s.io/client-go => k8s.io/client-go v0.21.0
)
