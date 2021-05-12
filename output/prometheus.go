package output

import (
	"fmt"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fastjson"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
	Outputs["prometheus"] = PrometheusFactory
}

type PrometheusConfig struct {
	Salt   string
	Listen string
	Label  string
}

type TraefikProm struct {
	registry *prometheus.Registry
	hits     *prometheus.CounterVec
}

func NewTraefikProm() *TraefikProm {
	t := &TraefikProm{
		registry: prometheus.NewRegistry(),
		hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_hit",
			Help: "HTTP hit, group by status",
		}, []string{"status"}),
	}
	t.registry.MustRegister(t.hits)

	return t
}

type PrometheusOutput struct {
	gatherers map[string]*TraefikProm
	config    *PrometheusConfig
}

func PrometheusFactory(rawCfg map[string]interface{}) (api.Output, error) {
	var cfg PrometheusConfig
	err := mapstructure.Decode(rawCfg, &cfg)

	if err != nil {
		return nil, err
	}
	return &PrometheusOutput{
		gatherers: make(map[string]*TraefikProm),
		config:    &cfg,
	}, nil
}

func (p *PrometheusOutput) Close() error {
	return nil
}

func (p *PrometheusOutput) Write(ts time.Time, line string, meta map[string]interface{}) error {
	keyRaw, ok := meta[p.config.Label]
	if !ok { // it's an anonymous log
		return nil
	}
	key, ok := keyRaw.(string)
	if !ok {
		return fmt.Errorf("the label is not a string %v %T", keyRaw, keyRaw)
	}

	_, ok = p.gatherers[key]
	if !ok {
		p.gatherers[key] = NewTraefikProm()
	}
	status := fastjson.GetInt([]byte(line), "OriginStatus")
	p.gatherers[key].hits.WithLabelValues(fmt.Sprintf("%dxx", status/100)).Inc()
	return nil
}
