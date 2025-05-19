package router

import (
	"go-backend-starter-template/internal/config"
	"net/http"
)

func New(config *config.Config) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app_name=" + config.App.Name + " version=" + config.App.Version))
	})

	return mux
}
