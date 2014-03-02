package alcatraz

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/command_runner"
	"github.com/pivotal-cf-experimental/garden/linux_backend/port_pool"
)

type DockerBackend struct {
	skeletonPath string
	depotPath    string

	runner command_runner.CommandRunner

	portPool *port_pool.PortPool

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

		portPool: port_pool.New(51000, 10000),

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

	containerPort, err := backend.portPool.Acquire()
	if err != nil {
		return nil, err
	}

	containerPath := path.Join(backend.depotPath, id)

	cp := &exec.Cmd{
		Path: "cp",
		Args: []string{"-r", backend.skeletonPath, containerPath},
	}

	if err := backend.runner.Run(cp); err != nil {
		return nil, err
	}

	keygen := &exec.Cmd{
		Path: "ssh-keygen",
		Args: []string{
			"-f", path.Join(containerPath, "ssh", "id_rsa"),
			"-t", "rsa", "-N", "",
		},
	}

	if err := backend.runner.Run(keygen); err != nil {
		return nil, err
	}

	authorize := &exec.Cmd{
		Path: "cp",
		Args: []string{
			path.Join(containerPath, "ssh", "id_rsa.pub"),
			path.Join(containerPath, "ssh_auth", "authorized_keys"),
		},
	}

	if err := backend.runner.Run(authorize); err != nil {
		return nil, err
	}

	cmd := &exec.Cmd{
		Path: "docker",

		Args: []string{
			"run", "-d",
			"-v", path.Join(containerPath, "ssh_auth") + ":/root/.ssh:ro",
			"-v", path.Join(containerPath, "ssh_auth") + ":/home/vcap/.ssh:ro",
			"-p", fmt.Sprintf("%d:22", containerPort),
			"--name", id,
			"alcatraz",
		},
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
		containerPort,
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

	log.Println("destroyed:", container.ID())

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
