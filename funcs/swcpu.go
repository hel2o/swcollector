package funcs

import (
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwCpu struct {
	Ip      string
	CpuUtil int
	UseTime int64
}

func CpuMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	chs := make([]chan SwCpu, len(AliveIp))
	for i, ip := range AliveIp {
		if ip != "" {
			chs[i] = make(chan SwCpu)
			go cpuMetrics(ip, chs[i])
		}
	}
	var useTime = make(map[string]int64, len(chs))

	for _, ch := range chs {
		swCpu, ok := <-ch
		if !ok {
			continue
		}
		useTime[swCpu.Ip] = swCpu.UseTime

		L = append(L, GaugeValueIp(time.Now().Unix(), swCpu.Ip, "switch.CpuUtilization", swCpu.CpuUtil))
	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)
	log.Printf("UpdateCpuUtilization complete. Process time %s. Used max time is %s, Latency=%ds.", endTime.Sub(startTime), maxIp, maxUseTime)

	return L
}

func cpuMetrics(ip string, ch chan SwCpu) {
	var startTime, endTime int64
	startTime = time.Now().Unix()
	var swCpu SwCpu

	cpuUtili, err := sw.CpuUtilization(ip, g.GetCommunity(ip), 2000, g.Config().Switch.SnmpRetry)
	endTime = time.Now().Unix()
	swCpu.UseTime = endTime - startTime

	if err != nil {
		if g.Config().Debug {
			log.Println(err)
		}
		close(ch)
		return
	}

	swCpu.Ip = ip
	swCpu.CpuUtil = cpuUtili

	ch <- swCpu

	return
}
