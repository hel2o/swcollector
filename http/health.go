package http

import (
	"github.com/hel2o/swcollector/g"
	"net/http"
)

func configHealthRoutes() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(g.VERSION))
	})
}
