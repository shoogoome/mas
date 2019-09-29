package http_err

import (
	"fmt"
)

func GetFileFail() (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5401
	ctx.Msg = "获取文件失败"
	return ctx
}

func CalculateHashError() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5402
	ctx.Msg = "计算文件hash失败"
	return ctx
}

func TokenFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5403
	ctx.Msg = "token无效或过期，请重新获取token"
	return ctx
}

func DamageToRawData() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5404
	ctx.Msg = "原始文件损坏，请重新上传"
	return ctx
}

func UploadFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5405
	ctx.Msg = "上传失败"
	return ctx
}

func ChuckExists() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5405
	ctx.Msg = "分片存在"
	return ctx
}

func StorageUnexpectedTermination(err error) (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5406
	ctx.Msg = fmt.Sprintf("存储意外终止: %v", err)
	return ctx
}

func FileIsNotExists() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5407
	ctx.Msg = "文件不存在"
	return ctx
}

func ResendOver() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5408
	ctx.Msg = "分片重发超出设定"
	return ctx
}

func DownloadFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5409
	ctx.Msg = "下载失败"
	return ctx
}

func FileIsNotInit() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5410
	ctx.Msg = "文件未初始化"
	return ctx
}

func FileIsPersistence() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5411
	ctx.Msg = "文件已存在"
	return ctx
}

func ChuckSizeOverRegulations() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5412
	ctx.Msg = "上传分块大小大于规定最大分块上传大小"
	return ctx
}

func FileIsNotPersistence() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5413
	ctx.Msg = "文件未完成上传"
	return ctx
}