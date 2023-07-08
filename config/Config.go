package config

type Config struct {
	AppType        string                 `yaml:"appType"`
	LogLevel       int                    `yaml:"logLevel"`
	LogFile        string                 `yaml:"logFile"`
	Databases      []DbConfig             `yaml:"databases"`
	Redis          []RedisConfig          `yaml:"redis"`
	KafkaProducers []KafkaProducerConfig  `yaml:"kafkaProducers"`
	ServicesName   string                 `yaml:"servies_name"`
	KafkaConsumers []*KafkaConsumerConfig `yaml:"kafkaConsumers"`
	HttpGin        struct {
		Port               int  `yaml:"port"`
		AppendPprof        bool `yaml:"appendPprof"`
		IPConcurrencyLimit int  `yaml:"iPConcurrencyLimit"`
	} `yaml:"httpginServer"`
	GrpcServer struct {
		Port          int  `yaml:"port"`
		TraceTracking bool `yaml:"traceTracking"`
	} `yaml:"grpcServer"`
	Micro struct {
		IsReg          bool  `yaml:"is_reg"`           //  is_reg : true  #是否有服务注册
		IsHeartPush    bool  `yaml:"is_heart_push"`    //  is_heart_push : true  #是否心跳推送
		HeartPushSpeed int64 `yaml:"heart_push_speed"` //  is_reg : true  #是否有服务注册
	} `yaml:"micro"` //micro :
}

//  heart_push_speed : 5  #心跳推送速度 单位 秒

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

type KafkaProducerConfig struct {
	Name                 string   `yaml:"name"`
	Brokers              []string `yaml:"brokers"`
	Topic                string   `yaml:"topic"`
	Compression          string   `yaml:"compression"`
	Timeout              int      `yaml:"timeout"`
	BatchSize            int      `yaml:"batchSize"`
	BatchTime            int      `yaml:"batchTime"`
	WaitForAll           bool     `yaml:"waitForAll"`
	MaxRetries           int      `yaml:"maxRetries"`
	RetryBackoff         int      `yaml:"retryBackoff"`
	ReturnSuccesses      bool     `yaml:"returnSuccesses"`
	NewManualPartitioner bool     `yaml:"newManualPartitioner"`
}
