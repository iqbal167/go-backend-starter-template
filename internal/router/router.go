package router

import (
	"go-backend-starter-template/internal/config"
	"net/http"
)

func New(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app_name=" + cfg.App.Name + " version=" + cfg.App.Version))
	})

	return mux
}
