package users

import (
	"database/sql"
	"errors"
	"fmt"
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
		ID INTEGER PRIMARY KEY,
		TOKEN TEXT,
		PASSWORD TEXT,
		EXPIRED DATE,
		DOMAIN_CODE TEXT,
		QUOTA INT,
		RELAY TEXT,
		ADBLOCK INT2,
		VPN TEXT
	);`

	_, err := uts.client.Exec(query)
	return err
}

func (uts *usersTableStruct) NewUser(id int) error {
	pass := password.MustGenerate(8, 2, 0, false, false)
	query := fmt.Sprintf(`INSERT INTO users VALUES (
		%d,
		'%s',
		'%s',
		'%s',
		'',
		1000,
		'',
		1,
		''
	);`, id, pass, uuid.New().String(), time.Now().Format("2006-01-02"))

	_, err := uts.client.Exec(query)
	return err
}

func (uts *usersTableStruct) GetUserByIdOrToken(id, token any) (*UserStruct, error) {
	var (
		user  UserStruct
		query string
	)

	if id != nil {
		query = fmt.Sprintf("SELECT * FROM users WHERE id = %d;", id.(uint64))
	} else if token != nil {
		query = fmt.Sprintf("SELECT * FROM users WHERE token = '%s';", token.(string))
	} else {
		return &user, errors.New("query failed")
	}

	row := uts.client.QueryRow(query)

	err := row.Scan(
		&user.ID,
		&user.Token,
		&user.Password,
		&user.Expired,
		&user.DomainCode,
		&user.Quota,
		&user.Relay,
		&user.Adblock,
		&user.VPN,
	)

	return &user, err
}
