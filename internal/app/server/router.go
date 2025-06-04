package server

import (
	"net/http"
)

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app_name=" + s.config.App.Name + " version=" + s.config.App.Version))
	})

	return mux
}
