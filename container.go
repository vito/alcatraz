package alcatraz

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path"
	"strconv"
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
	"github.com/pivotal-cf-experimental/garden/command_runner"
	"github.com/pivotal-cf-experimental/garden/linux_backend/process_tracker"

	"github.com/vito/alcatraz/lxc_cgroups_manager"
)

type DockerContainer struct {
	id     string
	handle string
	path   string
	port   uint32

	runner command_runner.CommandRunner

	cgroupsManager *lxc_cgroups_manager.LXCCgroupsManager

	processTracker *process_tracker.ProcessTracker
}

func NewDockerContainer(
	id, handle, path string,
	port uint32,
	runner command_runner.CommandRunner,
	cgroupsManager *lxc_cgroups_manager.LXCCgroupsManager,
) *DockerContainer {
	return &DockerContainer{
		id:     id,
		handle: handle,
		path:   path,
		port:   port,

		runner: runner,

		cgroupsManager: cgroupsManager,

		processTracker: process_tracker.New(path, runner),
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
	return container.rsync(src, "vcap@127.0.0.1:"+dst)
}

func (container *DockerContainer) CopyOut(src, dst, owner string) error {
	log.Println(container.id, "copying out from", src, "to", dst)

	err := container.rsync("vcap@127.0.0.1:"+src, dst)
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
	limit := fmt.Sprintf("%d", limits.LimitInBytes)

	// memory.memsw.limit_in_bytes must be >= memory.limit_in_bytes
	//
	// however, it must be set after memory.limit_in_bytes, and if we're
	// increasing the limit, writing memory.limit_in_bytes first will fail.
	//
	// so, write memory.limit_in_bytes before and after
	container.cgroupsManager.Set("memory", "memory.limit_in_bytes", limit)
	container.cgroupsManager.Set("memory", "memory.memsw.limit_in_bytes", limit)

	err := container.cgroupsManager.Set("memory", "memory.limit_in_bytes", limit)
	if err != nil {
		return err
	}

	return nil
}

func (container *DockerContainer) CurrentMemoryLimits() (backend.MemoryLimits, error) {
	limitInBytes, err := container.cgroupsManager.Get("memory", "memory.limit_in_bytes")
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	numericLimit, err := strconv.ParseUint(limitInBytes, 10, 0)
	if err != nil {
		return backend.MemoryLimits{}, err
	}

	return backend.MemoryLimits{uint64(numericLimit)}, nil
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
	log.Println(container.id, "running process:", spec.Script)

	user := "vcap"
	if spec.Privileged {
		user = "root"
	}

	wsh := &exec.Cmd{
		Path: "ssh",
		Args: []string{
			"-q",
			"-i", path.Join(container.path, "ssh", "id_rsa"),
			"-p", fmt.Sprintf("%d", container.port),
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			user + "@127.0.0.1",
			"/bin/bash",
		},
		Stdin: bytes.NewBufferString(spec.Script),
	}

	return container.processTracker.Run(wsh)
}

func (container *DockerContainer) Attach(processID uint32) (<-chan backend.ProcessStream, error) {
	log.Println(container.id, "attaching to process", processID)
	return container.processTracker.Attach(processID)
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
	rsync := &exec.Cmd{
		Path: "rsync",
		Args: []string{
			"-e",
			fmt.Sprintf(
				"ssh -i %s -p %d -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null",
				path.Join(container.path, "ssh", "id_rsa"),
				container.port,
			),
			"-r",
			"-p",
			"--links",
			src,
			dst,
		},
	}

	return container.runner.Run(rsync)
}
