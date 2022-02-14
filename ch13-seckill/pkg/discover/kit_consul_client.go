package discover

import (
	"fmt"
	"log"
	"strconv"

	"micro_server/ch13-seckill/pkg/common"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func (consulClient *DiscoveryClientInstance) Register(instanceId, svcHost, healthCheckUrl, svcPort string, svcName string, weight int, meta map[string]string, tags []string, logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcHost)

	fmt.Println(weight)
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    svcName,
		Address: svcHost,
		Port:    port,
		Meta:    meta,
		Tags:    tags,
		Weights: &api.AgentWeights{
			Passing: weight,
		},
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + svcHost + ":" + strconv.Itoa(port) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	// 2.发送服务注册到consul中
	err := consulClient.client.Register(serviceRegistration)
	if err != nil {
		logger.Println("注册服务失败")
		return false
	}
	logger.Println("注册服务成功")
	return true
}

func (consulClient *DiscoveryClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {

	// 构建包含服务实例 ID 的元数据结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	// 发送服务注销请求
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil {
		if logger != nil {
			logger.Println("注销服务失败")
		}
		return false
	}
	if logger != nil {
		logger.Println("注销服务成功!")
	}

	return true
}

func (consulClient *DiscoveryClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []*common.ServiceInstance {
	// 尝试从缓存中获取
	instanceList, ok := consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	}

	// 申请锁
	consulClient.mutex.Lock()
	defer consulClient.mutex.Unlock()

	// 再次检查是否已经缓存
	instanceList, ok = consulClient.instancesMap.Load(serviceName)
	if ok {
		return instanceList.([]*common.ServiceInstance)
	} else {
		// 协程执行注册consul
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil {
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok {
					return // 数据异常， 忽略
				}

				// 服务实例为0
				if len(v) == 0 {
					consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
				}

				var healthServices []*common.ServiceInstance
				for _, service := range v {
					// todo
					if service.Checks.AggregatedStatus() == api.HealthPassing {
						healthServices = append(healthServices, newServiceInstance(service.Service))
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}

	// 根据服务名请求实例列表
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		consulClient.instancesMap.Store(serviceName, []*common.ServiceInstance{})
		if logger != nil {
			logger.Println("发现服务错误")
		}
		return nil
	}

	instances := make([]*common.ServiceInstance, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = newServiceInstance(entries[i].Service)
	}
	consulClient.instancesMap.Store(serviceName, instances)
	return instances
}

func newServiceInstance(service *api.AgentService) *common.ServiceInstance {
	rpcPort := service.Port - 1
	if service.Meta != nil {
		if rpcPortString, ok := service.Meta["rpcPort"]; ok {
			rpcPort, _ = strconv.Atoi(rpcPortString)
		}
	}

	return &common.ServiceInstance{
		Host:     service.Address,
		Port:     service.Port,
		GrpcPort: rpcPort,
		Weight:   service.Weights.Passing,
	}
}
