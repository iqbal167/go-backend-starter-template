package gorm

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB *gorm.DB
}

func New(conn *sql.DB) *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to open gorm session: %v", err)
	}

	return db
}
