package funcs

import (
	"log"
	"time"

	"github.com/hel2o/swcollector/g"
	"github.com/hel2o/sw"
	"github.com/open-falcon/common/model"
)

type SwTemp struct {
	Ip   string
	Temp int
}

func TempMetrics() (L []*model.MetricValue) {

	chs := make([]chan SwTemp, len(AliveIp))
	for i, ip := range AliveIp {
		if ip != "" {
			chs[i] = make(chan SwTemp)
			go tempMetrics(ip, chs[i])
		}
	}

	for _, ch := range chs {
		swTemp, ok := <-ch
		if !ok {
			continue
		}
		L = append(L, GaugeValueIp(time.Now().Unix(), swTemp.Ip, "switch.Temperature", swTemp.Temp))

	}

	return L
}

func tempMetrics(ip string, ch chan SwTemp) {
	var swTemp SwTemp

	temp, err := sw.Temperature(ip, g.Config().Switch.Community, g.Config().Switch.SnmpTimeout, g.Config().Switch.SnmpRetry)
	if err != nil {
		log.Println(err)
		close(ch)
		return
	}

	swTemp.Ip = ip
	swTemp.Temp = temp
	ch <- swTemp

	return
}
