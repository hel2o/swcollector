package funcs

import (
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwPing struct {
	Ip   string
	Ping float64
}

func PingMetrics() (L []*model.MetricValue) {
	vpns := g.Config().Switch.VpnRange
	vpnAndAlive := AliveIp
	for _, vpn := range vpns {
		vpnAndAlive = append(vpnAndAlive, vpn)
	}
	chs := make([]chan SwPing, len(vpnAndAlive))
	for i, ip := range vpnAndAlive {
		if ip != "" {
			chs[i] = make(chan SwPing)
			go pingMetrics(ip, chs[i])
		}
	}

	for _, ch := range chs {
		swPing := <-ch
		if swPing.Ping == -1 {
			log.Println(swPing.Ip, swPing.Ping)
		}
		L = append(L, GaugeValueIp(time.Now().Unix(), swPing.Ip, "switch.Ping", swPing.Ping))
	}	
	return L
}

func pingMetrics(ip string, ch chan SwPing) {
	var swPing SwPing
	timeout := g.Config().Switch.PingTimeout
	retry := g.Config().Switch.PingRetry
	fastPingMode := g.Config().Switch.FastPingMode
	rtt, err := sw.PingRtt(ip, timeout, retry, fastPingMode)
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
