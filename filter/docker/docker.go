package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
	"github.com/factorysh/traefik-log-multiplexer/api"
	"github.com/factorysh/traefik-log-multiplexer/filter"
	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/valyala/fastjson"
)

var (
	dockerContainerEvents = promauto.NewCounter(prometheus.CounterOpts{
		Name: "demultiplexer_docker_events",
		Help: "The total number of docker events",
	})
)

func init() {
	if filter.Filters == nil {
		filter.Filters = map[string]filter.FilterFactory{}
	}
	filter.Filters["docker"] = DockerFactory

}

type DockerProvider struct {
	containers *Containers
	watcher    *visitor.Watcher
	client     *client.Client
}

func DockerFactory(cfg map[string]interface{}) (api.Filter, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return New(docker), nil
}

func (dp *DockerProvider) create(container *types.ContainerJSON) error {
	dp.containers.Add(container)
	dockerContainerEvents.Inc()
	return nil
}

func (dp *DockerProvider) visitor(action string, container *types.ContainerJSON) {
	switch action {
	case visitor.START:
		dp.containers.Add(container)
		dockerContainerEvents.Inc()
	case visitor.STOP:
		dp.containers.Remove(container.ID)
		dockerContainerEvents.Inc()
	case visitor.DIE:
		dp.containers.Remove(container.ID)
		dockerContainerEvents.Inc()
	}
}

func New(client *client.Client) *DockerProvider {
	dp := &DockerProvider{
		containers: NewContainers(),
		watcher:    visitor.New(client),
		client:     client,
	}
	dp.watcher.VisitCurrentCointainer(dp.create)
	dp.watcher.WatchFor(dp.visitor)

	return dp
}

func (dp *DockerProvider) Start(ctx context.Context) error {
	if hub := sentry.CurrentHub(); hub != nil {
		info, err := dp.client.Info(context.TODO())
		if err != nil {
			return err
		}
		hub.AddBreadcrumb(&sentry.Breadcrumb{
			Category: "Filter: docker",
			Data: map[string]interface{}{
				"server_version": info.ServerVersion,
				"client_version": dp.client.ClientVersion(),
			},
		}, nil)
	}
	return dp.watcher.Start(ctx)
}

func (dp *DockerProvider) Filter(ctx context.Context, ts time.Time, data *fastjson.Value, meta map[string]interface{}) error {
	backend := data.GetStringBytes("BackendAddr")
	if len(backend) == 0 {
		return nil
	}
	container := dp.containers.ByListen(string(backend))
	if container == nil {
		return fmt.Errorf("can't find backend %s", string(backend))
	}
	labels := []string{"sh.factory.project"}
	for _, label := range labels {
		meta[label] = container.Config.Labels[label]
	}
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.Scope().SetTag("project", container.Config.Labels["sh.factory.project"])
	}
	return nil
}
