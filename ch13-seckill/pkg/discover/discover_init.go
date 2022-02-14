package discover

import (
	"micro_server/ch13-seckill/pkg/bootstrap"
)

var ConsulService DiscoveryClient

func Register() {
	// 实例创建失败，直接停止服务
	if ConsulService == nil {
		panic(0)
	}

	instanceId := bootstrap.DiscoveryConfig.InstanceId
	if instanceId == "" {

	}
}
