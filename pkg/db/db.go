package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

type DB struct {
	*sqlx.DB
}

func New(config *Config) (*DB, error) {
	conn, err := sqlx.Open(config.Driver, config.DSN)
	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}
