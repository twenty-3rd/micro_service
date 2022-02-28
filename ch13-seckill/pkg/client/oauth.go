package client

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"micro_server/ch13-seckill/pb"
	"micro_server/ch13-seckill/pkg/discover"
	"micro_server/ch13-seckill/pkg/loadbalance"
)

// OAuthClient 如果要给一个类定义方法， 请用接口， 如果要给一个类定义属性， 请用结构体覆盖这个类
type OAuthClient interface {
	CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
	tracer      opentracing.Tracer
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, tracer opentracing.Tracer, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse) // 复制并获取引用
	if err := impl.manager.DecoratorInvoke("/pb.OauthService/CheckToken", "token_check", tracer, ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}

func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance, tracer opentracing.Tracer) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "Oauth"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &OAuthClientImpl{
		manager: &DefaultClientManager{
			serviceName:     serviceName,
			loadBalance:     lb,
			discoveryClient: discover.ConsulService,
			logger:          discover.Logger,
		},
		serviceName: serviceName,
		loadBalance: lb,
		tracer:      tracer,
	}, nil
}
