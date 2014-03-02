package alcatraz

import (
	"log"
	"os/exec"
	"path"
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/command_runner"
)

type DockerContainer struct {
	id     string
	handle string
	path   string

	runner command_runner.CommandRunner
}

func NewDockerContainer(
	id, handle, path string,

	runner command_runner.CommandRunner,
) *DockerContainer {
	return &DockerContainer{
		id:     id,
		handle: handle,
		path:   path,

		runner: runner,
	}
}

func (container *DockerContainer) ID() string {
	return container.id
}

func (container *DockerContainer) Handle() string {
	return container.handle
}

func (container *DockerContainer) GraceTime() time.Duration {
	log.Println("TODO GraceTime")
	return 0
}

func (container *DockerContainer) Stop(kill bool) error {
	log.Println("TODO Stop")
	return nil
}

func (container *DockerContainer) Info() (backend.ContainerInfo, error) {
	log.Println("TODO Info")
	return backend.ContainerInfo{}, nil
}

func (container *DockerContainer) CopyIn(src, dst string) error {
	log.Println(container.id, "copying in from", src, "to", dst)
	return container.rsync(src, "vcap@container:"+dst)
}

func (container *DockerContainer) CopyOut(src, dst, owner string) error {
	log.Println(container.id, "copying out from", src, "to", dst)

	err := container.rsync("vcap@container:"+src, dst)
	if err != nil {
		return err
	}

	if owner != "" {
		chown := &exec.Cmd{
			Path: "chown",
			Args: []string{"-R", owner, dst},
		}

		err := container.runner.Run(chown)
		if err != nil {
			return err
		}
	}

	return nil
}

func (container *DockerContainer) LimitBandwidth(limits backend.BandwidthLimits) error {
	log.Println("TODO LimitBandwidth")
	return nil
}

func (container *DockerContainer) CurrentBandwidthLimits() (backend.BandwidthLimits, error) {
	log.Println("TODO CurrentBandwidthLimits")
	return backend.BandwidthLimits{}, nil
}

func (container *DockerContainer) LimitDisk(limits backend.DiskLimits) error {
	log.Println("TODO LimitDisk")
	return nil
}

func (container *DockerContainer) CurrentDiskLimits() (backend.DiskLimits, error) {
	log.Println("TODO CurrentDiskLimits")
	return backend.DiskLimits{}, nil
}

func (container *DockerContainer) LimitMemory(limits backend.MemoryLimits) error {
	log.Println("TODO LimitMemory")
	return nil
}

func (container *DockerContainer) CurrentMemoryLimits() (backend.MemoryLimits, error) {
	log.Println("TODO CurrentMemoryLimits")
	return backend.MemoryLimits{}, nil
}

func (container *DockerContainer) LimitCPU(limits backend.CPULimits) error {
	log.Println("TODO LimitCPU")
	return nil
}

func (container *DockerContainer) CurrentCPULimits() (backend.CPULimits, error) {
	log.Println("TODO CurrentCPULimits")
	return backend.CPULimits{}, nil
}

func (container *DockerContainer) Run(spec backend.ProcessSpec) (uint32, <-chan backend.ProcessStream, error) {
	log.Println("TODO Run")
	return 0, nil, nil
}

func (container *DockerContainer) Attach(processID uint32) (<-chan backend.ProcessStream, error) {
	log.Println("TODO Attach")
	return nil, nil
}

func (container *DockerContainer) NetIn(hostPort uint32, containerPort uint32) (uint32, uint32, error) {
	log.Println("TODO NetIn")
	return hostPort, containerPort, nil
}

func (container *DockerContainer) NetOut(network string, port uint32) error {
	log.Println("TODO NetOut")
	return nil
}

func (container *DockerContainer) rsync(src, dst string) error {
	wshPath := path.Join(container.path, "bin", "wsh")
	sockPath := path.Join(container.path, "run", "wshd.sock")

	rsync := &exec.Cmd{
		Path: "rsync",
		Args: []string{
			"-e", wshPath + " --socket " + sockPath + " --rsh",
			"-r",
			"-p",
			"--links",
			src,
			dst,
		},
	}

	return container.runner.Run(rsync)
}
