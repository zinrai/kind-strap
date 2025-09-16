package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/zinrai/kind-strap/pkg/config"
	"github.com/zinrai/kind-strap/pkg/tasks"
	"github.com/zinrai/kind-strap/pkg/utils"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Handle subcommands
	switch os.Args[1] {
	case "version":
		fmt.Printf("kind-strap version %s\n", version)
		os.Exit(0)

	case "help":
		printUsage()
		os.Exit(0)

	case "apply":
		runApplyCommand()

	case "delete":
		runDeleteCommand()

	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func runApplyCommand() {
	flagSet := flag.NewFlagSet("apply", flag.ExitOnError)
	configFile := flagSet.String("f", "kind-strap.yaml", "Configuration file path")
	dryRun := flagSet.Bool("dry-run", false, "Show commands without executing them")

	// Custom usage for apply command
	flagSet.Usage = func() {
		fmt.Println("Usage: kind-strap apply [flags]")
		fmt.Println("\nApply tasks defined in configuration file")
		fmt.Println("\nFlags:")
		flagSet.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  # Apply with default configuration")
		fmt.Println("  kind-strap apply")
		fmt.Println("\n  # Apply with custom configuration")
		fmt.Println("  kind-strap apply -f my-tasks.yaml")
		fmt.Println("\n  # Preview what will be executed")
		fmt.Println("  kind-strap apply -dry-run")
	}

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	logger := utils.NewLogger()

	if *dryRun {
		logger.Warning("DRY-RUN MODE: Commands will be shown but not executed")
		fmt.Println()
	}

	if err := runApply(*configFile, *dryRun, logger); err != nil {
		logger.Error("Failed to apply tasks: %v", err)
		os.Exit(1)
	}
}

func runDeleteCommand() {
	flagSet := flag.NewFlagSet("delete", flag.ExitOnError)
	configFile := flagSet.String("f", "kind-strap.yaml", "Configuration file path")
	dryRun := flagSet.Bool("dry-run", false, "Show commands without executing them")

	// Custom usage for delete command
	flagSet.Usage = func() {
		fmt.Println("Usage: kind-strap delete [flags]")
		fmt.Println("\nDelete tasks defined in configuration file")
		fmt.Println("\nFlags:")
		flagSet.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  # Delete with default configuration")
		fmt.Println("  kind-strap delete")
		fmt.Println("\n  # Delete with custom configuration")
		fmt.Println("  kind-strap delete -f my-tasks.yaml")
		fmt.Println("\n  # Preview what will be deleted")
		fmt.Println("  kind-strap delete -dry-run")
	}

	if err := flagSet.Parse(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	logger := utils.NewLogger()

	if *dryRun {
		logger.Warning("DRY-RUN MODE: Commands will be shown but not executed")
		fmt.Println()
	}

	if err := runDelete(*configFile, *dryRun, logger); err != nil {
		logger.Error("Failed to delete tasks: %v", err)
		os.Exit(1)
	}
}

func runApply(configFile string, dryRun bool, logger *utils.Logger) error {
	logger.Info("Reading configuration from %s", configFile)

	cfg, err := config.ReadFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logger.Info("Starting to apply %d tasks", len(cfg.Tasks))
	if err := tasks.Apply(ctx, cfg, dryRun, logger); err != nil {
		return err
	}

	if dryRun {
		logger.Info("Dry-run completed. No changes were made.")
	} else {
		logger.Success("All tasks applied successfully!")
	}
	return nil
}

func runDelete(configFile string, dryRun bool, logger *utils.Logger) error {
	logger.Info("Reading configuration from %s", configFile)

	cfg, err := config.ReadFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	logger.Info("Starting to delete %d tasks", len(cfg.Tasks))
	if err := tasks.Delete(ctx, cfg, dryRun, logger); err != nil {
		return err
	}

	if dryRun {
		logger.Info("Dry-run completed. No changes were made.")
	} else {
		logger.Success("All tasks deleted successfully!")
	}
	return nil
}

func printUsage() {
	fmt.Print(`kind-strap - Bootstrap your kind cluster with operators and manifests

Usage:
  kind-strap [command] [flags]

Available Commands:
  apply       Apply tasks defined in configuration file
  delete      Delete tasks defined in configuration file
  version     Show version information
  help        Show this help message

Flags:
  -f string       Configuration file path (default "kind-strap.yaml")
  -dry-run        Show commands without executing them

Examples:
  # Apply tasks from default configuration
  kind-strap apply

  # Apply tasks from custom configuration
  kind-strap apply -f my-tasks.yaml

  # Preview what will be executed
  kind-strap apply -dry-run

  # Delete tasks
  kind-strap delete

  # Show version
  kind-strap version

For more information, visit: https://github.com/zinrai/kind-strap
`)
}
