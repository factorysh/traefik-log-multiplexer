package main

import (
	"context"
	"fmt"
	"os"

	"github.com/factorysh/traefik-log-multiplexer/conf"
	"github.com/factorysh/traefik-log-multiplexer/route"
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

	cfg, err := conf.New()
	if err != nil {
		log.WithError(err).Error()
		return
	}
	ctx := context.Background()
	r := route.New(ctx)
	err = r.Read(cfg.LogPath)
	l := log.WithField("logPath", cfg.LogPath)
	if err != nil {
		l.WithError(err).Error()
		return
	}
}
