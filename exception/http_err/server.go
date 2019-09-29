package http_err

import "fmt"

func GetEnvKeyFail() (ctx MASExceptBase)  {
	ctx.Status = false
	ctx.Code = 5001
	ctx.Msg = "读取key失败"
	return ctx
}

func ServerNumLess(num int) (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5002
	ctx.Msg = fmt.Sprintf("正常运行系统数量不达标, 当前运行系统数量: %d", num)
	return ctx
}