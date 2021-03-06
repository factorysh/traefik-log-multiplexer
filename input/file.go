package input

import (
	"context"

	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/getsentry/sentry-go"
	"github.com/influxdata/tail"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

func init() {
	if Inputs == nil {
		Inputs = make(map[string]InputFactory)
	}
	Inputs["file"] = NewFileInput
}

type FileInput struct {
	path   string
	writer api.Writer
}

type FileInputConfig struct {
	Path string
}

func (f *FileInput) To(writer api.Writer) {
	f.writer = writer
}

func NewFileInput(rawCfg map[string]interface{}) (api.Input, error) {
	var cfg FileInputConfig
	err := mapstructure.Decode(rawCfg, &cfg)
	if err != nil {
		return nil, err
	}
	return &FileInput{
		path: cfg.Path,
	}, nil
}

func (f *FileInput) Start(ctx context.Context) error {
	l := log.WithField("path", f.path)
	if hub := sentry.CurrentHub(); hub != nil {
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "Input: file",
			Data: map[string]interface{}{
				"path": f.path,
			},
		}, &sentry.BreadcrumbHint{})
	}
	for {
		t, err := tail.TailFile(f.path, tail.Config{Follow: true})
		if err != nil {
			return err
		}
		l.Info("tail")
		again := true
		for again {
			select {
			case <-ctx.Done():
				l.Info("Stop reading")
				return nil
			case line := <-t.Lines:
				if line == nil {
					l.Warn("Oups the file has vanished")
					again = false
					break
				}
				if line.Err != nil {
					l.WithField("line", line).WithError(err).Info("Tail")
					continue
				}
				// a fresh context, just for this line
				err = f.writer.Write(context.TODO(), line.Time, line.Text)
				if err != nil {
					l.WithField("line", line).WithError(err).Info("Write")
					continue
				}
			}
		}
	}
}
