package datastruck

type Config struct {
	Consul consul `yaml:"consul"`
}

type consul struct {
	Host   string `yaml:"host"`
	Format string `yaml:"format"`
}
