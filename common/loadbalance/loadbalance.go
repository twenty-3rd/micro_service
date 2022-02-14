package loadbalance

import (
	"errors"
	"math/rand"

	"github.com/hashicorp/consul/api"
)

// 负载均衡器
type LoadBalance interface {
	SelectService(service []*api.AgentService) (*api.AgentService, error)
}

// 随机负载均衡
type RandomLoadBalance struct {
}

func (LoadBalance *RandomLoadBalance) SelectService(services []*api.AgentService) (*api.AgentService, error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("service instances are not exist")
	}
	return services[rand.Intn(len(services))], nil
}
