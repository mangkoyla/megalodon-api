package users

import (
	"database/sql"
	"time"

	database "github.com/FoolVPN-ID/megalodon-api/modules/db"
	"github.com/google/uuid"
	"github.com/sethvargo/go-password/password"
)

type usersTableStruct struct {
	client *sql.DB
}

func MakeUsersTableClient() *usersTableStruct {
	db := database.MakeDatabase()
	return &usersTableStruct{
		client: db.GetClient(),
	}
}

func (uts *usersTableStruct) CreateTableSafe() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		token STRING,
		password STRING,
		expired DATE,
		server_code STRING,
		quota INT,
		relay STRING,
		adblock INT2,
		vpn STRING
	);`

	_, err := uts.client.Exec(query)
	return err
}

func (uts *usersTableStruct) NewUser(id uint64) error {
	pass := password.MustGenerate(8, 2, 0, false, false)
	stmt, err := uts.client.Prepare(`INSERT INTO users VALUES (
		?,
		?,
		?,
		?,
		'',
		1000,
		'',
		1,
		''
	);`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, pass, uuid.New().String(), time.Now().Format("2006-01-02"))
	return err
}

func (uts *usersTableStruct) GetUserByIdOrToken(id, token any) (*UserStruct, error) {
	var (
		user UserStruct
	)

	row := uts.client.QueryRow("SELECT * FROM users WHERE id = ? OR token = ?;", id, token)

	err := row.Scan(
		&user.ID,
		&user.Token,
		&user.Password,
		&user.Expired,
		&user.ServerCode,
		&user.Quota,
		&user.Relay,
		&user.Adblock,
		&user.VPN,
	)

	return &user, err
}
