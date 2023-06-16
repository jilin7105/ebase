package config

type Config struct {
	AppType  string `yaml:"appType"`
	LogLevel int    `yaml:"logLevel"`
	LogFile  string `yaml:"logFile"`
}
