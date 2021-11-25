package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/feitian124/goapi/nils"

	"github.com/feitian124/goapi/graph/generated"
	"github.com/feitian124/goapi/graph/model"
)

func (r *queryResolver) Tables(ctx context.Context, pattern *string) ([]*model.TableInfo, error) {
	p := nils.String(pattern)
	ts, err := r.DB.Tables(p)
	if err != nil {
		return nil, err
	}
	var tbs []*model.TableInfo
	for _, t := range ts {
		tb := &model.TableInfo{
			Name:      t.Name,
			Type:      t.Type,
			Comment:   t.Comment,
			Def:       t.Def,
			CreatedAt: t.CreatedAt,
		}
		tbs = append(tbs, tb)
	}
	return tbs, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
