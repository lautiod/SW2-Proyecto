package docker

type Service struct {
	Name        string `json:"name"`
	ContainerID string `json:"container_id"`
}
