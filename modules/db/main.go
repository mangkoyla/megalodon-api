package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
)

var (
	dbConn *sql.DB
	dbName = "megalodon-local.db"
	dbPath string
	dbDir  string
)

type databaseStruct struct {
	connector *libsql.Connector
	client    *sql.DB
	dbURL     string
	dbToken   string
}

func Init() string {
	dir, err := os.MkdirTemp("", "megalodon-*")
	if err != nil {
		panic(err)
	}

	// Need to clear db dir when program exit
	dbDir = dir
	dbPath = filepath.Join(dir, dbName)

	return dbDir
}

func MakeDatabase() *databaseStruct {
	db := databaseStruct{
		dbURL:   os.Getenv("TURSO_DATABASE_URL"),
		dbToken: os.Getenv("TURSO_AUTH_TOKEN"),
	}

	if dbPath == "" {
		panic(errors.New("database not initialized"))
	}

	db.client = db.connect()

	return &db
}

func (db *databaseStruct) connect() *sql.DB {
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, db.dbURL,
		libsql.WithAuthToken(db.dbToken),
		libsql.WithSyncInterval(5*time.Minute),
	)
	if err != nil {
		panic(err)
	}

	db.connector = connector
	dbConn = sql.OpenDB(connector)

	return dbConn
}

func (db *databaseStruct) Close() {
	db.connector.Close()
	db.client.Close()
	os.RemoveAll(dbDir)
}
