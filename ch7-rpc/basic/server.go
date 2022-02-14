package main

import (
	"log"
	service "micro_server/ch7-rpc/basic/string-service"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	// 1.实现服务
	stringService := new(service.StringService)
	// 2.注册rpc
	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("Register error: ", registerError)
	}
	rpc.HandleHTTP()
	server, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(server, nil)
}
