package db

import (
	"strings"

	"github.com/feitian124/goapi/database/drivers"
	"github.com/feitian124/goapi/database/drivers/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/xo/dburl"
)

// DB stands for a connection to database, including current schema
type DB struct {
	Name    string  `json:"name"`
	Version string  `json:"version"`
	Url     string  `json:"url"`
	Schema  *Schema `json:"schema"`
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
			Name: parts[1],
			DB:   db,
		},
	}

	return d, nil
}

func (d *DB) Close() error {
	if d != nil && d.Schema != nil && d.Schema.DB != nil {
		if err := d.Schema.DB.Close(); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// UseSchema do nothing if the schema has been used, otherwise change the schema
func (d *DB) UseSchema(name string) (*Schema, error) {
	if d.Schema != nil && d.Schema.Name == name {
		return d.Schema, nil
	}

	if err := d.Close(); err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO wip, update d.Url with param name
	s, err := Open(d.Url)
	if err != nil {
		return nil, err
	}
	return s.Schema, nil
}
