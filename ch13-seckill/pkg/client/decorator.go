package client

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
	"micro_server/ch13-seckill/pkg/discover"
	"micro_server/ch13-seckill/pkg/loadbalance"
	"strconv"
	"time"
)

var ErrPRCService = errors.New("没有rpc服务")

var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}

type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName     string
	logger          *log.Logger
	discoveryClient discover.DiscoveryClient
	loadBalance     loadbalance.LoadBalance
	after           []InvokerAfterFunc
	before          []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string, tracer opentracing.Tracer, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

	// 1.回调函数
	for _, fn := range manager.before {
		if err = fn(); err != nil {
			return err
		}
	}

	// 2. 使用Hystrix的Do方法构造断路器
	if err = hystrix.Do(hystrixName, func() error {
		// 3.服务发现
		instances := manager.discoveryClient.DiscoverServices(manager.serviceName, manager.logger)
		// 4.负载均衡
		if instance, err := manager.loadBalance.SelectService(instances); err == nil {
			// 5.获得RPC端口并发送RPC请求
			if instance.GrpcPort > 0 {
				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(),
					grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(genTracer(tracer), otgrpc.LogPayloads())), grpc.WithTimeout(1*time.Second)); err == nil {
					if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				return ErrRPCService
			}
		} else {
			return err
		}
		return nil
	}, func(e error) error {
		return e
	}); err != nil {
		return err
	} else {
		//  6.调用after回调函数
		for _, fn := range manager.after {
			if err = fn(); err != nil {
				return err
			}
		}
		return nil
	}
}
