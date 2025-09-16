package executor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/utils"
)

type KustomizeExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewKustomizeExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *KustomizeExecutor {
	return &KustomizeExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (k *KustomizeExecutor) ApplyKustomize(ctx context.Context, task config.Task) error {
	if task.KustomizeConfig == nil {
		return fmt.Errorf("kustomize config is nil for task %s", task.Name)
	}

	path := task.KustomizeConfig.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("kustomize path does not exist: %s", path)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := fmt.Sprintf("kubectl apply -k %s", absPath)
	_, err = k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to apply kustomize manifests: %w", err)
	}

	return nil
}

func (k *KustomizeExecutor) DeleteKustomize(ctx context.Context, task config.Task) error {
	if task.KustomizeConfig == nil {
		return fmt.Errorf("kustomize config is nil for task %s", task.Name)
	}

	path := task.KustomizeConfig.Path
	if _, err := os.Stat(path); os.IsNotExist(err) {
		k.logger.Warning("Kustomize path does not exist: %s, skipping deletion", path)
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	cmd := fmt.Sprintf("kubectl delete -k %s --ignore-not-found", absPath)
	_, err = k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to delete kustomize manifests: %w", err)
	}

	return nil
}

func (k *KustomizeExecutor) VerifyInstallation(ctx context.Context, task config.Task) error {
	if task.Namespace == "" {
		// Namespace not specified, skip verification
		return nil
	}

	// Check if namespace exists and has resources
	cmd := fmt.Sprintf("kubectl get all -n %s", task.Namespace)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("installation verification failed: %w", err)
	}

	return nil
}
