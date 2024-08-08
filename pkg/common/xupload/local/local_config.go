package local

var config *Config

type Config struct {
	Domain   string `yaml:"domain" validate:"required"`
	FilePath string `yaml:"file_path" validate:"required"`
}

func SetUpConfig(cfg *Config) {
	config = cfg
}
