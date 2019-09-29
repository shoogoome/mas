package http_err

func TokenVerificationFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5301
	ctx.Msg = "token验证失败"
	return ctx
}

func ModifyTokenFail() (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5302
	ctx.Msg = "修改token失败"
	return ctx
}

func SystemTokenVerificationFail() (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5303
	ctx.Msg = "系统token验证失败"
	return ctx
}
