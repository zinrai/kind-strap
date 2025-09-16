package config

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

type Task struct {
	Name            string           `yaml:"name"`
	Type            string           `yaml:"type"` // "helm", "kubectl", or "kustomize"
	Namespace       string           `yaml:"namespace"`
	HelmConfig      *HelmConfig      `yaml:"helmConfig,omitempty"`
	KubectlConfig   *KubectlConfig   `yaml:"kubectlConfig,omitempty"`
	KustomizeConfig *KustomizeConfig `yaml:"kustomizeConfig,omitempty"`
}

type HelmConfig struct {
	Repo       RepoInfo `yaml:"repo"`
	Chart      string   `yaml:"chart"`
	Version    string   `yaml:"version,omitempty"`
	ValuesFile string   `yaml:"valuesFile,omitempty"`
}

type RepoInfo struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type KubectlConfig struct {
	ManifestFile string `yaml:"manifestFile"`
}

type KustomizeConfig struct {
	Path string `yaml:"path"`
}

func (t *Task) Validate() error {
	if t.Name == "" {
		return ErrEmptyTaskName
	}

	switch t.Type {
	case "helm":
		if t.HelmConfig == nil {
			return ErrMissingHelmConfig
		}
		if t.HelmConfig.Repo.Name == "" || t.HelmConfig.Repo.URL == "" {
			return ErrInvalidHelmRepo
		}
		if t.HelmConfig.Chart == "" {
			return ErrMissingHelmChart
		}
	case "kubectl":
		if t.KubectlConfig == nil {
			return ErrMissingKubectlConfig
		}
		if t.KubectlConfig.ManifestFile == "" {
			return ErrMissingManifestFile
		}
	case "kustomize":
		if t.KustomizeConfig == nil {
			return ErrMissingKustomizeConfig
		}
		if t.KustomizeConfig.Path == "" {
			return ErrMissingKustomizePath
		}
	default:
		return ErrUnsupportedTaskType
	}

	return nil
}
