package env

import "github.com/palantir/go-githubapp/githubapp"

type Config struct {
	Server HTTPConfig       `yaml:"server"`
	Github githubapp.Config `yaml:"github"`

	ArgoCdUrl          string `yaml:"argoCdUrl"`
	PrivateKeyFilePath string `yaml:"privateKeyFilePath"`
}
