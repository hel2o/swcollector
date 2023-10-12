package http

import (
	"fmt"
	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
	"log"
	"net/http"
	"strings"
)

func configAdminRoutes() {
	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		if g.IsTrustable(r.RemoteAddr) {
			g.SetReloadType(true)
			log.Println("config will be reload in next interval")
			RenderDataJson(w, "reload type on")
		} else {
			w.Write([]byte("no privilege"))
		}
	})
	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		RenderDataJson(w, funcs.AliveIp)
	})
	http.HandleFunc("/vendor", func(w http.ResponseWriter, r *http.Request) {
		var s []string
		sw.VendorMap.Range(func(ip, vendor any) bool {
			s = append(s, fmt.Sprintf("%s: %s", ip, vendor))
			return true
		})
		w.Write([]byte(strings.Join(s, "\n")))
	})
}
