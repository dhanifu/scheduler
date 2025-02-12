package repository

import (
	"fmt"
	"go-scheduler/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func NewPostgreConnection(config *config.Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DatabasePostgresHost,
		config.DatabasePostgresPort,
		config.DatabasePostgresUsername,
		config.DatabasePostgresPassword,
		config.DatabasePostgresName,
	)

	db := sqlx.MustOpen("postgres", dsn)
	err := db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
