package route

import (
	"context"
	"sync"
	"testing"

	"github.com/influxdata/tail"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx, cancel := context.WithCancel(context.Background())
	r := New(ctx)
	r.SetProjectBackend("demo", "web")
	r.SetProjectBackend("demo", "chat")
	assert.Len(t, r.backends, 2)

	lines := make(chan *tail.Line)
	go r.readLines(lines)
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		for range r.projects["demo"] {
			wait.Done()
		}
	}()
	lines <- &tail.Line{
		Text: `{"Waza": "aussi"}`,
	}
	lines <- &tail.Line{
		Text: `{"BackendName": "web"}`,
	}
	wait.Wait()
	cancel()
}
