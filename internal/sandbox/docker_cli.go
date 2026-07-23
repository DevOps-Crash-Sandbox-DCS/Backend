package sandbox

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type DockerCLI struct {
	timeout time.Duration
}

func NewDockerCLI(timeout time.Duration) *DockerCLI {
	return &DockerCLI{
		timeout: timeout,
	}
}

func (d *DockerCLI) RunContainer(
	ctx context.Context,
	containerName string,
	image string,
) error {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"run",
		"-d",
		"--name", containerName,

		"--memory", "256m",
		"--cpus", "0.5",
		"--pids-limit", "128",

		"--security-opt", "no-new-privileges",

		image,
		"sleep",
		"3600",
	}

	result, err := d.run(ctx, args...)
	if err != nil {
		return errors.New(strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (d *DockerCLI) ContainerExists(ctx context.Context, containerName string) bool {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"inspect",
		containerName,
	}

	_, err := d.run(ctx, args...)

	return err == nil
}

func (d *DockerCLI) Exec(
	ctx context.Context,
	containerName string,
	command string,
) (*CommandResult, error) {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"exec",
		containerName,
		"sh",
		"-lc",
		command,
	}

	result, err := d.run(ctx, args...)
	if err != nil {
		return result, nil
	}

	return result, nil
}

func (d *DockerCLI) RemoveContainer(ctx context.Context, containerName string) error {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"rm",
		"-f",
		containerName,
	}

	result, err := d.run(ctx, args...)
	if err != nil {
		return errors.New(strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (d *DockerCLI) InspectContainer(ctx context.Context, containerName string) error {
	ctx, cancel := context.WithTimeout(ctx, d.timeout)
	defer cancel()

	args := []string{
		"inspect",
		containerName,
	}

	result, err := d.run(ctx, args...)
	if err != nil {
		return errors.New(strings.TrimSpace(result.Stderr))
	}

	return nil
}

func (d *DockerCLI) run(ctx context.Context, args ...string) (*CommandResult, error) {
	cmd := exec.CommandContext(ctx, "docker", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		exitCode = 1

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	}

	result := &CommandResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}

	return result, err
}
