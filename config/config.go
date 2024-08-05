package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(localhost:3306)/apotek")
	if err != nil {
		return err
	}
	return nil
}
