package cron

import (
	"strings"

	"github.com/hel2o/management-system/tools"

	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
	_ "github.com/hel2o/swcollector/http"
	"github.com/open-falcon/common/model"

	"log"
	"math"
	"time"
)

func Collect() {
	if !g.Config().Transfer.Enabled {
		return
	}

	if g.Config().Transfer.Addr == "" {
		return
	}

	for _, v := range funcs.Mappers {
		go collect(int64(v.Interval), v.Fs)
	}
}

func collect(sec int64, fns []func() []*model.MetricValue) {
	for {
		go MetricToTransfer(sec, fns)
		time.Sleep(time.Duration(sec) * time.Second)
	}
}

func MetricToTransfer(sec int64, fns []func() []*model.MetricValue) {
	var mvs []*model.MetricValue

	for _, fn := range fns {
		items := fn()
		if items == nil {
			continue
		}

		if len(items) == 0 {
			continue
		}

		for _, mv := range items {
			mvs = append(mvs, mv)
		}
	}

	startTime := time.Now()

	//分批次传给transfer和N9E
	n := 30000
	lenMvs := len(mvs)

	div := lenMvs / n
	mod := math.Mod(float64(lenMvs), float64(n))

	var mvsSend []*model.MetricValue
	for i := 1; i <= div+1; i++ {
		if i < div+1 {
			mvsSend = mvs[n*(i-1) : n*i]
		} else {
			mvsSend = mvs[n*(i-1) : (n*(i-1))+int(mod)]
		}
		time.Sleep(100 * time.Millisecond)
		if g.Config().Transfer.N9e {
			var n9eSend []*tools.N9eMetric
			for _, v := range mvsSend {
				var tagsMap = make(map[string]string)
				for _, tags := range strings.Split(v.Tags, ",") {
					t := strings.Split(tags, "=")
					if len(t) == 2 {
						tagsMap[t[0]] = t[1]
					}
				}
				var alias string
				if p := strings.LastIndex(g.HostConfig().Hosts[v.Endpoint], "-"); p > -1 {
					alias = g.HostConfig().Hosts[v.Endpoint][:p]
				} else {
					continue
				}
				tagsMap["alias"] = alias
				tagsMap["ident"] = v.Endpoint

				n9eSend = append(n9eSend, &tools.N9eMetric{
					Metric:       strings.Replace(v.Metric, ".", "_", -1),
					Timestamp:    v.Timestamp,
					ValueUnTyped: v.Value,
					Tags:         tagsMap,
				})
			}
			var startPushTime = time.Now()
			if err, r := tools.PushDataToN9E(n9eSend); err != nil {
				log.Println("Push data to N9E error", err, r)
			} else {
				log.Printf(
					"<= N9E <Total=%v, Success=%d, Fail=%d, Latency=%v, Message:%s>\n",
					len(n9eSend),
					r.Success, r.Fail,
					time.Now().Sub(startPushTime),
					r.Msg,
				)
			}
		}

	}

	endTime := time.Now()
	log.Println("INFO : Send metrics N9E running in the background. Process time :", endTime.Sub(startTime), "Send metrics :", len(mvs))
}
