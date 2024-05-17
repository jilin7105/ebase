package config

type Config struct {
	AppType        string                 `yaml:"appType"`
	LogLevel       int                    `yaml:"logLevel"`
	LogFile        string                 `yaml:"logFile"`
	Databases      []DbConfig             `yaml:"databases"`
	Mongo          []MongoConfig          `yaml:"mongoConfig"`
	Es             []EsConfig             `yaml:"es"`
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
	LinkTrack LinkTracking `yaml:"linkTracking"`
}

type LinkTracking struct {
	IsOpen            bool   `yaml:"is_open"`             //is_open 是否开启链路追踪
	IsLog             bool   `yaml:"is_log"`              //is_log 是否开启日志
	KafkaProducerName string `yaml:"kafka_producer_name"` //kafka 生产者名称
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
	Type         string `yaml:"type"`
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

// es 配置信息
type EsConfig struct {
	Hosts   []string `yaml:"hosts"`
	User    string   `yaml:"user"`
	Pass    string   `yaml:"pass"`
	CloudId string   `yaml:"cloudId"`
	ApiKey  string   `yaml:"apiKey"`
	Type    string   `yaml:"type"`
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
}

type MongoConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
	Name     string `yaml:"name"`
}
