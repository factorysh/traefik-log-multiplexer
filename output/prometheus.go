package output

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/mitchellh/mapstructure"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
	Outputs["prometheus"] = PrometheusFactory
}

type PrometheusConfig struct {
	Salt  string
	Addr  string
	Label string
}

type TraefikProm struct {
	registry  *prometheus.Registry
	hits      *prometheus.CounterVec
	latencies *prometheus.HistogramVec
	mac       []byte
}

func NewTraefikProm(project, key []byte) *TraefikProm {
	t := &TraefikProm{
		registry: prometheus.NewRegistry(),
		hits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "http_hit",
			Help: "HTTP hit, group by status",
		}, []string{"status"}),
		latencies: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name: "http_latencies",
			Help: "HTTP latencies, group by status and method",
		}, []string{"status", "method"}),
		mac: hmac.New(sha256.New, key).Sum(project),
	}
	t.registry.MustRegister(t.hits, t.latencies)
	return t
}

type PrometheusOutput struct {
	gatherers map[string]*TraefikProm
	config    *PrometheusConfig
	server    *http.Server
	listener  net.Listener
}

func PrometheusFactory(rawCfg map[string]interface{}) (api.Output, error) {
	var cfg PrometheusConfig
	err := mapstructure.Decode(rawCfg, &cfg)

	if err != nil {
		return nil, err
	}

	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	p := &PrometheusOutput{
		gatherers: make(map[string]*TraefikProm),
		config:    &cfg,
		listener:  l,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
	}
	mux.Handle("/metrics/", p)
	go p.server.Serve(l)
	return p, nil
}

func (p *PrometheusOutput) Close() error {
	return p.server.Close()
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
		p.gatherers[key] = NewTraefikProm([]byte(key), []byte(p.config.Salt))
	}
	status := fastjson.GetInt([]byte(line), "OriginStatus")
	statusXXX := fmt.Sprintf("%dxx", status/100)
	p.gatherers[key].hits.WithLabelValues(fmt.Sprintf("%d", status)).Inc()
	method := fastjson.GetString([]byte(line), "RequestMethod")
	duration := fastjson.GetFloat64([]byte(line), "Duration")
	p.gatherers[key].latencies.WithLabelValues(statusXXX, method).Observe(duration)
	return nil
}

func (p *PrometheusOutput) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var tokRaw string
	if p.config.Salt != "" {
		tokRaw = r.URL.Query().Get("t")
		if len(tokRaw) == 0 {
			w.WriteHeader(401)
			return
		}
	}

	slugs := strings.Split(r.URL.Path, "/")
	if len(slugs) < 3 {
		w.WriteHeader(404)
		return
	}
	t, ok := p.gatherers[slugs[2]]
	if !ok {
		w.WriteHeader(404)
		return
	}
	if p.config.Salt != "" {
		tok, err := base64.StdEncoding.DecodeString(tokRaw)
		if err != nil {
			log.WithField("token", tokRaw).WithError(err).Info("Can't decode base64 token")
			w.WriteHeader(400)
			return
		}
		if !hmac.Equal(t.mac, tok) {
			w.WriteHeader(401)
			return
		}
	}
	promhttp.HandlerFor(t.registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}
