package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/factorysh/docker-visitor/visitor"
)

type DockerProvider struct {
	containers *Containers
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

func New(ctx context.Context, client *client.Client) *DockerProvider {
	dp := &DockerProvider{
		containers: NewContainers(),
	}
	w := visitor.New(client)
	w.VisitCurrentCointainer(dp.create)
	w.WatchFor(dp.visitor)

	go w.Start(ctx)
	return dp
}
