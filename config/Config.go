package config

type Config struct {
	AppType        string                 `yaml:"appType"`
	LogLevel       int                    `yaml:"logLevel"`
	LogFile        string                 `yaml:"logFile"`
	Databases      []DbConfig             `yaml:"databases"`
	Redis          []RedisConfig          `yaml:"redis"`
	ServicesName   string                 `yaml:"servies_name"`
	KafkaConsumers []*KafkaConsumerConfig `yaml:"kafkaConsumers"`
	HttpGin        struct {
		Port        int  `yaml:"port"`
		AppendPprof bool `yaml:"appendPprof"`
	} `yaml:"httpgin"`
}

type DbConfig struct {
	Name         string `yaml:"name"`
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Dbname       string `yaml:"dbname"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
}

type RedisConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"poolSize"`
}

type KafkaConsumerConfig struct {
	Name              string   `yaml:"name"`
	Brokers           []string `yaml:"brokers"`
	Topics            []string `yaml:"topics"`
	AutoOffsetReset   string   `yaml:"autoOffsetReset"`
	GroupID           string   `yaml:"groupID"`
	MaxWaitTime       int      `yaml:"maxWaitTime"`
	SessionTimeout    int      `yaml:"sessionTimeout"`
	HeartbeatInterval int      `yaml:"heartbeatInterval"`
}
