package env

type ArgoCliConfig struct {
	Command        string `yaml:"command"`
	Server         string `yaml:"server"`
	AdditionalArgs string `yaml:"additional_args"`
}
