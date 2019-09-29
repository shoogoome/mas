package models

import "mas/models/physicalTransmission"

type RedisChucks struct {
	ChuckInfo map[string]string `json:"chuck_info"`
}

type ShardsStatus struct {
	// 存储ip地址
	Ip     string
	// 存储状态
	Status bool
	// 分片序号
	Index  int
	// gRPC客户端
	Client physicalTransmission.PhysicalTransmissionClient
}
