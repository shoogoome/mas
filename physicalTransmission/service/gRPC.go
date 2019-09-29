package service

import (
	"context"
	"mas/models"
	pb "mas/models/physicalTransmission"
	"mas/physicalTransmission"
)

// 发送上传
func GRPCUpload(client pb.PhysicalTransmissionClient, shard []byte, index int, fileInfo models.FileInfo, statusMap chan models.ShardsStatus) {
	// 上传物理数据
	_, except := client.Upload(context.Background(), &pb.ShardChuckDataInfo{
		FileData: shard,
		Metadata: &pb.ShardChuckMetaData{
			FileHash: fileInfo.FileHash,
			Index:    int64(index),
			Shard:    true,
		},
	})
	status := models.ShardsStatus{
		Ip:    fileInfo.StorageServerIp[len(fileInfo.StorageServerIp)-1],
		Index: index,
		Client: client,
	}
	if except != nil {
		status.Status = false
	} else {
		status.Status = true
	}
	statusMap <- status
}

// 删除分片
func GRPCDeleteChuck(serverIp []string, hash string) {
	clients, _, _ := physicalTransmission.NewAppointGrpcConnection(serverIp)

	// TODO: 后续添加日志 该接口的报错信息将放在日志
	for i := 0; i < len(serverIp); i++ {
		go func(index int) {
			client := <- clients
			_, _ = client.DeleteChuck(context.Background(), &pb.ShardChuckMetaData{
				FileHash: hash,
				Index: int64(index),
			})
		} (i)
	}
}

// 删除数据
func GRPCDeleteShard(serverIp []string, hash string) {
	clients, _, _ := physicalTransmission.NewAppointGrpcConnection(serverIp)

	// TODO: 后续添加日志 该接口的报错信息将放在日志
	for i := 0; i < len(serverIp); i++ {
		go func(index int) {
			client := <- clients
			_, _ = client.DeleteShard(context.Background(), &pb.ShardChuckMetaData{
				FileHash: hash,
				Index: int64(index),
			})
		} (i)
	}
}






