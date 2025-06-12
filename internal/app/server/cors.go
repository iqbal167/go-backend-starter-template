package server

import (
	"github.com/rs/cors"
)

func (s *Http) cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   s.Option.Cors.AllowedOrigins,
		AllowedMethods:   s.Option.Cors.AllowedMethods,
		AllowedHeaders:   s.Option.Cors.AllowedHeaders,
		AllowCredentials: s.Option.Cors.AllowCredentials,
		MaxAge:           s.Option.Cors.MaxAge,
	})
}
