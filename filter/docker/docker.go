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
	"github.com/valyala/fastjson"
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
	return nil
}

func (dp *DockerProvider) visitor(action string, container *types.ContainerJSON) {
	switch action {
	case visitor.START:
		dp.containers.Add(container)
	case visitor.STOP:
		dp.containers.Remove(container.ID)
	case visitor.DIE:
		dp.containers.Remove(container.ID)
	}
}

func New(client *client.Client) *DockerProvider {
	dp := &DockerProvider{
		containers: NewContainers(),
		watcher:    visitor.New(client),
	}
	dp.watcher.VisitCurrentCointainer(dp.create)
	dp.watcher.WatchFor(dp.visitor)

	return dp
}

func (dp *DockerProvider) Start(ctx context.Context) error {
	return dp.watcher.Start(ctx)
}

func (dp *DockerProvider) Filter(ts time.Time, data *fastjson.Value, meta map[string]interface{}) error {
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
	return nil
}
