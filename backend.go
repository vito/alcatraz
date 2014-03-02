package alcatraz

import (
	"github.com/pivotal-cf-experimental/garden/backend"
)

type DockerBackend struct {
}

func New() *DockerBackend {
	return &DockerBackend{}
}

func (backend *DockerBackend) Setup() error {
	panic("TODO Setup")
	return nil
}

func (backend *DockerBackend) Start() error {
	panic("TODO Start")
	return nil
}

func (backend *DockerBackend) Stop() {
	panic("TODO Stop")
}

func (backend *DockerBackend) Create(spec backend.ContainerSpec) (backend.Container, error) {
	panic("TODO Create")
	return nil, nil
}

func (backend *DockerBackend) Destroy(handle string) error {
	panic("TODO Destroy")
	return nil
}

func (backend *DockerBackend) Containers() ([]backend.Container, error) {
	panic("TODO Containers")
	return nil, nil
}

func (backend *DockerBackend) Lookup(handle string) (backend.Container, error) {
	panic("TODO Lookup")
	return nil, nil
}
