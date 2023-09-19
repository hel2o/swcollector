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
	MemUtili uint64
	UseTime  time.Duration
}

func MemMetrics() (L []*model.MetricValue) {
	if len(AliveIp) > 0 {
		startTime := time.Now()
		chs := make([]chan SwMem, len(AliveIp))
		for i, ip := range AliveIp {
			if ip != "" {
				chs[i] = make(chan SwMem)
				go memMetrics(ip, chs[i])
			}
		}
		var useTime = make(map[string]time.Duration, len(chs))

		for _, ch := range chs {
			swMem, ok := <-ch
			if !ok {
				continue
			}
			useTime[swMem.Ip] = swMem.UseTime
			L = append(L, GaugeValueIp(time.Now().Unix(), swMem.Ip, SwcollectorTakeSec, swMem.UseTime.Seconds(), "type=memory"))
			L = append(L, GaugeValueIp(time.Now().Unix(), swMem.Ip, "switch.MemUtilization", swMem.MemUtili))
		}
		endTime := time.Now()
		maxIp, maxUseTime := findMaxUseTime(useTime)
		log.Printf("Update MemUtilization complete. Process time %s. Used max time is %s, Latency=%s.", endTime.Sub(startTime), maxIp, maxUseTime.String())
	}
	return L
}

func memMetrics(ip string, ch chan SwMem) {
	var startTime time.Time
	startTime = time.Now()
	var swMem SwMem

	memUtili, err := sw.MemUtilization(ip, g.GetCommunity(ip), 2000, g.Config().Switch.SnmpRetry)
	swMem.UseTime = time.Since(startTime)

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
