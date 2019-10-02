package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	"mas/dao"
	"mas/exception/http_err"
	"mas/models"
	"mas/physicalTransmission"
	"mas/physicalTransmission/service"
	"mas/utils/config"
	"mas/utils/gzipUtils"
	"mas/utils/hashUtils"
	"mas/utils/rs"
	"mas/utils/token"
	"net/http"

	pb "mas/models/physicalTransmission"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

// 生成上传令牌
// header: systemToken
// get: hash
// @router /api/token/upload [get]
func (this *FileSystemController) GenerateUploadToken() {
	this.generateToken(tokenUtils.Upload)
}

// 生成下载令牌
// header: systemToken
// get: hash
// @router /api/token/download [get]
func (this *FileSystemController) GenerateDownloadToken() {
	this.generateToken(tokenUtils.Download)
}

// 获取文件信息
// header: systemToken
// get: hash
// @router /api/file/info [get]
func (this *FileSystemController) GetFileInfo() {
	this.Verification()
	hash := this.GetString("hash")
	this.ReturnJSON(dao.GetFileInfo(hash))
}

// 单文件上传
// header: token
// form-data: file: 数据文件
// @router /api/file/upload/single [post]
func (this *FileSystemController) UploadSingle() {
	hash := this.LoadHash(tokenUtils.Upload)
	// 查看文件是否存在
	if dao.SearchFile(hash) {
		this.ReturnJSON(map[string]string{
			"status": "success",
		})
		return
	}
	// 获取file
	file, headers, err := this.GetFile("file")
	if err != nil {
		this.Exception(http_err.GetFileFail())
	}
	// 计算真实文件hash
	var dd bytes.Buffer
	reader := io.TeeReader(file, &dd)
	fileHash, except := hashUtils.CalculateHash(reader)
	if except != nil {
		this.Exception(except)
	}
	// hash不匹配则报token不匹配错误
	if fileHash != hash {
		this.Exception(http_err.TokenFail())
	}
	// 构建文件基础信息
	fileInfo := models.FileInfo{
		FileName:    headers.Filename,
		CreateTime:  time.Now().Unix(),
		FileSize:    int64(dd.Len()),
		FileHash:    hash,
		Persistence: false,
	}
	// 保存文件
	dao.SaveFileInfo(fileInfo)
	this.saveFile(dd.Bytes(), fileInfo, hash)
	this.ReturnJSON(fileInfo)
}

// 初始化文件信息
// 给分片上传作准备，且在获取上传token之后进行
// header: token
// get: name 文件名称
// get: size 文件大小
// @router /api/file/upload/init [get]
func (this *FileSystemController) InitFileInfo() {

	hash := this.LoadHash(tokenUtils.Upload)

	size := this.GetString("size")
	sizeInt, except := strconv.Atoi(size)
	if except != nil {
		this.Exception(http_err.LackParams("文件长度[size]"))
	}
	this.ReturnJSON(models.FileInfo{
		FileHash:    hash,
		FileName:    this.GetString("name"),
		FileSize:    int64(sizeInt),
		Persistence: false,
	})
}

// 分片上传
// header: token
// get: chunk 分块序号 从1开始
// form-data: file 分块数据
// @router /api/file/upload/chuck [post]
func (this *FileSystemController) ChunkUpload() {

	hash := this.LoadHash(tokenUtils.Upload)
	chuck := this.GetString("chuck")

	// 查询文件大小是否符合要求
	file, h, err := this.GetFile("file")
	if err != nil {
		this.Exception(http_err.UploadFail())
	}
	if h.Size > config.SystemConfig.Server.ChuckMaxSize {
		this.Exception(http_err.ChuckSizeOverRegulations())
	}

	// 查询文件是否已初始化
	fileInfo := dao.GetFileInfo(hash)
	if fileInfo == nil {
		this.Exception(http_err.FileIsNotInit())
	}
	// 判断文件是否已持久化
	if fileInfo.Persistence {
		this.Exception(http_err.FileIsPersistence)
	}

	redisConn := this.RedisConn()
	defer redisConn.Close()
	chuckInfoString, err := redis.String(redisConn.Do("get", hash))
	// 查询是否已存储
	var chuckInfo models.RedisChucks
	if err == nil {
		err = json.Unmarshal([]byte(chuckInfoString), &chuckInfo)
		if err != nil {
			this.Exception(http_err.UploadFail())
		}
		_, ok := chuckInfo.ChuckInfo[chuck]
		if ok {
			this.Exception(http_err.ChuckExists())
		}
	} else {
		chuckInfo = models.RedisChucks{
			ChuckInfo: map[string]string{},
		}
	}

	fileName := path.Join(
		config.SystemConfig.Server.FileTempPath,
		fmt.Sprintf("%s.%s", hash, chuck),
	)
	chuckWrite, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		this.Exception(http_err.UploadFail())
	}
	defer chuckWrite.Close()
	// 存储分片数据
	_, err = io.Copy(chuckWrite, file)
	if err != nil {
		_ = os.Remove(fileName)
		this.Exception(http_err.UploadFail())
	}
	chuckInfo.ChuckInfo[chuck] = config.SystemConfig.Server.Server
	chuckInfoByte, _ := json.Marshal(chuckInfo)
	_, _ = redisConn.Do("set", hash, string(chuckInfoByte[:]))
	// 上传进度
	speed := strconv.Itoa(len(chuckInfo.ChuckInfo))
	this.ReturnJSON(map[string]string{
		"status": "success",
		"speed":  speed,
	})
}

// 完成上传
// header: systemToken
// header: token
// @router /api/file/upload/finish [get]
func (this *FileSystemController) Finish() {
	this.Verification()
	hash := this.LoadHash(tokenUtils.Upload)

	// 查询文件是否已初始化
	fileInfo := dao.GetFileInfo(hash)
	if fileInfo == nil {
		this.Exception(http_err.FileIsNotInit())
	}
	// 判断文件是否已持久化
	if fileInfo.Persistence {
		this.Exception(http_err.FileIsPersistence)
	}

	redisConn := this.RedisConn()
	defer redisConn.Close()
	chuckInfoString, err := redis.String(redisConn.Do("get", hash))
	if err != nil {
		this.Exception(http_err.UploadFail())
	}
	// 查询是否已存储
	var chuckInfo models.RedisChucks
	if len(chuckInfoString) == 0 {
		this.Exception(http_err.UploadFail())
	}
	err = json.Unmarshal([]byte(chuckInfoString), &chuckInfo)
	if err != nil {
		this.Exception(http_err.UploadFail())
	}
	// 读取文件
	chuckNum := len(chuckInfo.ChuckInfo)
	var mu sync.RWMutex
	lock := make(chan int)
	var chucks = make([][]byte, chuckNum+1)

	// 获取gRPC服务连接
	var clientMap = make(map[string]pb.PhysicalTransmissionClient)
	clients, ips, except := physicalTransmission.NewRandomGrpcConnection()
	if except != nil {
		this.Exception(except)
	}
	for i := 0; i < len(ips); i++ {
		clientMap[ <-ips] = <-clients
	}

	serverIPs := make([]string, 0)
	// 远端下载分块数据
	for chuck, ip := range chuckInfo.ChuckInfo {

		serverIPs = append(serverIPs, ip)
		go func(c string, nip string, lock chan int) {
			chuckInt, except := strconv.Atoi(chuck)
			if except != nil {
				lock <- 0
				return
			}
			// 远端下载
			client := clientMap[ip]
			shardChuckDataInfo, except := client.Download(context.Background(), &pb.ShardChuckMetaData{
				FileHash: hash,
				Index:    int64(chuckInt),
				Shard:    false,
			})
			if except != nil {
				lock <- 0
				return
			}
			mu.Lock()
			chucks[chuckInt] = shardChuckDataInfo.FileData
			mu.Unlock()
			lock <- 1
		}(chuck, ip, lock)
	}
	// 读取所有分片次数
	for i := 0; i < chuckNum; i++ {
		<-lock
	}
	// 合并所有分块
	fileByte := []byte("")
	allFileByte := bytes.Join(chucks, fileByte)
	// 计算真实文件hash
	var dd bytes.Buffer
	reader := io.TeeReader(bytes.NewBuffer(allFileByte), &dd)
	fileHash, except := hashUtils.CalculateHash(reader)
	if except != nil {
		this.Exception(except)
	}
	// hash不匹配则报token不匹配错误
	if fileHash != hash {
		// 删除分片
		service.GRPCDeleteChuck(serverIPs, hash)
		// 清除redis记录
		_, _ = redisConn.Do("del", hash)
		this.Exception(http_err.TokenFail())
	}
	ddByte := dd.Bytes()
	// 构建文件基础信息

	// 保存文件
	fileInfo.FileSize = int64(len(ddByte))
	this.saveFile(ddByte, *fileInfo, hash)
	// 删除临时分块
	service.GRPCDeleteChuck(serverIPs, hash)
	this.ReturnJSON(fileInfo)
}

// 文件下载
// header: token 下载令牌
// @router /api/file/upload/download [get]
func (this *FileSystemController) Download() {
	hash := this.LoadHash(tokenUtils.Download)

	// 查询文件信息
	fileInfo := dao.GetFileInfo(hash)
	if fileInfo == nil{
		this.Exception(http_err.FileIsNotExists())
	}
	if !fileInfo.Persistence {
		this.Exception(http_err.FileIsNotPersistence())
	}

	shards := make([][]byte, models.RsConfig.AllShards)
	// 获取分片数据
	var mu sync.RWMutex
	var lock = make(chan int)
	clients, _, except := physicalTransmission.NewAppointGrpcConnection(fileInfo.StorageServerIp)
	for index := range fileInfo.StorageServerIp {
		go func(index int, lock chan int) {
			client := <- clients
			shardChuckDataInfo, except := client.Download(context.Background(), &pb.ShardChuckMetaData{
				FileHash: hash,
				Index: int64(index),
				Shard: true,
			})
			if except != nil {
				lock <- 0
				return
			}
			mu.Lock()
			shards[index] = shardChuckDataInfo.FileData
			mu.Unlock()
			lock <- 1
		}(index, lock)
	}
	// 等待读取所有分片次数
	for i := 0; i < models.RsConfig.AllShards; i++ {
		<- lock
	}
	// 获取原文件
	var file io.ReadSeeker
	decode := rs.NewDecoder(shards, fileInfo.StorageServerIp)
	dd, except := decode.Decode(hash)
	if except != nil {
		service.GRPCDeleteShard(fileInfo.StorageServerIp, hash)
		this.Exception(except)
	}
	// gunzip
	if config.SystemConfig.Server.Gzip {
		dd, except = gzipUtils.GunzipFile(dd)
		if except != nil {
			service.GRPCDeleteShard(fileInfo.StorageServerIp, hash)
			this.Exception(except)
		}
	}
	file = bytes.NewReader(dd)
	// 输出文件
	this.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fileInfo.FileName)
	this.Ctx.Output.Header("Content-Length", fmt.Sprintf("%d", len(dd)))
	http.ServeContent(this.Ctx.Output.Context.ResponseWriter, this.Ctx.Output.Context.Request, fileInfo.FileName, time.Now(), file)
}

