package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type OutputMockup struct {
	lastTs time.Time
}

func (o *OutputMockup) Write(ctx context.Context, ts time.Time, line string, meta map[string]interface{}) error {
	o.lastTs = ts
	return nil
}

func (o *OutputMockup) Close() error {
	return nil
}

func TestJsonEngine(t *testing.T) {
	ctx := context.TODO()
	o := &OutputMockup{}
	e := NewJsonEngine(o)
	now := time.Now()
	err := e.Write(ctx, now, "{}")
	assert.NoError(t, err)
	assert.Equal(t, now, o.lastTs)
}
