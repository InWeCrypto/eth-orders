package main

import (
	"flag"

	"github.com/dynamicgo/aliyunlog"
	"github.com/dynamicgo/config"
	"github.com/dynamicgo/slf4go"
	orders "github.com/inwecrypto/eth-orders"
	_ "github.com/lib/pq"
)

var logger = slf4go.Get("eth-orders")
var configpath = flag.String("conf", "./orders.json", "geth indexer config file path")

func main() {

	flag.Parse()

	conf, err := config.NewFromFile(*configpath)

	if err != nil {
		logger.ErrorF("load eth indexer config err , %s", err)
		return
	}

	factory, err := aliyunlog.NewAliyunBackend(conf)

	if err != nil {
		logger.ErrorF("create aliyun log backend err , %s", err)
		return
	}

	slf4go.Backend(factory)

	server, err := orders.NewAPIServer(conf)

	if err != nil {
		logger.ErrorF("load eth config err , %s", err)
		return
	}

	if err := server.Run(); err != nil {
		logger.ErrorF("run api server err , %s", err)
	}
}
