package funcs

import (
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwMem struct {
	Ip       string
	MemUtili int
	UseTime  int64
}

func MemMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	chs := make([]chan SwMem, len(AliveIp))
	for i, ip := range AliveIp {
		if ip != "" {
			chs[i] = make(chan SwMem)
			go memMetrics(ip, chs[i])
		}
	}
	var useTime = make(map[string]int64, len(chs))

	for _, ch := range chs {
		swMem, ok := <-ch
		if !ok {
			continue
		}
		useTime[swMem.Ip] = swMem.UseTime

		L = append(L, GaugeValueIp(time.Now().Unix(), swMem.Ip, "switch.MemUtilization", swMem.MemUtili))
	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)

	log.Printf("UpdateMemUtilization complete. Process time %s. Used max time is %s, Latency=%ds.", endTime.Sub(startTime), maxIp, maxUseTime)

	return L
}

func memMetrics(ip string, ch chan SwMem) {
	var startTime, endTime int64
	startTime = time.Now().Unix()
	var swMem SwMem

	memUtili, err := sw.MemUtilization(ip, g.Config().Switch.Community, 2000, g.Config().Switch.SnmpRetry)
	endTime = time.Now().Unix()
	swMem.UseTime = endTime - startTime

	if err != nil {
		if g.Config().Debug {
			log.Println(err)
		}
		close(ch)
		return
	}

	swMem.Ip = ip
	swMem.MemUtili = memUtili

	ch <- swMem

	return
}
