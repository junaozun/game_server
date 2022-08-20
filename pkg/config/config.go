package config

// Config 配置
type Config struct {
	Server ServerConfig
	DB     DBConfig
	ETCD   EtcdConfig
}

type ServerConfig struct {
	Port  string
	Debug bool
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Debug    bool   `yaml:"debug"`
}

type EtcdConfig struct {
	Servers        string `yaml:"servers"`
	DialTimeout    int64  `yaml:"dial_timeout"`
	RequestTimeout int64  `yaml:"request_timeout"`
}
