package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/factorysh/traefik-log-multiplexer/conf"
	"github.com/factorysh/traefik-log-multiplexer/demultiplexer"
	_ "github.com/factorysh/traefik-log-multiplexer/filter/docker"
	_ "github.com/factorysh/traefik-log-multiplexer/input"
	_ "github.com/factorysh/traefik-log-multiplexer/output"
	"github.com/factorysh/traefik-log-multiplexer/version"
	"github.com/getsentry/sentry-go"
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

	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			// Either set your DSN here or set the SENTRY_DSN environment variable.
			Dsn: dsn,

			// Either set environment and release here or set the SENTRY_ENVIRONMENT
			// and SENTRY_RELEASE environment variables.
			Environment: "",
			Release:     fmt.Sprintf("traefik-log-multiplexer@%s", version.Version()),
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug:            true,
			AttachStacktrace: true,
		})
		if err != nil {
			log.WithError(err).Error()
			return
		}
		// Flush buffered events before the program terminates.
		// Set the timeout to the maximum duration the program can afford to wait.
		defer sentry.Flush(2 * time.Second)
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
