package lxc_cgroups_manager

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/pivotal-cf-experimental/garden/command_runner"
)

type LXCCgroupsManager struct {
	containerID string

	runner command_runner.CommandRunner
}

func New(containerID string, runner command_runner.CommandRunner) *LXCCgroupsManager {
	return &LXCCgroupsManager{
		containerID: containerID,

		runner: runner,
	}
}

func (manager *LXCCgroupsManager) Set(subsystem, name, value string) error {
	get := &exec.Cmd{
		Path: "lxc-cgroup",
		Args: []string{"--name", manager.containerID, name, value},
	}

	return manager.runner.Run(get)
}

func (manager *LXCCgroupsManager) Get(subsystem, name string) (string, error) {
	out := new(bytes.Buffer)

	get := &exec.Cmd{
		Path:   "lxc-cgroup",
		Args:   []string{"--name", manager.containerID, name},
		Stdout: out,
	}

	err := manager.runner.Run(get)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(out.String(), "\n"), nil
}
