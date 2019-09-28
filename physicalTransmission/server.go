package mian

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	pb "mas/models/physicalTransmission"
	"mas/utils"
	"os"
	"path"
)

type server struct{}

// 上传保存数据
func (server) Upload(context context.Context, shardChuckDataInfo *pb.ShardChuckDataInfo) (*pb.ShardChuckMetaData, error) {

	fileName := path.Join(
		utils.SystemConfig.Server.FileRootPath,
		fmt.Sprintf("%s.%d", shardChuckDataInfo.Metadata.FileHash, shardChuckDataInfo.Metadata.Index),
	)
	fileWrite, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0664);
	if err != nil {
		return nil, err
	}
	defer fileWrite.Close()
	_, err = io.Copy(fileWrite, bytes.NewReader(shardChuckDataInfo.FileData));
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// 下载数据
func (server) Download(context context.Context, shardChuckMetaData *pb.ShardChuckMetaData) (*pb.ShardChuckDataInfo, error) {

	var filePath string
	if shardChuckMetaData.Shard {
		filePath = utils.SystemConfig.Server.FileRootPath
	} else {
		filePath = utils.SystemConfig.Server.FileTempPath
	}

	fileName := fmt.Sprintf("%s.%d", shardChuckMetaData.FileHash, shardChuckMetaData.Index)
	fileBytes := getFile(filePath, fileName)
	// TODO: 异常处理
	return &pb.ShardChuckDataInfo{
		FileData: fileBytes,
		Metadata: nil,
	}, nil

}

// 删除分片数据
func (server) DeleteShard(context context.Context, shardChuckMetaData *pb.ShardChuckMetaData) (*pb.ShardChuckMetaData, error) {
	return nil, deleteFile(utils.SystemConfig.Server.FileRootPath, shardChuckMetaData.FileHash, shardChuckMetaData.Index)
}

// 删除分块数据
func (server) DeleteChuck(context context.Context, shardChuckMetaData *pb.ShardChuckMetaData) (*pb.ShardChuckMetaData, error) {
	return nil, deleteFile(utils.SystemConfig.Server.FileTempPath, shardChuckMetaData.FileHash, shardChuckMetaData.Index)
}

// 指定路径获取文件
func getFile(fileRootPath string, fileName string) []byte {

	filePath := path.Join(
		fileRootPath,
		fileName,
	)
	fileBytes, err := ioutil.ReadFile(filePath);
	if err != nil {
		return nil
	}
	return fileBytes
}

// 删除指定路径文件
func deleteFile(fileRootPath string, name string, index int64) error {
	fileName := path.Join(
		fileRootPath,
		fmt.Sprintf("%s.%d", name, index))
	err := os.Remove(fileName)
	return err
}
