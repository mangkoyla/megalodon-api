package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var (
	dbConn *sql.DB
)

type databaseStruct struct {
	client  *sql.DB
	dbURL   string
	dbToken string
}

func MakeDatabase() *databaseStruct {
	db := databaseStruct{
		dbURL:   os.Getenv("TURSO_DATABASE_URL"),
		dbToken: os.Getenv("TURSO_AUTH_TOKEN"),
	}

	db.client = db.connect()

	return &db
}

func (db *databaseStruct) connect() *sql.DB {
	if dbConn == nil {
		dbConn, _ = sql.Open("libsql", fmt.Sprintf("%s?authToken=%s", db.dbURL, db.dbToken))
	} else {
		dbConn.Ping()
	}

	return dbConn
}

func (db *databaseStruct) GetClient() *sql.DB {
	return db.client
}

func (db *databaseStruct) Close() {
	db.client.Close()
}
