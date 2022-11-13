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
	CpuUtil uint64
	UseTime time.Duration
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
	var useTime = make(map[string]time.Duration, len(chs))

	for _, ch := range chs {
		swCpu, ok := <-ch
		if !ok {
			continue
		}
		useTime[swCpu.Ip] = swCpu.UseTime
		L = append(L, GaugeValueIp(time.Now().Unix(), swCpu.Ip, SwcollectorTakeSec, swCpu.UseTime.Seconds(), "type=cpu"))
		L = append(L, GaugeValueIp(time.Now().Unix(), swCpu.Ip, "switch.CpuUtilization", swCpu.CpuUtil))
	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)
	log.Printf("Update CpuUtilization complete. Process time %s. Used max time is %s, Latency=%s.", endTime.Sub(startTime), maxIp, maxUseTime.String())
	return L
}

func cpuMetrics(ip string, ch chan SwCpu) {
	var startTime time.Time
	startTime = time.Now()
	var swCpu SwCpu
	cpuUtili, err := sw.CpuUtilization(ip, g.GetCommunity(ip), 5000, g.Config().Switch.SnmpRetry)
	swCpu.UseTime = time.Since(startTime)
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
