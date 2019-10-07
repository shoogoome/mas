package http_err

import "fmt"

func UnmarshalBodyError () (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5201
	ctx.Msg = "读取body数据错误"
	return ctx
}

func MarshalFail() (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5202
	ctx.Msg = "序列化失败"
	return ctx
}

func LengthIsNotAllow(ob string, min int, max int) (ctx MASExceptBase) {

	msg := fmt.Sprintf("%s参数长度", ob)
	if min != -1 {
		msg += fmt.Sprintf("不得小于%d", min)
	}
	if max != -1 {
		msg += fmt.Sprintf("不得大于%d", max)
	}

	ctx.Status = false
	ctx.Code = 5203
	ctx.Msg = msg
	return ctx
}

func LackParams(par string) (ctx MASExceptBase){
	ctx.Status = false
	ctx.Code = 5204
	ctx.Msg = fmt.Sprintf("%s 参数为必填参数", par)
	return ctx
}