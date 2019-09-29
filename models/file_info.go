package models

type FileInfo struct {
	// 文件大小
	FileSize int64 `json:"size" bson:"size"`
	// 文件名称
	FileName string `json:"name" bson:"name"`
	// 文件hash
	FileHash string `json:"hash" bson:"hash"`
	// 文件是否已持久化
	Persistence bool `json:"persistence" bson:"persistence"`
	// 创建时间
	CreateTime int64 `json:"create_time" bson:"create_time"`
	// 存储服务ip地址
	StorageServerIp []string `json:"server_ip" yaml:"server_ip"`
}
