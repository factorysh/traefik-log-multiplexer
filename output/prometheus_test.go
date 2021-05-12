package output

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrometheus(t *testing.T) {
	p, err := PrometheusFactory(map[string]interface{}{
		"label":  "com.docker.compose.project",
		"salt":   "queide5Weisoof7booWae7UeM2uataip",
		"listen": "127.0.0.1:0",
	})
	assert.NoError(t, err)
	err = p.Write(time.Now(), `{"OriginStatus": 200}`, map[string]interface{}{
		"com.docker.compose.project": "demo",
	})
	assert.NoError(t, err)
	pp, ok := p.(*PrometheusOutput)
	assert.True(t, ok)
	fmt.Println(pp.gatherers["demo"].hits.MetricVec)
	//assert.True(t, false)
}
