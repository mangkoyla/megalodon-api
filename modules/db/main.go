package database

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/tursodatabase/go-libsql"
)

var (
	dbConn *sql.DB
	dbName = "megalodon-local.db"
	dbPath string
	dbDir  string
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

	if dbPath == "" {
		dir, err := os.MkdirTemp("", "megalodon-*")
		if err != nil {
			panic(err)
		}

		// Need to clear db dir when program exit
		dbDir = dir
		dbPath = filepath.Join(dir, dbName)
	}

	db.client = db.connect()

	return &db
}

func (db *databaseStruct) connect() *sql.DB {
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, db.dbURL, libsql.WithAuthToken(db.dbToken))
	if err != nil {
		panic(err)
	}

	dbConn = sql.OpenDB(connector)

	return dbConn
}
