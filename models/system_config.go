package models


type SystemConfig struct {
	// MongoDB配置
	MongoDB MongoDBConfig `yaml:"mongodb" json:"mongodb"`
	// Redis配置
	Redis RedisConfig `yaml:"redis" json:"redis"`
	// RabbitMQ配置
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq" json:"rabbitmq"`
	// 系统配置
	Server ServerConfig
}

type MongoDBConfig struct {
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
	Host string `yaml:"host" json:"host"`
	Port int `yaml:"port" json:"port"`
	DBName string `yaml:"db_name" json:"db_name"`
}

type RedisConfig struct {
	Host string `yaml:"host" json:"host"`
	Password string `yaml:"password" json:"password"`
	Port int `yaml:"port" json:"port"`
}

type RabbitMQConfig struct {
	Host string `yaml:"host" json:"host"`
	Queue string `yaml:"queue" json:"queue"`
}

type ServerConfig struct {
	// 分片存储根目录
	FileRootPath string
	// 分块存储根目录
	FileTempPath string
	// 是否启动gzip
	Gzip bool
	// 分块上传最大大小
	ChuckMaxSize int64
	// 系统token
	Token string
	// hmac加密key
	Key string
	// 当前服务ip
	Server string
	// 最大分片重发次数
	Resend int
	// gRPC连接重试
	GrpcRetry int
	// gRPC服务端口
	GrpcPort string
}
