package docker

type Service struct {
	Name       string `json:"name"`
	Containers string `json:"container"`
}
