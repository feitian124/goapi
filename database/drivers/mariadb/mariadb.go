package mariadb

import (
	"database/sql"

	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/drivers/mysql"
)

type Mariadb struct {
	mysql.Mysql
}

// New return new Mariadb
func New(db *sql.DB, opts ...drivers.Option) (*Mariadb, error) {
	m, err := mysql.New(db, opts...)
	if err != nil {
		return nil, err
	}
	m.EnableMariaMode()
	return &Mariadb{*m}, nil
}
