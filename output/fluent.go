package output

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/imdario/mergo"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
	Outputs["fluent"] = FluentOutputFactory
}

type FluentOutputConfig struct {
	Labels []string
	Tag    string
	Config *fluent.Config
}

type FluentOutput struct {
	fluent *fluent.Fluent
	config *FluentOutputConfig
}

func FluentOutputFactory(rawCfg map[string]interface{}) (api.Output, error) {
	cfg := &FluentOutputConfig{
		Labels: make([]string, 0),
		Config: &fluent.Config{},
	}
	labels, ok := rawCfg["labels"]
	if ok {
		labelsRaw, ok := labels.([]interface{})
		if !ok {
			cfg.Labels, ok = labels.([]string)
			if !ok {
				return nil, fmt.Errorf("bad labels format : %v %T", labels, labels)
			}
		} else {
			cfg.Labels = make([]string, len(labelsRaw))
			for i, blob := range labelsRaw {
				cfg.Labels[i], ok = blob.(string)
				if !ok {
					return nil, fmt.Errorf("bad label format : %v %T", blob, blob)
				}
			}
		}
	}
	fluentCfg, ok := rawCfg["config"]
	if ok {
		err := mergo.Map(cfg.Config, fluentCfg)
		if err != nil {
			return nil, err
		}
		mergo.Merge(cfg.Config, &fluent.Config{
			FluentHost:         "127.0.0.1",
			SubSecondPrecision: true,
		})
	}
	f, err := fluent.New(*cfg.Config)
	if err != nil {
		return nil, err
	}
	return &FluentOutput{
		fluent: f,
		config: cfg,
	}, nil
}

func (f *FluentOutput) Close() error {
	return f.fluent.Close()
}

func (f *FluentOutput) Write(ts time.Time, line string, meta map[string]interface{}) error {
	var message map[string]interface{}
	err := json.Unmarshal([]byte(line), &message)
	if err != nil {
		return err
	}
	timeRaw, ok := message["time"]
	if !ok {
		return fmt.Errorf("where is my time")
	}
	t, ok := timeRaw.(string)
	if !ok {
		return fmt.Errorf("bad time format : %v", timeRaw)
	}
	tt, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	for _, label := range f.config.Labels {
		data[label] = meta[label]
	}
	message["docker_provider"] = data
	return f.fluent.PostWithTime(f.config.Tag, tt, message)
}
