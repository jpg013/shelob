package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// Conn is the connection handle for the database
	Conn *sql.DB
)

// OpenConnection opens a mysql database connection
func OpenConnection() {
	var err error
	// Open up our database connection.
	connStr := MakeConnectionString()
	fmt.Println("connection string")
	fmt.Println(connStr)
	Conn, err = sql.Open("mysql", connStr)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
}

// MakeConnectionString transforms database config into myslq connection string.
func MakeConnectionString() string {
	// Load env
	dbAddress := os.Getenv("DB_ADDRESS")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbProtocol := os.Getenv("DB_PROTOCOL")
	dbPort := os.Getenv("DB_PORT")

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s", dbUser, dbPassword, dbProtocol, dbAddress, dbPort, dbName)
}
