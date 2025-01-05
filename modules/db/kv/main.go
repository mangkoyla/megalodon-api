package kv

import (
	"database/sql"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
)

type kvTableStruct struct {
	client *sql.DB
}

func MakeKVTableClient() *kvTableStruct {
	db := database.MakeDatabase()
	return &kvTableStruct{
		client: db.GetClient(),
	}
}

func (kts *kvTableStruct) CreateTableSafe() error {
	query := `CREATE TABLE IF NOT EXISTS kv (
		id INTEGER PRIMARY KEY,
		key STRING,
		value STRING
	);`

	_, err := kts.client.Exec(query)
	return err
}

func (kts *kvTableStruct) GetValueFromKVByKey(key string) (*string, error) {
	var value string
	row := kts.client.QueryRow("SELECT value FROM kv WHERE key = ?;", key)

	err := row.Scan(
		&value,
	)

	return &value, err
}
