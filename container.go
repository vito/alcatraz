package alcatraz

import (
	"time"

	"github.com/pivotal-cf-experimental/garden/backend"
)

type DockerContainer struct {
}

func (container *DockerContainer) ID() string {
	return ""
}

func (container *DockerContainer) Handle() string {
	return ""
}

func (container *DockerContainer) GraceTime() time.Duration {
	panic("TODO GraceTime")
	return 0
}

func (container *DockerContainer) Start() error {
	panic("TODO Start")
	return nil
}

func (container *DockerContainer) Stop(kill bool) error {
	panic("TODO Stop")
	return nil
}

func (container *DockerContainer) Info() (backend.ContainerInfo, error) {
	panic("TODO Info")
	return backend.ContainerInfo{}, nil
}

func (container *DockerContainer) CopyIn(src, dst string) error {
	panic("TODO CopyIn")
	return nil
}

func (container *DockerContainer) CopyOut(src, dst, owner string) error {
	panic("TODO CopyOut")
	return nil
}

func (container *DockerContainer) LimitBandwidth(limits backend.BandwidthLimits) error {
	panic("TODO LimitBandwidth")
	return nil
}

func (container *DockerContainer) CurrentBandwidthLimits() (backend.BandwidthLimits, error) {
	panic("TODO CurrentBandwidthLimits")
	return backend.BandwidthLimits{}, nil
}

func (container *DockerContainer) LimitDisk(limits backend.DiskLimits) error {
	panic("TODO LimitDisk")
	return nil
}

func (container *DockerContainer) CurrentDiskLimits() (backend.DiskLimits, error) {
	panic("TODO CurrentDiskLimits")
	return backend.DiskLimits{}, nil
}

func (container *DockerContainer) LimitMemory(limits backend.MemoryLimits) error {
	panic("TODO LimitMemory")
	return nil
}

func (container *DockerContainer) CurrentMemoryLimits() (backend.MemoryLimits, error) {
	panic("TODO CurrentMemoryLimits")
	return backend.MemoryLimits{}, nil
}

func (container *DockerContainer) LimitCPU(limits backend.CPULimits) error {
	panic("TODO LimitCPU")
	return nil
}

func (container *DockerContainer) CurrentCPULimits() (backend.CPULimits, error) {
	panic("TODO CurrentCPULimits")
	return backend.CPULimits{}, nil
}

func (container *DockerContainer) Run(spec backend.ProcessSpec) (uint32, <-chan backend.ProcessStream, error) {
	panic("TODO Run")
	return 0, nil, nil
}

func (container *DockerContainer) Attach(processID uint32) (<-chan backend.ProcessStream, error) {
	panic("TODO Attach")
	return nil, nil
}

func (container *DockerContainer) NetIn(hostPort uint32, containerPort uint32) (uint32, uint32, error) {
	panic("TODO NetIn")
	return hostPort, containerPort, nil
}

func (container *DockerContainer) NetOut(network string, port uint32) error {
	panic("TODO NetOut")
	return nil
}
