package main

import (
	"flag"
	pb "micro_server/ch7-rpc/stream-pb"
	string_service "micro_server/ch7-rpc/stream/string-service"
	"net"

	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringService)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
