package output

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/influxdata/tail"
	log "github.com/sirupsen/logrus"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]NewOutput)
	}
	Outputs["file"] = func(config map[string]interface{}) (Output, error) {
		pathRaw, ok := config["path_pattern"]
		if !ok {
			return nil, errors.New("path_pattern key is mandatory")
		}
		path, ok := pathRaw.(string)
		if !ok {
			return nil, errors.New("path_pattern must be a string")
		}
		return NewFile(path), nil
	}
}

// File writes logs to file
type File struct {
	pathPattern string
	files       map[string]*os.File
	lock        sync.RWMutex
}

// NewFile returns new File, with a path pattern
func NewFile(pathPattern string) *File {
	return &File{
		pathPattern: pathPattern,
		files:       make(map[string]*os.File),
	}
}

func (f *File) Read(project string, line *tail.Line) {
	f.lock.RLock()
	file, ok := f.files[project]
	f.lock.RUnlock()
	if !ok {
		var err error
		file, err = os.Open(fmt.Sprintf(f.pathPattern, project))
		if err != nil {
			log.WithField("project", project).WithError(err).Error()
			return
		}
		f.lock.Lock()
		f.files[project] = file
		f.lock.Unlock()
	}
	file.Write([]byte(line.Text))
	// TODO flush anytime ?
}

// RemoveProject removes a project
func (f *File) RemoveProject(project string) {
	l := log.WithField("project", project)
	f.lock.Lock()
	defer f.lock.Unlock()
	file, ok := f.files[project]
	if ok {
		err := file.Close()
		if err != nil {
			l.WithError(err).Error()
		}
		delete(f.files, project)
		l.Info("Removed")
	} else {
		l.Warning("Try to remove unknown project")
	}
}
