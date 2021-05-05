package output

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/influxdata/tail"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func init() {
	if Outputs == nil {
		Outputs = make(map[string]OutputFactory)
	}
	Outputs["file"] = FileOutputFactory
}

type FileOutputConfig struct {
	PathPattern string `yaml:"path_pattern"`
}

func FileOutputFactory(rawCfg map[string]interface{}) (api.Output, error) {
	var cfg FileOutputConfig
	err := mapstructure.Decode(rawCfg, &cfg)
	if err != nil {
		return nil, err
	}
	return NewFile(cfg.PathPattern), nil

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

func (f *File) Write(ts time.Time, line string, meta map[string]interface{}) error {
	return nil
}

// FIXME
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
