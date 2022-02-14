package setup

import (
	"fmt"
	conf "github.com/longjoy/micro-go-book/ch13-seckill/pkg/config"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

func InitZk() {
	// todo ip可能有问题
	var hosts = []string("127.0.0.1")
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return
	}
	conf.Zk.ZkConn = Conn
	conf.Zk.SecProductKey = "/product"
}
