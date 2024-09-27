package env

type ArgoConfig struct {
	Command        string `yaml:"command"`
	Server         string `yaml:"server"`
	ApiBaseUrl     string `yaml:"apiBaseUrl"`
	AdditionalArgs string `yaml:"AdditionalArguments"`
}
