package http_err

func StorageServerInsufficient() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5701
	ctx.Msg = "存储服务数量少于最小分片数量"
	return ctx
}

