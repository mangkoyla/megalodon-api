package servers

import (
	"database/sql"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
)

type serversTableStruct struct {
	client *sql.DB
}

func MakeServersTableClient() *serversTableStruct {
	db := database.MakeDatabase()
	return &serversTableStruct{
		client: db.GetClient(),
	}
}

func (sts *serversTableStruct) CreateTableSafe() error {
	query := `CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY,
		code STRING,
		domain STRING,
		ip STRING,
		country STRING,
		users_count INT,
		users_max INT
	);`

	_, err := sts.client.Exec(query)
	return err
}

func (dts *serversTableStruct) GetDomainByCode(code string) (*ServerStruct, error) {
	var domain ServerStruct

	row := dts.client.QueryRow("SELECT * FROM users WHERE code = ?;", code)

	err := row.Scan(
		&domain.ID,
		&domain.Code,
		&domain.Domain,
		&domain.IP,
		&domain.Country,
		&domain.UsersCount,
		&domain.UsersMax,
	)

	return &domain, err
}
