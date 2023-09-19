package http

import (
	"github.com/hel2o/swcollector/g"
	"log"
	"net/http"
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
}
