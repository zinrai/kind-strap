# kind-strap

Bootstrap your [kind](https://kind.sigs.k8s.io/) cluster with operators and manifests using a single YAML definition.

## Overview

`kind-strap` is a task runner that sets up Kubernetes resources on kind clusters. Define your tasks in a YAML file to install operators, apply manifests, and configure your development environment.

## Features

- **Single YAML file** defines all tasks
- **Three deployment methods**: Helm, kubectl, and Kustomize
- **Sequential execution** with reverse-order uninstall
- **Command transparency**: Shows all executed commands
- **Dry-run mode**: Preview changes before applying
- **No state management**: Stateless operation

## Use Cases

- **Local Development**: Set up consistent development environments
- **Testing**: Create reproducible test environments
- **Learning**: Experiment with Kubernetes operators
- **CI/CD**: Automate kind cluster setup in pipelines
- **Demos**: Prepare demo environments quickly

## Requirements

- [kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [helm](https://helm.sh/)
- [Docker](https://www.docker.com/)

## Quick Start

1. Install kind-strap:

```bash
$ go install github.com/zinrai/kind-strap@latest
```

2. Create your kind cluster:

```bash
$ kind create cluster
```

3. Apply tasks

Using the Kyverno example:

```bash
$ kind-strap apply -f examples/kyverno/kind-strap.yaml
```

4. Clean up

```bash
$ kind-strap delete -f examples/kyverno/kind-strap.yaml
```

## Task Types

kind-strap supports three task types:

- **helm**: Install Helm charts - [example](./examples/kyverno/kind-strap.yaml)
- **kubectl**: Apply Kubernetes manifests - [example](./examples/tekton/kind-strap.yaml)
- **kustomize**: Deploy using Kustomize - [example](./examples/nginx/kind-strap.yaml)

See the [examples](./examples) directory for complete configurations.

## Commands

### apply

Apply tasks defined in the configuration file:

```bash
# Default configuration file (kind-strap.yaml)
$ kind-strap apply

# Custom configuration file
$ kind-strap apply -f my-tasks.yaml

# Preview commands without executing
$ kind-strap apply -dry-run
```

### delete

Remove tasks in reverse order:

```bash
# Default configuration file
$ kind-strap delete

# Custom configuration file
$ kind-strap delete -f my-tasks.yaml

# Preview commands without executing
$ kind-strap delete -dry-run
```

### version

```bash
$ kind-strap version
```

### help

```bash
$ kind-strap help
$ kind-strap apply -h
$ kind-strap delete -h
```

## Flags

| Flag       | Description                     | Default           |
|------------|---------------------------------|-------------------|
| `-f`       | Configuration file path         | `kind-strap.yaml` |
| `-dry-run` | Show commands without executing | `false`           |

## Output Example

```
$ kind-strap apply -f examples/kyverno/kind-strap.yaml

[15:04:05] INFO: Reading configuration from examples/kyverno/kind-strap.yaml
[15:04:05] CMD: kind get clusters
[15:04:05] CMD: kubectl cluster-info
[15:04:05] INFO: Successfully connected to kind cluster
[15:04:05] INFO: Starting to apply 1 tasks
[15:04:05] INFO: [1/1] Applying task: kyverno (type: helm)
[15:04:05] CMD: helm repo list
[15:04:05] CMD: helm repo add kyverno https://kyverno.github.io/kyverno/
[15:04:06] CMD: helm repo update
[15:04:07] CMD: helm upgrade --install kyverno kyverno/kyverno --namespace kyverno --create-namespace --version 3.1.1 --values ./examples/kyverno/values.yaml
[15:04:15] CMD: helm status kyverno --namespace kyverno
[15:04:15] OK: Task 'kyverno' applied successfully
[15:04:15] OK: All tasks applied successfully!
```

## License

This project is licensed under the [MIT License](./LICENSE).
