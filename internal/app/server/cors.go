package server

import (
	"github.com/rs/cors"
)

func (s *Server) cors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   s.config.CORSConfig.AllowedOrigins,
		AllowedMethods:   s.config.CORSConfig.AllowedMethods,
		AllowedHeaders:   s.config.CORSConfig.AllowedHeaders,
		AllowCredentials: s.config.CORSConfig.AllowCredentials,
		MaxAge:           s.config.CORSConfig.MaxAge,
	})
}
