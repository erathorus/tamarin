package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"gitlab.com/lattetalk/lattetalk/config"
)

var DB *sql.DB

func init() {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		config.Config.DB.User,
		config.Config.DB.Password,
		config.Config.DB.DBName,
		config.Config.DB.SSLMode,
	)
	var err error
	if DB, err = sql.Open("postgres", dsn); err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}
}
