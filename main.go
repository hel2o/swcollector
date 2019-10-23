package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hel2o/swcollector/cron"
	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
	"github.com/hel2o/swcollector/http"
	"github.com/hel2o/swcollector/rpc"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	check := flag.Bool("check", false, "check collector")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}
	g.ParseConfig(*cfg)
	if g.Config().SwitchHosts.Enabled {
		hostCfg := g.Config().SwitchHosts.Hosts
		g.ParseHostConfig(hostCfg)
	}
	if g.Config().CustomMetrics.Enabled {
		custMetrics := g.Config().CustomMetrics.Template
		g.ParseCustConfig(custMetrics)
	}
	g.ModifyRlimit()
	g.StartSSL()
	g.InitRootDir()
	g.InitLocalIps()
	g.InitLocalIp()
	rpc.InitRpcClients()

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}

	funcs.NewLastifMap()
	funcs.BuildMappers()

	cron.Collect()

	go http.Start()
	go rpc.RpcServerStart()
	select {}

}
