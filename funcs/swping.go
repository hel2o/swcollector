package funcs

import (
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwPing struct {
	Ip      string
	Ping    float64
	UseTime int64
}

func PingMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	vpns := g.Config().Switch.VpnRange
	ipRange := g.Config().Switch.IpRange
	allIp := append(ipRange, vpns...)
	chs := make([]chan SwPing, len(allIp))
	for i, ip := range allIp {
		if ip != "" {
			chs[i] = make(chan SwPing)
			go pingMetrics(ip, chs[i])
		}
	}
	var useTime = make(map[string]int64, len(chs))
	for _, ch := range chs {
		swPing := <-ch
		useTime[swPing.Ip] = swPing.UseTime

		if swPing.Ping == -1 {
			if g.Config().Debug {
				log.Println(swPing.Ip, swPing.Ping)
			}
		}
		L = append(L, GaugeValueIp(time.Now().Unix(), swPing.Ip, "switch.Ping", swPing.Ping))
	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)
	log.Printf("UpdatePing complete. Process time %s. Used max time is %s, Latency=%ds.", endTime.Sub(startTime), maxIp, maxUseTime)

	return L
}

func pingMetrics(ip string, ch chan SwPing) {
	var swPing SwPing
	var startTime, endTime int64

	startTime = time.Now().Unix()
	timeout := g.Config().Switch.PingTimeout
	retry := g.Config().Switch.PingRetry
	fastPingMode := g.Config().Switch.FastPingMode
	rtt, err := sw.PingRtt(ip, timeout, retry, fastPingMode)

	endTime = time.Now().Unix()
	swPing.UseTime = endTime - startTime

	if err != nil {
		log.Println(ip, err)
		swPing.Ip = ip
		swPing.Ping = -1
		ch <- swPing
		return
	}
	swPing.Ip = ip
	swPing.Ping = rtt
	ch <- swPing
	return

}
