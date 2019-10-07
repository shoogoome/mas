package system

import (
	"mas/physicalTransmission"
)

// 活跃信号反馈
// @router /api/server/signal [get]
func (this *SystemController) Signal() {
	this.Verification()
	// 返回当前服务ip
	this.ReturnJSON(map[string][]string {
		"ip": physicalTransmission.GetRandomServerIp(),
	})
}

