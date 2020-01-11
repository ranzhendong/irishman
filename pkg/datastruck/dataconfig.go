package datastruck

//config.yaml
type Config struct {
	Server   server   `yaml:"server"`
	Etcd     etcd     `yaml:"etcd"`
	Reserved []string `yaml:"reserved"`
}

type server struct {
	Bind         string `yaml:"bind"`
	ReadTimeout  int    `yaml:"readtimeout"`
	WriteTimeout int    `yaml:"writetimeout"`
}

type etcd struct {
	Host    string `yaml:"host"`
	Format  string `yaml:"format"`
	Timeout int    `yaml:"timeout"`
}
