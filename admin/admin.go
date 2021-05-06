package admin

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Admin struct {
	server *http.Server
}

func New(listen string, prometheus bool) *Admin {
	a := &Admin{
		server: &http.Server{
			Addr: listen,
		},
	}
	mux := http.NewServeMux()
	if prometheus {
		mux.Handle("/metrics", promhttp.Handler())
	}
	a.server.Handler = mux
	return a
}

func (a *Admin) Start(ctx context.Context) {
	go a.server.ListenAndServe()

	go func() {
		<-ctx.Done()
		a.server.Shutdown(context.TODO())
	}()

}
