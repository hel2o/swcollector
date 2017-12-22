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

	for _, ch := range chs {
		swCpu, ok := <-ch
		if !ok {
			continue
		}
		L = append(L, GaugeValueIp(time.Now().Unix(), swCpu.Ip, "switch.CpuUtilization", swCpu.CpuUtil))
	}
	endTime := time.Now()
	log.Printf("UpdateCpuUtilization complete. Process time %s.", endTime.Sub(startTime))

	return L
}

func cpuMetrics(ip string, ch chan SwCpu) {
	var swCpu SwCpu

	cpuUtili, err := sw.CpuUtilization(ip, g.Config().Switch.Community, g.Config().Switch.SnmpTimeout, g.Config().Switch.SnmpRetry)
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
