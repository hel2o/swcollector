package funcs

import (
	"fmt"
	"log"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/g"
	"github.com/open-falcon/common/model"
)

type SwSlotStatus struct {
	Ip        string
	SlotValue []sw.SlotValue
	UseTime   time.Duration
}

func FanMetrics() (L []*model.MetricValue) {
	if len(AliveIp) > 0 {
		startTime := time.Now()
		chs := make([]chan SwSlotStatus, len(AliveIp))
		for i, ip := range AliveIp {
			if ip != "" {
				chs[i] = make(chan SwSlotStatus)
				go swGeneralMetrics(ip, chs[i], sw.FanStatus)
			}
		}
		var useTime = make(map[string]time.Duration, len(chs))

		for _, ch := range chs {
			swFan, ok := <-ch
			if !ok {
				continue
			}
			useTime[swFan.Ip] = swFan.UseTime
			L = append(L, GaugeValueIp(time.Now().Unix(), swFan.Ip, SwcollectorTakeSec, swFan.UseTime.Seconds(), "type=Fan"))
			for i, status := range swFan.SlotValue {
				L = append(L, GaugeValueIp(time.Now().Unix(), swFan.Ip, "switch.Fan", status.Value, fmt.Sprintf("slot=%d", i)))
			}
		}
		endTime := time.Now()
		maxIp, maxUseTime := findMaxUseTime(useTime)
		log.Printf("Update FanStatus complete. Process time %s. Used max time is %s, Latency=%s.", endTime.Sub(startTime), maxIp, maxUseTime.String())
	}
	return L
}

func swGeneralMetrics(ip string, ch chan SwSlotStatus, f func(ip, community string, timeout, retry int) ([]sw.SlotValue, error)) {
	var startTime time.Time
	startTime = time.Now()
	var slotStatus SwSlotStatus
	vs, err := f(ip, g.GetCommunity(ip), 5000, g.Config().Switch.SnmpRetry)
	slotStatus.UseTime = time.Since(startTime)
	if err != nil {
		if g.Config().Debug {
			log.Println(err)
		}
		close(ch)
		return
	}
	slotStatus.Ip = ip
	slotStatus.SlotValue = vs
	ch <- slotStatus
	return
}
