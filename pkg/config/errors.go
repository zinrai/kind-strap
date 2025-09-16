package config

import "errors"

var (
	ErrEmptyTaskName          = errors.New("task name cannot be empty")
	ErrUnsupportedTaskType    = errors.New("unsupported task type, must be 'helm', 'kubectl', or 'kustomize'")
	ErrMissingHelmConfig      = errors.New("helmConfig is required for helm type tasks")
	ErrInvalidHelmRepo        = errors.New("helm repository name and URL must be specified")
	ErrMissingHelmChart       = errors.New("helm chart name must be specified")
	ErrMissingKubectlConfig   = errors.New("kubectlConfig is required for kubectl type tasks")
	ErrMissingManifestFile    = errors.New("manifestFile is required for kubectl type tasks")
	ErrMissingKustomizeConfig = errors.New("kustomizeConfig is required for kustomize type tasks")
	ErrMissingKustomizePath   = errors.New("path is required for kustomize type tasks")
	ErrFileNotFound           = errors.New("configuration file not found")
	ErrInvalidYAML            = errors.New("invalid YAML format")
)
