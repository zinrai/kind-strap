package executor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/utils"
)

type HelmExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewHelmExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *HelmExecutor {
	return &HelmExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (h *HelmExecutor) AddRepository(ctx context.Context, repo config.RepoInfo) error {
	if !h.repositoryExists(ctx, repo.Name) {
		_, err := h.cmdExecutor.Execute(ctx, fmt.Sprintf("helm repo add %s %s", repo.Name, repo.URL))
		if err != nil {
			return fmt.Errorf("failed to add helm repository: %w", err)
		}
	}

	// Update repositories
	_, err := h.cmdExecutor.Execute(ctx, "helm repo update")
	if err != nil {
		return fmt.Errorf("failed to update helm repositories: %w", err)
	}

	return nil
}

// repositoryExists reports whether a repo with the exact name is configured.
// Matching the NAME column (not a substring of the whole listing) avoids a
// false positive when the name appears inside another repo's URL.
func (h *HelmExecutor) repositoryExists(ctx context.Context, name string) bool {
	output, err := h.cmdExecutor.Execute(ctx, "helm repo list -o json")
	if err != nil {
		// helm exits non-zero when no repositories are configured yet.
		return false
	}

	var repos []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(output), &repos); err != nil {
		return false
	}

	for _, r := range repos {
		if r.Name == name {
			return true
		}
	}
	return false
}

func (h *HelmExecutor) InstallChart(ctx context.Context, task config.Task) error {
	if task.HelmConfig == nil {
		return fmt.Errorf("helm config is nil for task %s", task.Name)
	}

	// Add repository
	if err := h.AddRepository(ctx, task.HelmConfig.Repo); err != nil {
		return err
	}

	// Build helm upgrade --install command
	cmd := fmt.Sprintf("helm upgrade --install %s %s/%s --namespace %s --create-namespace",
		task.Name,
		task.HelmConfig.Repo.Name,
		task.HelmConfig.Chart,
		task.Namespace)

	if task.HelmConfig.Version != "" {
		cmd += fmt.Sprintf(" --version %s", task.HelmConfig.Version)
	}

	if task.HelmConfig.ValuesFile != "" {
		if !h.cmdExecutor.FileExists(task.HelmConfig.ValuesFile) {
			return fmt.Errorf("values file not found: %s", task.HelmConfig.ValuesFile)
		}
		cmd += fmt.Sprintf(" --values %s", task.HelmConfig.ValuesFile)
	}

	_, err := h.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to install helm chart: %w", err)
	}

	return nil
}

func (h *HelmExecutor) UninstallChart(ctx context.Context, task config.Task) error {
	cmd := fmt.Sprintf("helm uninstall %s --namespace %s", task.Name, task.Namespace)
	_, err := h.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to uninstall helm chart: %w", err)
	}

	return nil
}

func (h *HelmExecutor) VerifyInstallation(ctx context.Context, task config.Task) error {
	// Check helm release status
	releaseCmd := fmt.Sprintf("helm status %s --namespace %s", task.Name, task.Namespace)
	_, err := h.cmdExecutor.Execute(ctx, releaseCmd)
	if err != nil {
		return fmt.Errorf("helm release verification failed: %w", err)
	}

	return nil
}
