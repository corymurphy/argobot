package env

import "github.com/palantir/go-githubapp/githubapp"

type Config struct {
	Server HTTPConfig       `yaml:"server"`
	Github githubapp.Config `yaml:"github"`

	ArgoCdApiUrl       string `yaml:"argoCdApiUrl"`
	ArgoCdWebUrl       string `yaml:"argoCdWebUrl"`
	PrivateKeyFilePath string `yaml:"privateKeyFilePath"`
	EnableLocking      bool   `yaml:"enableLocking"`
}
