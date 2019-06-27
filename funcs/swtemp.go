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
	Temp    int
	UseTime int64
}

func TempMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	chs := make([]chan SwTemp, len(AliveIp))
	for i, ip := range AliveIp {
		if ip != "" {
			chs[i] = make(chan SwTemp)
			go tempMetrics(ip, chs[i])
		}
	}
	var useTime = make(map[string]int64, len(chs))

	for _, ch := range chs {
		swTemp, ok := <-ch
		if !ok {
			continue
		}
		useTime[swTemp.Ip] = swTemp.UseTime

		L = append(L, GaugeValueIp(time.Now().Unix(), swTemp.Ip, "switch.Temperature", swTemp.Temp))

	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)
	log.Printf("UpdateTemperature complete. Process time %s. Used max time is %s, Latency=%ds.", endTime.Sub(startTime), maxIp, maxUseTime)

	return L
}

func tempMetrics(ip string, ch chan SwTemp) {
	var swTemp SwTemp
	var startTime, endTime int64
	startTime = time.Now().Unix()
	temp, err := sw.Temperature(ip, g.Config().Switch.Community, 2000, g.Config().Switch.SnmpRetry)
	endTime = time.Now().Unix()
	swTemp.UseTime = endTime - startTime

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
