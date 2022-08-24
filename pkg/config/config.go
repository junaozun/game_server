package config

// GameConfig 配置
type GameConfig struct {
	Logic  *ServerConfig `yaml:"logic"`
	Cross  *ServerConfig `yaml:"cross"`
	Pvp    *ServerConfig `yaml:"pvp"`
	Battle *ServerConfig `yaml:"battle"`
	Common *CommonConfig `yaml:"common"`
}

type ServerConfig struct {
	Port  string       `yaml:"port"`
	Debug bool         `yaml:"debug"`
	Mysql *MysqlConfig `yaml:"mysql"`
}

type MysqlConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

type CommonConfig struct {
	Cluster string       `yaml:"cluster"`
	Etcd    *EtcdConfig  `yaml:"etcd"`
	NATS    *NatsConfig  `yaml:"nats"`
	Kafka   *KafkaConfig `yaml:"kafka"`
}

type EtcdConfig struct {
	Servers        string `yaml:"servers"`
	DialTimeout    int64  `yaml:"dial_timeout"`
	RequestTimeout int64  `yaml:"request_timeout"`
}

type KafkaConfig struct {
	Broker     string `yaml:"broker"`
	Frequency  int    `yaml:"frequency"`
	MaxMessage int    `yaml:"max_message"`
}

type NatsConfig struct {
	Server         string `yaml:"server"`          // nats://127.0.0.1:4222,nats://127.0.0.1:4223
	RequestTimeout int32  `yaml:"request_timeout"` // 请求超时（秒）
	ReconnectWait  int64  `yaml:"reconnect_wait"`  // 重连间隔
	MaxReconnects  int32  `yaml:"max_reconnects"`  // 重连次数
}
