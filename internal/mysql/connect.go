package mysql

import (
	"context"
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
		Timeout:              time.Second * 5,
		ReadTimeout:          time.Second * 5,
		WriteTimeout:         time.Second * 5,
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	for i := 0; i <= 2; i++ {
		conn, _ := sqlx.Open("nrmysql", cfg.FormatDSN())

		err := conn.PingContext(context.TODO())
		if err == nil {
			return conn, nil
		}

		fmt.Println("failed to ping, sleep 1 second and try again")
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("failed to connect to database after 3 attempts")
}
