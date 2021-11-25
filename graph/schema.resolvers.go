package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/feitian124/goapi/db/mysql"
	"github.com/feitian124/goapi/graph/generated"
	"github.com/feitian124/goapi/graph/model"
	"github.com/feitian124/goapi/nils"
)

func (r *queryResolver) Tables(ctx context.Context, pattern *string) ([]*mysql.TableInfo, error) {
	p := nils.String(pattern)
	ts, err := r.DB.Tables(p)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *queryResolver) Table(ctx context.Context, name string) (*model.Table, error) {
	panic(fmt.Errorf("not implemented"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
