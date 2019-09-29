package dao

import (
	"MAS/exception/http_err"
	"MAS/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"mas/utils/mongo"
)

// 从数据库获取文件信息
func GetFileInfo(hash string) (fileInfo *models.FileInfo) {

	collection := mongo.MongoConn.Collection("fileserver")
	record := collection.FindOne(context.Background(), &bson.D{{
		"hash", hash,
	}})

	err := record.Decode(&fileInfo)
	if err != nil {
		return nil
	}
	return fileInfo
}

// 插入文件数据
func SaveFileInfo(fileInfo models.FileInfo) interface{} {
	// 写入到数据库
	collection := mongo.MongoConn.Collection("fileserver")
	_, err := collection.InsertOne(
		context.Background(),
		//&bson.M {
		//	"hash": fileInfo.Hash,
		//	"size": fileInfo.Size,
		//	"server_ip": fileInfo.ServerIp,
		//}
		fileInfo,
	);
	if err != nil {
		return http_err.SaveFileInfoError(err)
	}
	return nil
}

// 更新文件数据
func UpdateFileInfo(fileInfo models.FileInfo) interface{} {
	// 写入到数据库
	collection := mongo.MongoConn.Collection("fileserver")
	_, err := collection.UpdateOne(
		context.Background(),
		&bson.D{
			{"hash", fileInfo.FileHash},
		},
		&bson.M {
			"hash": fileInfo.FileHash,
			"size": fileInfo.FileSize,
			"server_ip": fileInfo.StorageServerIp,
		},
	)
	if err != nil {
		return http_err.SaveFileInfoError(err)
	}
	return nil
}

// 查询文件是否存在
func SearchFile(hash string) bool {
	fileInfo := GetFileInfo(hash)
	if fileInfo != nil {
		if !fileInfo.Persistence {
			return false
		}
		return true
	}
	return false
}
