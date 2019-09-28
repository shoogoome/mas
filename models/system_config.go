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
	FileRootPath string
	FileTempPath string
	Gzip bool
	ServerIp []string
	SignalUrl string
	StorageUrl string
	StorageChuckUrl string
	Token string
	Key string
	Server string
	ServerNum int
	Resend int
	GrpcClientNumber int
}
