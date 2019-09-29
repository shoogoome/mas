package http_err

func GzipFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5601
	ctx.Msg = "压缩文件失败"
	return ctx
}

func GunzipFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5602
	ctx.Msg = "解压缩文件失败"
	return ctx
}
