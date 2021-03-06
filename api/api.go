package api

import (
	"context"
	"time"

	"github.com/valyala/fastjson"
)

type Input interface {
	To(Writer)
	Start(context.Context) error
}

type Writer interface {
	Write(ctx context.Context, ts time.Time, line string) error
}
type Filter interface {
	Filter(ctx context.Context, ts time.Time, data *fastjson.Value, meta map[string]interface{}) error
	Start(context.Context) error
}

type Output interface {
	Write(ctx context.Context, ts time.Time, line string, meta map[string]interface{}) error
	Close() error
}

type JsonEngine struct {
	parser  fastjson.Parser
	filters []Filter
	output  Output
}

func NewJsonEngine(output Output, filters ...Filter) *JsonEngine {
	return &JsonEngine{
		filters: filters,
		output:  output,
	}
}

func (j *JsonEngine) Write(ctx context.Context, ts time.Time, line string) error {
	p, err := j.parser.Parse(line)
	if err != nil {
		return err
	}
	meta := make(map[string]interface{})
	for _, filter := range j.filters {
		err = filter.Filter(ctx, ts, p, meta)
		if err != nil {
			return err
		}
	}
	return j.output.Write(ctx, ts, line, meta)
}
