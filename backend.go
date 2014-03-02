package alcatraz

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/command_runner"
)

type DockerBackend struct {
	skeletonPath string
	depotPath    string

	runner command_runner.CommandRunner

	containerIDs <-chan string

	containers      map[string]backend.Container
	containersMutex *sync.RWMutex
}

type UnknownHandleError struct {
	Handle string
}

func (e UnknownHandleError) Error() string {
	return "unknown handle: " + e.Handle
}

func New(skeletonPath string, depotPath string, runner command_runner.CommandRunner) *DockerBackend {
	containerIDs := make(chan string)

	go generateContainerIDs(containerIDs)

	return &DockerBackend{
		skeletonPath: skeletonPath,
		depotPath:    depotPath,

		runner: runner,

		containerIDs: containerIDs,

		containers:      make(map[string]backend.Container),
		containersMutex: new(sync.RWMutex),
	}
}

func (backend *DockerBackend) Setup() error {
	log.Println("Setup")
	return os.MkdirAll(backend.depotPath, 0755)
}

func (backend *DockerBackend) Start() error {
	log.Println("Start")
	return nil
}

func (backend *DockerBackend) Stop() {
	log.Println("Stop")
}

func (backend *DockerBackend) Create(spec backend.ContainerSpec) (backend.Container, error) {
	id := <-backend.containerIDs

	containerPath := path.Join(backend.depotPath, id)

	cp := &exec.Cmd{
		Path: "cp",
		Args: []string{"-a", backend.skeletonPath, containerPath},
	}

	if err := backend.runner.Run(cp); err != nil {
		return nil, err
	}

	out := new(bytes.Buffer)

	cmd := &exec.Cmd{
		Path: "docker",

		Args: []string{
			"run", "-d",

			"--name", id,

			// TODO: this should only be writable by privileged root
			"-v", containerPath + ":/share",

			"alcatraz",
		},

		Stdout: out,
	}

	if err := backend.runner.Run(cmd); err != nil {
		return nil, err
	}

	handle := id
	if spec.Handle != "" {
		handle = spec.Handle
	}

	log.Println("created:", handle)

	container := NewDockerContainer(
		id,
		handle,
		containerPath,
		backend.runner,
	)

	backend.containersMutex.Lock()
	backend.containers[handle] = container
	backend.containersMutex.Unlock()

	return container, nil
}

func (backend *DockerBackend) Destroy(handle string) error {
	container, err := backend.Lookup(handle)
	if err != nil {
		return err
	}

	cmd := &exec.Cmd{
		Path: "docker",
		Args: []string{"kill", container.ID()},
	}

	err = backend.runner.Run(cmd)
	if err != nil {
		return err
	}

	log.Println("Destroyed:", container.ID())

	return nil
}

func (backend *DockerBackend) Containers() (containers []backend.Container, err error) {
	backend.containersMutex.RLock()
	defer backend.containersMutex.RUnlock()

	for _, container := range backend.containers {
		containers = append(containers, container)
	}

	return containers, nil
}

func (backend *DockerBackend) Lookup(handle string) (backend.Container, error) {
	backend.containersMutex.RLock()
	defer backend.containersMutex.RUnlock()

	container, found := backend.containers[handle]
	if !found {
		return nil, UnknownHandleError{handle}
	}

	return container, nil
}

func generateContainerIDs(ids chan<- string) string {
	for containerNum := time.Now().UnixNano(); ; containerNum++ {
		containerID := []byte{}

		var i uint
		for i = 0; i < 11; i++ {
			containerID = strconv.AppendInt(
				containerID,
				(containerNum>>(55-(i+1)*5))&31,
				32,
			)
		}

		ids <- string(containerID)
	}
}
