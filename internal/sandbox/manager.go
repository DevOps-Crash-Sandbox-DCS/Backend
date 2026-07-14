package sandbox

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSandboxAlreadyExists = errors.New("sandbox already exists")
	ErrSandboxNotFound      = errors.New("sandbox not found")
	ErrInvalidSandboxImage  = errors.New("invalid sandbox image")
	ErrInvalidCommand       = errors.New("invalid command")
)

type Manager struct {
	docker *DockerCLI
	repo   *Repository
}

func NewManager(docker *DockerCLI, repo *Repository) *Manager {
	return &Manager{
		docker: docker,
		repo:   repo,
	}
}

func (m *Manager) Start(ctx context.Context, req StartSandboxRequest) (*Sandbox, error) {
	if req.SessionID == uuid.Nil {
		return nil, ErrSandboxNotFound
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		return nil, ErrInvalidSandboxImage
	}

	existing, err := m.repo.GetBySessionID(ctx, req.SessionID)
	if err == nil {
		if existing.Status == SandboxStatusRunning && m.docker.ContainerExists(ctx, existing.ContainerName) {
			return existing, ErrSandboxAlreadyExists
		}

		if existing.Status == SandboxStatusRunning {
			_ = m.repo.MarkFailed(ctx, req.SessionID)
		}

		return m.recreateExisting(ctx, existing, req)
	}

	if !errors.Is(err, ErrSandboxRecordNotFound) {
		return nil, err
	}

	containerName := buildContainerName(req.SessionID)

	if err := m.docker.RunContainer(ctx, containerName, image); err != nil {
		return nil, err
	}

	sandbox := &Sandbox{
		SessionID:     req.SessionID,
		ScenarioID:    strings.TrimSpace(req.ScenarioID),
		ContainerName: containerName,
		Image:         image,
		Status:        SandboxStatusRunning,
	}

	return m.repo.Create(ctx, sandbox)
}

func (m *Manager) EnsureStarted(ctx context.Context, req StartSandboxRequest) (*Sandbox, error) {
	existing, err := m.repo.GetBySessionID(ctx, req.SessionID)
	if err == nil {
		if existing.Status == SandboxStatusRunning && m.docker.ContainerExists(ctx, existing.ContainerName) {
			_ = m.repo.Touch(ctx, req.SessionID)

			return existing, nil
		}

		if existing.Status == SandboxStatusRunning {
			_ = m.repo.MarkFailed(ctx, req.SessionID)
		}

		return m.recreateExisting(ctx, existing, req)
	}

	if !errors.Is(err, ErrSandboxRecordNotFound) {
		return nil, err
	}

	return m.Start(ctx, req)
}

func (m *Manager) Exec(ctx context.Context, req ExecCommandRequest) (*CommandResult, error) {
	command := strings.TrimSpace(req.Command)
	if command == "" {
		return nil, ErrInvalidCommand
	}

	sandbox, err := m.repo.GetBySessionID(ctx, req.SessionID)
	if err != nil {
		if errors.Is(err, ErrSandboxRecordNotFound) {
			return nil, ErrSandboxNotFound
		}

		return nil, err
	}

	if sandbox.Status != SandboxStatusRunning {
		return nil, ErrSandboxNotFound
	}

	if !m.docker.ContainerExists(ctx, sandbox.ContainerName) {
		_ = m.repo.MarkFailed(ctx, req.SessionID)

		return nil, ErrSandboxNotFound
	}

	result, err := m.docker.Exec(ctx, sandbox.ContainerName, command)
	if err != nil {
		return nil, err
	}

	_ = m.repo.Touch(ctx, req.SessionID)

	return result, nil
}

func (m *Manager) Stop(ctx context.Context, sessionID uuid.UUID) error {
	sandbox, err := m.repo.GetBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, ErrSandboxRecordNotFound) {
			return ErrSandboxNotFound
		}

		return err
	}

	if sandbox.Status == SandboxStatusRunning && m.docker.ContainerExists(ctx, sandbox.ContainerName) {
		if err := m.docker.RemoveContainer(ctx, sandbox.ContainerName); err != nil {
			_ = m.repo.MarkFailed(ctx, sessionID)

			return err
		}
	}

	return m.repo.MarkStopped(ctx, sessionID)
}

func (m *Manager) Get(ctx context.Context, sessionID uuid.UUID) (*Sandbox, error) {
	sandbox, err := m.repo.GetBySessionID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, ErrSandboxRecordNotFound) {
			return nil, ErrSandboxNotFound
		}

		return nil, err
	}

	return sandbox, nil
}

func (m *Manager) CleanupOldRunning(ctx context.Context, olderThan time.Duration) (*CleanupResult, error) {
	items, err := m.repo.ListRunningOlderThan(ctx, olderThan)
	if err != nil {
		return nil, err
	}

	result := &CleanupResult{}

	for _, item := range items {
		if m.docker.ContainerExists(ctx, item.ContainerName) {
			if err := m.docker.RemoveContainer(ctx, item.ContainerName); err == nil {
				result.RemovedContainers++
			}
		}

		if err := m.repo.MarkStopped(ctx, item.SessionID); err == nil {
			result.StoppedInDB++
		}
	}

	return result, nil
}

func (m *Manager) recreateExisting(
	ctx context.Context,
	existing *Sandbox,
	req StartSandboxRequest,
) (*Sandbox, error) {
	containerName := existing.ContainerName
	if strings.TrimSpace(containerName) == "" {
		containerName = buildContainerName(req.SessionID)
	}

	if m.docker.ContainerExists(ctx, containerName) {
		_ = m.docker.RemoveContainer(ctx, containerName)
	}

	image := strings.TrimSpace(req.Image)
	if image == "" {
		image = existing.Image
	}

	if err := m.docker.RunContainer(ctx, containerName, image); err != nil {
		_ = m.repo.MarkFailed(ctx, req.SessionID)

		return nil, err
	}

	// В текущей схеме у session_id UNIQUE, поэтому проще переиспользовать запись.
	// Обновляем статус/last_seen через MarkRunning.
	if err := m.repo.MarkRunning(ctx, req.SessionID); err != nil {
		return nil, err
	}

	return m.repo.GetBySessionID(ctx, req.SessionID)
}

func buildContainerName(sessionID uuid.UUID) string {
	return fmt.Sprintf("dcs-sandbox-%s", sessionID.String())
}
