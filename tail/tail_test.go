package tail

import (
	"context"
	"testing"

	"github.com/influxdata/tail"
)

func TestRouter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r := New(ctx)
	r.SetProjectBackend("demo", "web")
	r.SetProjectBackend("demo", "chat")

	lines := make(chan *tail.Line)
	go r.readLines(lines)
	lines <- &tail.Line{
		Text: `{"BackendName": "web"}`,
	}
	cancel()
}
