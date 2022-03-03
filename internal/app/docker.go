package app

type DockerCompose struct {
	Version  string                            `yaml:"version"`
	Services map[string]map[string]interface{} `yaml:"services"`
	Networks map[string]interface{}            `yaml:"networks,omitempty"`
	Volumes  map[string]interface{}            `yaml:"volumes,omitempty"`
}
