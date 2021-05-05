package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type OutputMockup struct {
	lastTs time.Time
}

func (o *OutputMockup) Write(ts time.Time, line string, meta map[string]interface{}) error {
	o.lastTs = ts
	return nil
}

func TestJsonEngine(t *testing.T) {
	o := &OutputMockup{}
	e := NewJsonEngine(o)
	now := time.Now()
	err := e.Write(now, "{}")
	assert.NoError(t, err)
	assert.Equal(t, now, o.lastTs)
}
