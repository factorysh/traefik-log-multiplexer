package fusebox

import (
	"context"
	"fmt"
	"sync"

	"github.com/influxdata/tail"
)

type Fusebox struct {
	queue   []*tail.Line
	start   int
	size    int
	max     int
	lock    sync.Mutex
	waitFor chan struct{}
}

// New Fusebox
func New(size int) *Fusebox {
	return &Fusebox{
		queue: make([]*tail.Line, size),
		start: 0,
		max:   size,
	}
}

func (f *Fusebox) last() int {
	return (f.start + f.size) % f.max
}

// Push a line to the stack
func (f *Fusebox) Push(line *tail.Line) bool {
	if f.size >= f.max { // the queue is full, lets drop things
		return false
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.queue[f.last()] = line
	f.size++
	if f.waitFor != nil {
		f.waitFor <- struct{}{}
		close(f.waitFor)
	}
	return true
}

// Pop a line from the stack
func (f *Fusebox) Pop() *tail.Line {
	if f.size == 0 {
		return nil
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.size--
	l := f.queue[f.start]
	f.start = (f.start + 1) % f.max
	return l
}

func (f *Fusebox) BlockingPop(ctx context.Context) *tail.Line {
	l := f.Pop()
	if l != nil {
		return l
	}
	f.lock.Lock()
	f.waitFor = make(chan struct{})
	f.lock.Unlock()
	fmt.Println("Waiting")
	select {
	case <-f.waitFor:
	//
	case <-ctx.Done():
		close(f.waitFor)
		return nil
	}
	// assert Pop is not nil
	return f.Pop()
}

func (f *Fusebox) Read(line *tail.Line) {

}
