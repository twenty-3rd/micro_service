package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"micro_server/ch10-resiliency/use-string-service/config"
	"micro_server/common/discover"
	"micro_server/common/loadbalance"
)

// Service constants
const (
	StringServiceCommandName = "String.string"
	StringService            = "string"
)

type UseStringService struct {
	// 服务发现客户端
	discoveryClient discover.DiscoveryClient
	// 负载均衡器
	loadbalance loadbalance.LoadBalance
}

var (
	ErrHystrixFallbackExecute = errors.New("hystrix fall back execute")
)

func NewUseStringService(client discover.DiscoveryClient, lb loadbalance.LoadBalance) Service {
	hystrix.ConfigureCommand(StringServiceCommandName, hystrix.CommandConfig{
		// 触发最低阈值为5
		RequestVloumeThreshold: 5,
	})
	return &UseStringService{
		discoveryClient: client,
		loadbalance:     lb,
	}
}

type StringResponse struct {
	Result string `json:"result"`
	Error  error  `json:"error"`
}

func (s UseStringService) UseStringService(operationType, a, b string) (string, error) {
	var operationResult string
	err := hystrix.Do(StringServiceCommandName, func() error {
		instances := s.discoveryClient.DiscoveryServices(StringService, config.Logger)
		instanceList := make([]*api.AgentService, len(instances))
		for i := 0; i < len(instances); i++ {
			instanceList[i] = instances[i].(*api.AgentService)
		}
		// 选取实例
		selectInstance, err := s.loadbalance.SelectService(instanceList)
		if err != nil {
			config.Logger.Println(err.Error())
			return err
		}
		requestUrl := url.URL{
			Scheme: "http",
			Host:   selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:   "/op/" + operationType + "/" + a + "/" + b,
		}
		config.Logger.Printf("current string-service ID is %s and address:port is %s:%s\n", selectInstance.ID, selectInstance.Address, strconv.Itoa(selectInstance.Port))
		resp, err := http.Post(requestUrl.String(), "", nil)
		if err != nil {
			return err
		}
		result := &StringResponse{}

		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return err
		} else if result.Error != nil {
			return result.Error
		}

		operationResult = result.Result
		return nil
	}, func(e error) error {
		return ErrHystrixFallbackExecute
	})
	return operationResult, err
}

// 服务健康状态检查
func (s UseStringService) HealthCheck() bool {
	return true
}
