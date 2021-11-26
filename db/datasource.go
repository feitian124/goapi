package db

import "fmt"

type Datasource struct {
	UserName string
	Passwd   string
	Host     string
	Port     int
	DBName   string
	Comment  string
}

// ConnectString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func (ds *Datasource) ConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", ds.UserName, ds.Passwd, ds.Host, ds.Port, ds.DBName)
}
