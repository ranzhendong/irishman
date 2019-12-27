package datastruck

type Config struct {
	Etcd etcd `yaml:"etcd"`
}

type etcd struct {
	Host   string `yaml:"host"`
	Format string `yaml:"format"`
}
