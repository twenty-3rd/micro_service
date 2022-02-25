package setup

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	conf "micro_server/ch13-seckill/pkg/config"
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
