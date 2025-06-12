package rest

import (
	"go-backend-starter-template/internal/config"
	"net/http"
)

type Rest struct {
	Config *config.Config
}

func New(config *config.Config) *Rest {
	return &Rest{
		Config: config,
	}
}

func (rest *Rest) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("app_name=" + rest.Config.App.Name + " version=" + rest.Config.App.Version))
	})

	return mux
}
