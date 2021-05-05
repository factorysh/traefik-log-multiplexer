package main

import (
	"context"
	"fmt"
	"os"

	"github.com/factorysh/traefik-log-multiplexer/conf"
	"github.com/factorysh/traefik-log-multiplexer/demultiplexer"
	_ "github.com/factorysh/traefik-log-multiplexer/filter/docker"
	_ "github.com/factorysh/traefik-log-multiplexer/input"
	_ "github.com/factorysh/traefik-log-multiplexer/output"
	"github.com/factorysh/traefik-log-multiplexer/version"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Logrus hook for adding file name and line to logs
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version.Version())
		return
	}

	cfg, err := conf.Read()
	if err != nil {
		log.WithError(err).Error()
		return
	}

	l := log.WithField("path", cfg.Path)

	d, err := demultiplexer.New(cfg)
	if err != nil {
		l.WithError(err).Error()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = d.Start(ctx)
	if err != nil {
		l.WithError(err).Error()
		return
	}
}
