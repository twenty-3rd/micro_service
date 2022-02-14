package loadbalance

import (
	"errors"
	"math/rand"

	"micro_server/ch13-seckill/pkg/common"
)

/*
1.接口
2.
*/

type LoadBalance interface {
	SelectService(service []*common.ServiceInstance) (*common.ServiceInstance, error)
}

type RandomLoadBalance struct {
}

// 随机负载均衡
func (LoadBalance *RandomLoadBalance) SelectService(services []*common.ServiceInstance) (*common.ServiceInstance, error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("服务实例不存在")
	}
	// todo rand, rand.Intn
	return services[rand.Intn(len(services))], nil
}

type WeightRoundRobinLoadBalance struct {
}

// 权重平滑负载均衡
func (LoadBalance *WeightRoundRobinLoadBalance) SelectService(services []*common.ServiceInstance) (best *common.ServiceInstance, err error) {
	if services == nil || len(services) == 0 {
		return nil, errors.New("服务实例不存在")
	}
	total := 0
	for i := 0; i < len(services); i++ {
		w := services[i]
		if w == nil {
			continue
		}
		w.CurWeight += w.Weight

		total += w.Weight
		if best == nil || w.CurWeight > best.CurWeight {
			best = w
		}
	}

	if best == nil {
		return nil, nil
	}

	best.CurWeight -= total
	return best, nil
}
