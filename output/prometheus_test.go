package output

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestPrometheus(t *testing.T) {
	salt := "queide5Weisoof7booWae7UeM2uataip"
	p, err := PrometheusFactory(map[string]interface{}{
		"label": "com.docker.compose.project",
		"salt":  salt,
		"addr":  "127.0.0.1:0",
	})
	assert.NoError(t, err)
	defer p.Close()
	err = p.Write(context.TODO(), time.Now(), `{
		"OriginStatus": 200,
		"RequestMethod":"GET",
		"Duration":42
		}`, map[string]interface{}{
		"com.docker.compose.project": "demo",
	})
	assert.NoError(t, err)
	pp, ok := p.(*PrometheusOutput)
	assert.True(t, ok)
	var m dto.Metric
	err = pp.gatherers["demo"].hits.With(prometheus.Labels{"status": "2xx"}).Write(&m)
	assert.NoError(t, err)
	v := m.Counter.GetValue()
	assert.Equal(t, float64(1), v)

	s := pp.listener.Addr().String()

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics/demo", s))
	assert.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	tok := base64.StdEncoding.EncodeToString(hmac.New(sha256.New, []byte(salt)).Sum([]byte("demo")))
	resp, err = http.Get(fmt.Sprintf("http://%s/metrics/demo?t=%s", s, tok))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, strings.HasPrefix(resp.Header.Get("content-type"), "text/plain"))
	resp, err = http.Get(fmt.Sprintf("http://%s/metrics/plop?t=%s", s, tok))
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

}
