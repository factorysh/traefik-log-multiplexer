package output

import (
	"context"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/grafana/loki/clients/pkg/promtail"
)

func init() {
	Outputs["loki"] = LokiOutputFactory
}

type LokiOutput struct {
	client promtail.Promtail
}

func (l *LokiOutput) Write(ctx context.Context, ts time.Time, line string, meta map[string]interface{}) error {
	return nil
}

func (l *LokiOutput) Close() error {
	return nil
}

func LokiOutputFactory(cfg map[string]interface{}) (api.Output, error) {
	return &LokiOutput{}, nil
}
