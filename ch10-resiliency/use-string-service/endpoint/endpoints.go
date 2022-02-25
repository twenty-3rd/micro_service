package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"micro_server/ch10-resiliency/use-string-service/service"
)

// UseStringEndpoints CalculateEndpoint define endpoint
type UseStringEndpoints struct {
	UseStringEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

type UseStringRequest struct {
	RequestType string `json:"request_type"`
	A           string `json:"a"`
	B           string `json:"b"`
}

type UseStringResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func MakeUseStringEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseStringRequest)
		var (
			res, a, b string
			opError   error
		)
		a = req.A
		b = req.B
		res, opError = svc.UseStringService(req.RequestType, a, b)
		//if opError != nil {
		//	opErrorString = opError.Error()
		//}
		return UseStringResponse{Result: res}, opError
	}
}

// HealthRequest 康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
