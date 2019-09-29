package http_err

import "fmt"

func SaveFileInfoError(err error) (ctx MASExceptBase) {
	ctx.Status = false
	ctx.Code = 5501
	ctx.Msg = fmt.Sprintf("保存文件信息错误: %v", err)
	return ctx
}
