package env

type AppConfig struct {
	DataDirectory      string `yaml:"dataDirectory"`
	PrivateKeyFilePath string `yaml:"privateKeyFilePath"`
}
