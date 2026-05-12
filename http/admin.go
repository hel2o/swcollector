package http

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hel2o/sw"
	"github.com/hel2o/swcollector/funcs"
	"github.com/hel2o/swcollector/g"
)

func configAdminRoutes() {
	http.HandleFunc("/config/reload", func(w http.ResponseWriter, r *http.Request) {
		g.SetReloadType(true)
		log.Println("config will be reload in next interval")
		RenderDataJson(w, "reload type on")
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

	http.HandleFunc("/reloadssl", func(w http.ResponseWriter, r *http.Request) {
		err := g.KPR.ReloadCert()
		if err != nil {
			RenderDataJson(w, map[string]string{"msg": err.Error()})
			return
		}
		RenderDataJson(w, map[string]string{"msg": "ok"})
	})
}
