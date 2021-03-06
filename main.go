package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hel2o/swcollector/cron"
	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
	"github.com/hel2o/swcollector/http"
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
		hostcfg := g.Config().SwitchHosts.Hosts
		g.ParseHostConfig(hostcfg)
	}
	if g.Config().CustomMetrics.Enabled {
		custMetrics := g.Config().CustomMetrics.Template
		g.ParseCustConfig(custMetrics)
	}
	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	if *check {
		funcs.CheckCollector()
		os.Exit(0)
	}
	funcs.NewLastifMap()
	funcs.BuildMappers()

	cron.Collect()

	go http.Start()

	select {}

}
