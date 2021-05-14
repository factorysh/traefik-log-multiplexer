package output

import (
	"fmt"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/imdario/mergo"
	"github.com/valyala/fastjson"

	"github.com/fluent/fluent-logger-golang/fluent"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
	Outputs["fluent"] = FluentOutputFactory
}

type FluentOutputConfig struct {
	Labels   []string
	Tag      string
	Timezone string
	Config   *fluent.Config
}

type FluentOutput struct {
	fluent   *fluent.Fluent
	config   *FluentOutputConfig
	parser   *fastjson.Parser
	location *time.Location
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
	if cfg.Timezone == "" {
		cfg.Timezone = "UTC"
	}
	l, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		return nil, err
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
		fluent:   f,
		config:   cfg,
		parser:   &fastjson.Parser{},
		location: l,
	}, nil
}

func (f *FluentOutput) Close() error {
	return f.fluent.Close()
}

func (f *FluentOutput) Write(ts time.Time, line string, meta map[string]interface{}) error {
	v, err := f.parser.Parse(line)
	if err != nil {
		return err
	}
	o, err := v.Object()
	if err != nil {
		return err
	}
	t := o.Get("time")
	if t == nil {
		return fmt.Errorf("no time in %s", line)
	}

	st, err := t.StringBytes()
	if err != nil {
		return err
	}
	tt, err := time.ParseInLocation("2006-01-02T15:04:05Z", string(st), f.location)
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	for _, label := range f.config.Labels {
		data[label] = meta[label]
	}

	return f.fluent.PostWithTime(f.config.Tag, tt, &LogMarshaler{
		line: o,
		meta: data,
	})
}
