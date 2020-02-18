package fusebox

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/influxdata/tail"
	"github.com/stretchr/testify/assert"
)

func TestFusebox(t *testing.T) {
	f := New(3)
	assert.Equal(t, 0, f.start)
	assert.Equal(t, 0, f.size)
	ok := f.Push(&tail.Line{
		Text: "pim",
	})
	assert.True(t, ok)
	ok = f.Push(&tail.Line{
		Text: "pam",
	})
	assert.True(t, ok)
	ok = f.Push(&tail.Line{
		Text: "poum",
	})
	assert.True(t, ok)
	ok = f.Push(&tail.Line{
		Text: "The Captain",
	})
	assert.False(t, ok)
	debug(f.queue)
	l := f.Pop()
	assert.NotNil(t, l)
	assert.Equal(t, "pim", l.Text)
	assert.Equal(t, 2, f.size)
	ok = f.Push(&tail.Line{
		Text: "The Astronom",
	})
	assert.True(t, ok)
	debug(f.queue)
}

func TestBlocking(t *testing.T) {
	f := New(3)
	l := sync.WaitGroup{}
	l.Add(1)
	go func() {
		line := f.BlockingPop(context.Background())
		assert.NotNil(t, line)
		assert.Equal(t, "Fritz", line.Text)
		l.Done()

	}()
	time.Sleep(10 * time.Millisecond)
	f.Push(&tail.Line{
		Text: "Fritz",
	})
	l.Wait()
}

func debug(lines []*tail.Line) {
	for _, l := range lines {
		if l != nil {
			fmt.Print(l.Text, " ")
		}
	}
	fmt.Println()

}
