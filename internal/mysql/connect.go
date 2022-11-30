package mysql

import (
	"fmt"
	"time"

	driver "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	_ "github.com/newrelic/go-agent/_integrations/nrmysql"
)

func Connect(username, password, host, db string) (*sqlx.DB, error) {
	cfg := driver.Config{
		User:                 username,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 host,
		DBName:               db,
		Loc:                  time.UTC,
		Timeout:              time.Second * 2,
		ReadTimeout:          time.Second * 2,
		WriteTimeout:         time.Second * 2,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	conn, err := sqlx.Open("nrmysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to db: %w", err)
	}

	return conn, nil
}
