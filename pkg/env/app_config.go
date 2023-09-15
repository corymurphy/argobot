package env

type AppConfig struct {
	PullRequestPreamble string `yaml:"pull_request_preamble"`
	DataDirectory       string `yaml:"data_directory"`
	PrivateKeyFilePath  string `yaml:"private_key_file_path"`
}
