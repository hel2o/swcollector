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
	UseTime time.Duration
}

func PingMetrics() (L []*model.MetricValue) {
	startTime := time.Now()
	var ipRange []string
	for _, ipc := range g.Config().Switch.IpRange {
		ipRange = append(ipRange, ipc.Ip)
	}
	chs := make([]chan SwPing, len(ipRange))
	for i, ip := range ipRange {
		if ip != "" {
			chs[i] = make(chan SwPing)
			go pingMetrics(ip, chs[i])
		}
	}
	var useTime = make(map[string]time.Duration, len(chs))
	for _, ch := range chs {
		swPing := <-ch
		useTime[swPing.Ip] = swPing.UseTime

		if swPing.Ping == -1 {
			if g.Config().Debug {
				log.Println(swPing.Ip, swPing.Ping)
			}
		}
		L = append(L, GaugeValueIp(time.Now().Unix(), swPing.Ip, SwcollectorTakeSec, swPing.UseTime.Seconds(), "type=ping"))
		L = append(L, GaugeValueIp(time.Now().Unix(), swPing.Ip, "switch.Ping", swPing.Ping))
	}
	endTime := time.Now()
	maxIp, maxUseTime := findMaxUseTime(useTime)
	log.Printf("Update Ping complete. Process time %s. Used max time is %s, Latency=%s.", endTime.Sub(startTime), maxIp, maxUseTime.String())

	return L
}

func pingMetrics(ip string, ch chan SwPing) {
	var swPing SwPing
	var startTime time.Time

	startTime = time.Now()
	timeout := g.Config().Switch.PingTimeout
	retry := g.Config().Switch.PingRetry
	fastPingMode := g.Config().Switch.FastPingMode
	rtt, err := sw.PingRtt(ip, timeout, retry, fastPingMode)

	swPing.UseTime = time.Since(startTime)

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
