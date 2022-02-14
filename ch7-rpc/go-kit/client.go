package main

import (
	"context"
	"flag"
	"fmt"
	service "micro_server/ch7-rpc/go-kit/string-service"
	"micro_server/ch7-rpc/pb"
	"time"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	ctx := context.Background()
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("grpc dial err:", err)
	}
	defer conn.Close()

	svr := NewStringClient(conn)
	result, err := svr.Concat(ctx, "A", "B")
	if err != nil {
		fmt.Println("check error", err.Error())
	}
	fmt.Println("result = ", result)
}

func NewStringClient(conn *grpc.ClientConn) service.Service {

	var ep = grpctransport.NewClient(conn,
		"pb.StringService",
		"Concat",
		DecodeStringRequest,
		EncodeStringResponse,
		pb.StringResponse{},
	).Endpoint()

	userEp := service.StringEndpoints{
		StringEndpoint: ep,
	}
	return userEp
}

func DecodeStringRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return r, nil
}

func EncodeStringResponse(_ context.Context, r interface{}) (interface{}, error) {
	return r, nil
}
