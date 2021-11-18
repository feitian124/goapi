package db

import (
	"strings"

	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/drivers/mysql"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

type DB struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Schema  *Schema
}

func Open(url string) (*DB, error) {
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	parts := strings.Split(u.Short(), "/")
	if len(parts) < 2 {
		return nil, errors.Errorf("invalid url: parse %s -> %#v", url, u)
	}

	var opts []drivers.Option
	if u.Driver == "mysql" {
		values := u.Query()
		for k := range values {
			if k == "show_auto_increment" {
				opts = append(opts, mysql.ShowAutoIcrrement())
				values.Del(k)
			}
			if k == "hide_auto_increment" {
				opts = append(opts, mysql.HideAutoIcrrement())
				values.Del(k)
			}
		}
		u.RawQuery = values.Encode()
		url = u.String()
	}

	db, err := dburl.Open(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}

	d := &DB{
		Name: u.Driver,
		Schema: &Schema{
			Name: u.Driver,
			DB:   db,
		},
	}

	return d, nil
}
