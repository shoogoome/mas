package file

import (
	"fmt"
	"mas/dao"
	"mas/exception/http_err"
	"mas/models"
	"mas/physicalTransmission"
	"mas/physicalTransmission/service"
	"mas/utils/config"
	"mas/utils/gzipUtils"
	"mas/utils/rs"
	tokenUtils "mas/utils/token"
	"time"
)

// 生成token
func GenerateTokenService(tokenType int, hash string) (string, interface{}) {

	var fileToken = models.FileToken{
		Hash:       hash,
		TokenType:  tokenType,
		CreateTime: time.Now().Unix(),
		ExpireAt:   time.Now().Unix() + 86400,
	}
	token, except := tokenUtils.GenerateToken(fileToken)

	if token == "" {
		return "", except
	}
	return token, nil
}

// 保存文件
func SaveFile(ddbyte []byte, fileInfo models.FileInfo, hash string) interface{} {

	// gzip压缩
	if config.SystemConfig.Server.Gzip {
		ddbyte, _ = gzipUtils.GzipFile(ddbyte, fileInfo.FileName)
	}

	fmt.Println("[*] gzip success...")

	// server
	clients, ips, except := physicalTransmission.NewRandomGrpcConnection()
	if except != nil {
		return except
	}

	fmt.Println("[*] get grpc client success")

	// 文件数据切片
	encode := rs.NewEncoder(ddbyte)
	shards, except := encode.Encode()
	if except != nil {
		return except
	}

	fmt.Println("[*] data encode success")

	var statusMap = make(chan models.ShardsStatus, models.RsConfig.AllShards)
	// 数据分片发送至存储服务端
	for index, shard := range shards {
		ip := <-ips
		fileInfo.StorageServerIp = append(fileInfo.StorageServerIp, ip)
		client := <- clients
		go service.GRPCUpload(client, shard, index, hash, ip, statusMap)
	}

	fmt.Println("[*] send over")

	// 读取结果 如果有允许损坏分片数量之内的分片数量损坏时
	// 重新修复分片并再次上传

	// 分片计数
	count := models.RsConfig.AllShards
	// 允许单一分片重发次数
	resend := make([]int, count)
	for {
		fmt.Println("[*] accept message: ", count)
		// 读取分片传输数据
		var status = <- statusMap
		count -= 1
		if !status.Status {
			// 分片重发
			if resend[status.Index] < config.SystemConfig.Server.Resend {
				go service.GRPCUpload(status.Client,
					shards[status.Index],
					status.Index, hash, status.Ip, statusMap)
				resend[status.Index] += 1
			} else {
				// 重发次数超出设定 认定失败
				// 删除数据分片
				go service.GRPCDeleteShard(fileInfo.StorageServerIp, hash)
				return http_err.ResendOver()
			}
		}
		// 所有分片处理完毕
		if count <= 0 {
			break
		}
	}
	// 文件信息存入数据库
	fileInfo.Persistence = true
	except = dao.UpdateFileInfo(fileInfo)
	if except != nil {
		return except
	}

	fmt.Println("[*] over")

	return nil
}

