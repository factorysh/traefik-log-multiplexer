package demultiplexer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/valyala/fastjson"

	"github.com/factorysh/traefik-log-multiplexer/admin"
	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/factorysh/traefik-log-multiplexer/conf"
	"github.com/factorysh/traefik-log-multiplexer/filter"
	"github.com/factorysh/traefik-log-multiplexer/input"
	"github.com/factorysh/traefik-log-multiplexer/output"
)

type Demultiplexer struct {
	parser  *fastjson.Parser
	input   api.Input
	filters []api.Filter
	outputs []api.Output
	closing chan error
	admin   *admin.Admin
}

func (d *Demultiplexer) Write(ts time.Time, line string) error {
	meta := make(map[string]interface{})
	data, err := d.parser.Parse(line)
	if err != nil {
		return err
	}
	for _, f := range d.filters {
		err = f.Filter(ts, data, meta)
		if err != nil {
			return err
		}
	}
	for _, o := range d.outputs {
		err = o.Write(ts, line, meta)
		if err != nil {
			return err
		}
	}
	return nil
}

func New(cfg *conf.Config) (*Demultiplexer, error) {
	d := &Demultiplexer{
		parser:  &fastjson.Parser{},
		filters: make([]api.Filter, 0),
		outputs: make([]api.Output, 0),
		closing: make(chan error),
		admin:   admin.New(cfg.Admin.Listen, cfg.Admin.Prometheus),
	}
	if len(cfg.Input) != 1 {
		return nil, fmt.Errorf("one input, not %d", len(cfg.Input))
	}
	var err error
	for name, args := range cfg.Input {
		factory, ok := input.Inputs[name]
		if !ok {
			return nil, fmt.Errorf("unknown input : %s", name)
		}
		d.input, err = factory(args)
		if err != nil {
			return nil, err
		}
		d.input.To(d)
	}
	for _, f := range cfg.Filters {
		if len(f) != 1 {
			return nil, fmt.Errorf("use only a map with 1 key, not %x", len(f))
		}
		for fName, fCfg := range f {
			factory, ok := filter.Filters[fName]
			if !ok {
				return nil, fmt.Errorf("unknown filter : %s", fName)
			}
			ff, err := factory(fCfg)
			if err != nil {
				return nil, err
			}
			d.filters = append(d.filters, ff)
		}
	}
	for name, args := range cfg.Output {
		factory, ok := output.Outputs[name]
		if !ok {
			return nil, fmt.Errorf("unknown output : %s", name)
		}
		o, err := factory(args)
		if err != nil {
			return nil, err
		}
		d.outputs = append(d.outputs, o)
	}

	return d, nil
}

func (d *Demultiplexer) Start(ctx context.Context) error {
	if os.Getenv("SENTRY_DNS") != "" {
		defer sentry.Recover()
	}
	d.admin.Start(ctx)
	go func(i api.Input) {
		err := i.Start(ctx)
		if err != nil {
			d.closing <- err
		}
	}(d.input)
	for _, f := range d.filters {
		go func(f api.Filter) {
			err := f.Start(ctx)
			if err != nil {
				d.closing <- err
			}
		}(f)
	}
	return <-d.closing
}
