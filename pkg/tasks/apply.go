package tasks

import (
	"context"
	"fmt"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/executor"
	"github.com/zinrai/kind-strap/pkg/utils"
)

func Apply(ctx context.Context, cfg *config.Config, dryRun bool, logger *utils.Logger) error {
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

	totalTasks := len(cfg.Tasks)
	for i, task := range cfg.Tasks {
		logger.Info("[%d/%d] Applying task: %s (type: %s)",
			i+1, totalTasks, task.Name, task.Type)

		var err error
		switch task.Type {
		case "helm":
			err = helmExecutor.InstallChart(ctx, task)
			if err == nil && !dryRun {
				err = helmExecutor.VerifyInstallation(ctx, task)
			}
		case "kubectl":
			err = kubectlExecutor.ApplyManifest(ctx, task)
			if err == nil && !dryRun {
				err = kubectlExecutor.VerifyInstallation(ctx, task)
			}
		case "kustomize":
			err = kustomizeExecutor.ApplyKustomize(ctx, task)
			if err == nil && !dryRun {
				err = kustomizeExecutor.VerifyInstallation(ctx, task)
			}
		default:
			err = fmt.Errorf("unsupported task type: %s", task.Type)
		}

		if err != nil {
			return fmt.Errorf("failed to apply task '%s': %w", task.Name, err)
		}

		if dryRun {
			logger.Success("Task '%s' validated", task.Name)
		} else {
			logger.Success("Task '%s' applied successfully", task.Name)
		}
	}

	return nil
}
