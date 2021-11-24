package graph

import (
	"github.com/feitian124/goapi/db/mysql"
	"github.com/feitian124/goapi/graph/model"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB     *mysql.DB
	tables []*model.Table
}
