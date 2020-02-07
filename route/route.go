package route

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/influxdata/tail"
	log "github.com/sirupsen/logrus"
)

type Router struct {
	backends map[string]string // backend => project
	projects map[string]chan *tail.Line
	lock     sync.RWMutex
	context  context.Context
}

func New(ctx context.Context) *Router {
	return &Router{
		backends: make(map[string]string),
		projects: make(map[string]chan *tail.Line),
		context:  ctx,
	}
}

func (r *Router) project(backend string) (chan *tail.Line, bool) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	p, ok := r.backends[backend]
	if !ok { // backend unknown
		log.Warn("Unknown backend : ", backend)
		return nil, false
	}
	c, ok := r.projects[p]
	if ok {
		return c, true
	}
	log.Error("No project for backend ", backend)
	return nil, false
}

func (r *Router) SetProjectBackend(project, backend string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	_, ok := r.projects[project]
	if !ok {
		r.projects[project] = make(chan *tail.Line)
	}
	r.backends[backend] = project
}

func (r *Router) RemoveBackend(backend string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.backends, backend)
	// remove orphans projects
}

func (r *Router) readLines(lines chan *tail.Line) {
	log.Info("Reading")
	for {
		select {
		case <-r.context.Done():
			log.Info("Stop reading")
			return
		case line := <-lines:
			var blob map[string]interface{}
			err := json.Unmarshal([]byte(line.Text), &blob)
			if err != nil {
				log.WithError(err).Warn()
				continue
			}
			backendRaw, ok := blob["BackendName"]
			if !ok {
				err = errors.New("This log line hasn't BackendName key")
				log.WithError(err).Warn()
				continue
			}
			backend, ok := backendRaw.(string)
			if !ok {
				err = errors.New("BackendName is not a string")
				log.WithError(err).Warn()
				continue
			}
			reader, ok := r.project(backend)
			if ok {
				reader <- line
			}
		}
	}
}

func (r *Router) Read(path string) error {
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		return err
	}
	r.readLines(t.Lines)
	return nil
}
