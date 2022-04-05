package crud

import "database/sql"

type repositoryImpl[T Entity, ID comparable] struct {
	Conn *sql.DB
}

func NewRepository[T Entity, ID comparable](Conn *sql.DB) Repository {
	return &repositoryImpl[T, ID]{Conn}
}

func (r *repositoryImpl[T, ID]) Save(entity T) (T, error) {
	r.Conn.Exec("")
	return entity, nil
}
