package db

import (
	"database/sql"
	"fmt"

	//postgresql lib
	_ "github.com/lib/pq"
)

var db *sql.DB

//Connect connects to the database
func Connect(host string, port int, user, password, dbname string, sslmode bool) error {

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	if sslmode {
		connStr += " sslmode=require"
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

//Close database connection
func Close() {
	db.Close()
}
