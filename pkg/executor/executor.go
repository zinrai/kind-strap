package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/zinrai/kind-strap/pkg/utils"
)

type CommandExecutor struct {
	logger *utils.Logger
	dryRun bool
}

func NewCommandExecutor(logger *utils.Logger, dryRun bool) *CommandExecutor {
	return &CommandExecutor{
		logger: logger,
		dryRun: dryRun,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, command string) (string, error) {
	// Always show the command being executed
	if e.dryRun {
		e.logger.Command("[dry-run] %s", command)
		return "", nil
	}

	e.logger.Command("%s", command)

	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		if stderrStr != "" {
			// Show error output for debugging
			fmt.Printf("Error output:\n%s\n", stderrStr)
		}
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return stdoutStr, nil
}

func (e *CommandExecutor) CheckCommandExists(command string) bool {
	if e.dryRun {
		// In dry-run mode, assume commands exist
		return true
	}
	_, err := exec.LookPath(command)
	return err == nil
}

func (e *CommandExecutor) EnsureRequiredCommands() error {
	if e.dryRun {
		e.logger.Info("Skipping command checks in dry-run mode")
		return nil
	}

	requiredCommands := []string{"kubectl", "helm", "kind"}

	missing := []string{}
	for _, cmd := range requiredCommands {
		if !e.CheckCommandExists(cmd) {
			missing = append(missing, cmd)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("required commands not found: %s. Please install them first", strings.Join(missing, ", "))
	}

	return nil
}

func (e *CommandExecutor) CheckKindCluster(ctx context.Context) error {
	if e.dryRun {
		e.logger.Info("Skipping cluster check in dry-run mode")
		return nil
	}

	// Check if kind cluster exists
	output, err := e.Execute(ctx, "kind get clusters")
	if err != nil {
		return fmt.Errorf("failed to list kind clusters: %w", err)
	}

	if strings.TrimSpace(output) == "" {
		return fmt.Errorf("no kind clusters found. Please create a cluster first with: kind create cluster")
	}

	// Check kubectl connection
	_, err = e.Execute(ctx, "kubectl cluster-info")
	if err != nil {
		return fmt.Errorf("cannot connect to kubernetes cluster: %w", err)
	}

	e.logger.Info("Successfully connected to kind cluster")
	return nil
}

func (e *CommandExecutor) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
