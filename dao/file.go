package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mas/exception/http_err"
	"mas/models"
	"mas/utils/mongo"
	paramsUtils "mas/utils/params"
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
		fileInfo,
	)
	if err != nil {
		return http_err.SaveFileInfoError(err)
	}
	return nil
}

// 更新文件数据
func UpdateFileInfo(fileInfo models.FileInfo) interface{} {
	// 写入到数据库
	collection := mongo.MongoConn.Collection("fileserver")
	_, err := collection.UpdateMany(
		context.Background(),
		&bson.D{
			{"hash", fileInfo.FileHash},
		},
		&bson.M{
			"$set": fileInfo,
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

// 获取文件信息列表
func FileList(params paramsUtils.ParamsParser) (fileInfo []models.FileInfo, except interface{}) {
	collection := mongo.MongoConn.Collection("fileserver")

	filter := bson.M{}
	if params.Has("FileName") {
		filter["FileName"], except = params.Str("FileName", "文件名称"); if except != nil {
			return nil, except
		}
	}
	if params.Has("Persistence") {
		filter["Persistence"], except = params.Bool("Persistence", "是否持久化"); if except != nil {
			return nil, except
		}
	}
	if params.Has("FileHash") {
		filter["FileHash"], except = params.Str("FileHash", "文件hash"); if except != nil {
			return nil, except
		}
	}

	cursor, err := collection.Find(
		context.Background(),
		filter,
		options.Find().SetSort(bson.M{"CreateTime": -1}),
		)
	if err != nil {
		return nil, http_err.GetFileListFail()
	}
	err = cursor.All(context.Background(), fileInfo)
	if err != nil {
		return nil, http_err.GetFileListFail()
	}
	return fileInfo, except
}




