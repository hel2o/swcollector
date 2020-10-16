package http

import (
	"log"
	"net/http"
	"time"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
)

type IfInOutPDU struct {
	In  float64
	Out float64
}

func configApiRoutes() {
	http.HandleFunc("/api/ifstats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		r.ParseForm()
		ip := r.PostFormValue("ip")
		if ip == "" {
			return
		}

		ifStatsList, err := sw.ListIfStats(ip, g.GetCommunity(ip), g.Config().Switch.SnmpTimeout, []string{}, g.Config().Switch.SnmpRetry, g.Config().Switch.LimitCon, true, false, true, true, true, true, true, true)
		if err != nil {
			log.Println(err)
			return
		}
		s := map[string]interface{}{
			"data": ifStatsList,
		}
		RenderJson(w, s)
	})

	http.HandleFunc("/api/ifhcinout", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		r.ParseForm()
		ip := r.PostFormValue("ip")
		index := r.PostFormValue("index")
		if ip == "" || index == "" {
			return
		}
		var ifInOut IfInOutPDU
		var err error
		inOid := "1.3.6.1.2.1.31.1.1.1.6." + index
		outOid := "1.3.6.1.2.1.31.1.1.1.10." + index
		ifInOut.In, err = funcs.GetCustMetric(ip, inOid, g.Config().Switch.SnmpTimeout, g.Config().Switch.SnmpRetry)
		if err != nil {
			log.Println(err)
			return
		}
		ifInOut.Out, err = funcs.GetCustMetric(ip, outOid, g.Config().Switch.SnmpTimeout, g.Config().Switch.SnmpRetry)
		if err != nil {
			log.Println(err)
			return
		}
		s := map[string]interface{}{
			"data": ifInOut,
			"ts":   time.Now().Format("15:04:05"),
		}
		RenderJson(w, s)
	})

}
