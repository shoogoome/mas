package http_err


func RabbitmqConnectionFail () (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5801
	ctx.Msg = "rabbitmq连接失败"
	return ctx
}

func RabbitmqBindFail () (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5802
	ctx.Msg = "rabbitmq绑定失败"
	return ctx
}