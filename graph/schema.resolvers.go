package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/tigql/tigql/db/mysql"
	"github.com/tigql/tigql/graph/generated"
	"github.com/tigql/tigql/nils"
)

func (r *queryResolver) Tables(ctx context.Context, pattern *string) ([]*mysql.TableInfo, error) {
	p := nils.String(pattern)
	ts, err := r.DB.Tables(p)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *queryResolver) Table(ctx context.Context, name string) (*mysql.Table, error) {
	tb, err := r.DB.Table(name)
	if err != nil {
		return nil, err
	}
	return tb, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
