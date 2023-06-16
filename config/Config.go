package config

type Config struct {
	AppType   string        `yaml:"appType"`
	LogLevel  int           `yaml:"logLevel"`
	LogFile   string        `yaml:"logFile"`
	Databases []DbConfig    `yaml:"databases"`
	Redis     []RedisConfig `yaml:"redis"`
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
