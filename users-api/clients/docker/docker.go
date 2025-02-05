package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	Client *client.Client
}

func NewDockerClient() DockerClient {
	client, err := client.NewClientWithOpts(client.WithVersion("1.45"))
	if err != nil {
		log.Fatalf("error initializing docker client: %s", err.Error())
	}
	return DockerClient{
		Client: client,
	}
}

// ListarContenedoresActivos obtiene los contenedores en ejecución.
func (client DockerClient) GetContainers(ctx context.Context) ([]types.Container, error) {
	containers, err := client.Client.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("error getting container list: %w", err)
	}

	return containers, err
}
