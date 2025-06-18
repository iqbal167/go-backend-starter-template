package rest

import (
	"fmt"

	"go-backend-starter-template/internal/config"
	"go-backend-starter-template/internal/pkg/server"

	"net/http"
	"net/http/httptest"
	"testing"
)

func mockConfig() *config.Config {
	return &config.Config{
		App: config.AppConfig{
			Name:    "test",
			Version: "1.0.0",
		},
	}
}

func TestNewRouter(t *testing.T) {
	config := mockConfig()
	rest := New(config)
	routes := rest.Routes()

	srv := server.New(routes, &server.Option{})

	tests := []struct {
		name           string
		path           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Root path should return app info",
			path:           "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   fmt.Sprintf("app_name=%s version=%s", config.App.Name, config.App.Version),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			srv.Server.Handler.ServeHTTP(w, req)

			if status := w.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if tt.expectedBody != "" {
				if body := w.Body.String(); body != tt.expectedBody {
					t.Errorf("handler returned wrong body: got %v want %v",
						body, tt.expectedBody)
				}
			}
		})
	}
}
