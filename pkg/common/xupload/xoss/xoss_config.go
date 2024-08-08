package xoss

type Config struct {
	Endpoint    string `yaml:"endpoint"`
	Domain      string `yaml:"domain"`
	Ak          string `yaml:"ak"`
	SK          string `yaml:"sk"`
	Bucket      string `yaml:"bucket"`
	Region      string `yaml:"region"`
	StoragePath string `yaml:"storage_path"`
}
