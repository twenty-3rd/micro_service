package main

import (
	"flag"
	"github.com/go-kit/kit/log"
	zipkin "github.com/openzipkin-contrib/zipkin-go-opentracing"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"micro_server/ch13-seckill/pkg/bootstrap"
	"os"
)

func main() {

	// 创建环境变量
	var (
		zipkinURL = flag.String("zipkin.url", "http://127.0.0.1:9411/api/v2/spans", "Zipkin server url")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var zipkinTracer *zipkin.Tracer
	{
		var (
			err            error
			userNoopTracer = *zipkinURL == ""
			reporter       = zipkinhttp.NewReporter(*zipkinURL)
		)
		defer reporter.Close()

		zipkinEndpoint, _ := zipkin.NewEndpoint(bootstrap.HttpConfig.Host, bootstrap.HttpConfig.Port)
		zipkinTracer, err = zipkin.NewTracer(
			reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(useNoopTracer),
		)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		if !useNoopTracer {
			logger.Log("tracer", "Zipkin", "type", "Native", "URL", *zipkinURL)
		}
	}
	register.Register()

	tags := map[string]string{
		"component": "gateway_server",
	}
}
