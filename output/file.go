package output

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/factorysh/traefik-log-multiplexer/tpl"
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
	PathPattern string `mapstructure:"path_pattern"`
}

func FileOutputFactory(rawCfg map[string]interface{}) (api.Output, error) {
	var cfg FileOutputConfig
	err := mapstructure.Decode(rawCfg, &cfg)
	if err != nil {
		return nil, err
	}
	if cfg.PathPattern == "" {
		return nil, fmt.Errorf("path_pattern is mandatory")
	}
	return NewFile(cfg.PathPattern)

}

// File writes logs to file
type File struct {
	pathPattern string
	template    *tpl.Template
	files       map[string]*os.File
	lock        sync.RWMutex
}

// NewFile returns new File, with a path pattern
func NewFile(pathPattern string) (*File, error) {
	t, err := tpl.Parse([]byte(pathPattern))
	if err != nil {
		return nil, err
	}
	return &File{
		pathPattern: pathPattern,
		template:    t,
		files:       make(map[string]*os.File),
	}, nil
}

func (f *File) Write(ctx context.Context, ts time.Time, line string, meta map[string]interface{}) error {
	path, err := f.template.Execute(meta)
	if err != nil {
		return err
	}
	spath := string(path)
	f.lock.RLock()
	file, ok := f.files[spath]
	f.lock.RUnlock()

	if !ok {
		file, err = os.OpenFile(spath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			log.WithField("path", spath).WithError(err).Error()
			return err
		}
		f.lock.Lock()
		f.files[spath] = file
		f.lock.Unlock()
	}
	file.Write([]byte(line))
	// TODO flush anytime ?
	return nil
}

func (f *File) Close() error {
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, ff := range f.files {
		err := ff.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
