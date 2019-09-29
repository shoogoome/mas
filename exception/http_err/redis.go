package http_err

func RedisConnectExcept () (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5101
	ctx.Msg = "redis连接错误"
	return ctx
}

func RedisVerificationError () (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5102
	ctx.Msg = "redis验证错误"
	return ctx
}

