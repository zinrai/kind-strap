package tasks

import (
	"context"
	"fmt"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/executor"
	"github.com/zinrai/kind-strap/pkg/utils"
)

func Delete(ctx context.Context, cfg *config.Config, dryRun bool, logger *utils.Logger) error {
	cmdExecutor := executor.NewCommandExecutor(logger, dryRun)

	// Check prerequisites
	if err := cmdExecutor.EnsureRequiredCommands(); err != nil {
		return fmt.Errorf("prerequisite check failed: %w", err)
	}

	if err := cmdExecutor.CheckKindCluster(ctx); err != nil {
		return fmt.Errorf("kind cluster check failed: %w", err)
	}

	// Initialize executors
	helmExecutor := executor.NewHelmExecutor(cmdExecutor, logger)
	kubectlExecutor := executor.NewKubectlExecutor(cmdExecutor, logger)
	kustomizeExecutor := executor.NewKustomizeExecutor(cmdExecutor, logger)

	// Process tasks in reverse order
	totalTasks := len(cfg.Tasks)
	for i := totalTasks - 1; i >= 0; i-- {
		task := cfg.Tasks[i]
		logger.Info("[%d/%d] Deleting task: %s (type: %s)",
			totalTasks-i, totalTasks, task.Name, task.Type)

		var err error
		switch task.Type {
		case "helm":
			err = helmExecutor.UninstallChart(ctx, task)
		case "kubectl":
			err = kubectlExecutor.DeleteManifest(ctx, task)
		case "kustomize":
			err = kustomizeExecutor.DeleteKustomize(ctx, task)
		default:
			err = fmt.Errorf("unsupported task type: %s", task.Type)
		}

		if err != nil {
			if !dryRun {
				logger.Warning("Failed to delete task '%s': %v (continuing...)", task.Name, err)
			}
			continue
		}

		if dryRun {
			logger.Success("Task '%s' validated for deletion", task.Name)
		} else {
			logger.Success("Task '%s' deleted successfully", task.Name)
		}
	}

	return nil
}
