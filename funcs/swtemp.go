package funcs

import (
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwTemp struct {
	Ip      string
	Temp    uint64
	UseTime time.Duration
}

func TempMetrics() (L []*model.MetricValue) {
	if len(AliveIp) > 0 {
		startTime := time.Now()
		chs := make([]chan SwTemp, len(AliveIp))
		for i, ip := range AliveIp {
			if ip != "" {
				chs[i] = make(chan SwTemp)
				go tempMetrics(ip, chs[i])
			}
		}
		var useTime = make(map[string]time.Duration, len(chs))

		for _, ch := range chs {
			swTemp, ok := <-ch
			if !ok {
				continue
			}
			useTime[swTemp.Ip] = swTemp.UseTime
			L = append(L, GaugeValueIp(time.Now().Unix(), swTemp.Ip, "switch.Temperature", swTemp.Temp))
			L = append(L, GaugeValueIp(time.Now().Unix(), swTemp.Ip, SwcollectorTakeSec, swTemp.UseTime.Seconds(), "type=temperature"))
		}
		endTime := time.Now()
		maxIp, maxUseTime := findMaxUseTime(useTime)
		log.Printf("Update Temperature complete. Process time %s. Used max time is %s, Latency=%s.", endTime.Sub(startTime), maxIp, maxUseTime.String())
	}
	return L
}

func tempMetrics(ip string, ch chan SwTemp) {
	var swTemp SwTemp
	var startTime time.Time
	startTime = time.Now()

	temp, err := sw.Temperature(ip, g.GetCommunity(ip), 5000, g.Config().Switch.SnmpRetry)
	swTemp.UseTime = time.Since(startTime)

	if err != nil {
		if g.Config().Debug {
			log.Println(err)
		}
		close(ch)
		return
	}
	swTemp.Ip = ip
	swTemp.Temp = temp
	ch <- swTemp

	return
}
