package file

import (
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
func (this *FileSystemController) generateToken(tokenType int) {
	this.Verification()
	hash := this.GetString("hash")
	var fileToken = models.FileToken{
		Hash:       hash,
		TokenType:  tokenType,
		CreateTime: time.Now().Unix(),
		ExpireAt:   time.Now().Unix() + 86400,
	}
	token, except := tokenUtils.GenerateToken(fileToken)

	if token == "" {
		this.Exception(except)
	}
	this.ReturnJSON(map[string]string{
		"token": token,
	})
}

// 保存文件
func (this *FileSystemController) saveFile(ddbyte []byte, fileInfo models.FileInfo, hash string) {

	// gzip压缩
	if config.SystemConfig.Server.Gzip {
		ddbyte, _ = gzipUtils.GzipFile(ddbyte, fileInfo.FileName)
	}
	// server
	clients, ips, except := physicalTransmission.NewRandomGrpcConnection();
	if except != nil {
		this.Exception(except)
	}
	// 文件数据切片
	encode := rs.NewEncoder(ddbyte)
	shards, except := encode.Encode()
	if except != nil {
		this.Exception(except)
	}

	var statusMap = make(chan models.ShardsStatus, rs.RsConfig.AllShards)
	// 数据分片发送至存储服务端
	for index, shard := range shards {
		ip := <-ips
		fileInfo.StorageServerIp = append(fileInfo.StorageServerIp, ip)
		client := <-clients
		go service.GRPCUpload(client, shard, index, hash, ip, statusMap)
	}

	// 读取结果 如果有允许损坏分片数量之内的分片数量损坏时
	// 重新修复分片并再次上传

	// 分片计数
	count := rs.RsConfig.AllShards
	// 允许单一分片重发次数
	resend := make([]int, count)
	for {
		// 读取分片传输数据
		var status = <-statusMap
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
				this.Exception(http_err.ResendOver())
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
		this.Exception(except)
	}
}

