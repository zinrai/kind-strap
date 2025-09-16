package executor

import (
	"context"
	"fmt"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/utils"
)

type KubectlExecutor struct {
	cmdExecutor *CommandExecutor
	logger      *utils.Logger
}

func NewKubectlExecutor(cmdExecutor *CommandExecutor, logger *utils.Logger) *KubectlExecutor {
	return &KubectlExecutor{
		cmdExecutor: cmdExecutor,
		logger:      logger,
	}
}

func (k *KubectlExecutor) ApplyManifest(ctx context.Context, task config.Task) error {
	if task.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for task %s", task.Name)
	}

	cmd := fmt.Sprintf("kubectl apply -f %s", task.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to apply manifest: %w", err)
	}

	return nil
}

func (k *KubectlExecutor) DeleteManifest(ctx context.Context, task config.Task) error {
	if task.KubectlConfig == nil {
		return fmt.Errorf("kubectl config is nil for task %s", task.Name)
	}

	cmd := fmt.Sprintf("kubectl delete -f %s --ignore-not-found", task.KubectlConfig.ManifestFile)
	_, err := k.cmdExecutor.Execute(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to delete manifest: %w", err)
	}

	return nil
}

func (k *KubectlExecutor) VerifyInstallation(ctx context.Context, task config.Task) error {
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
