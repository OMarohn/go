package coaster

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"
	"strings"

	"kom.com/m/v2/src/kom.com/graph/generated"
	"kom.com/m/v2/src/kom.com/graph/model"
)

// CreateCoaster is the resolver for the createCoaster field.
func (r *mutationResolver) CreateCoaster(ctx context.Context, input model.NewCoaster) (*model.Coaster, error) {
	var newCoaster Coaster = Coaster{
		Name:        input.Name,
		Manufacture: strings.Clone(*input.Manufacture),
		Height:      *input.Height,
	}
	newCoaster.ID = fmt.Sprintf("id%d", rand.Intn(99999999))
	err := r.service.createCoaster(newCoaster)
	if err != nil {
		panic(fmt.Errorf("not implemented"))
	}

	return &model.Coaster{ID: newCoaster.ID, Name: newCoaster.Name, Manufacture: &newCoaster.Manufacture, Height: &newCoaster.Height}, nil
}

// Coasters is the resolver for the coasters field.
func (r *queryResolver) Coasters(ctx context.Context) ([]*model.Coaster, error) {
	allCoaster := r.service.getCoasters()
	var ret []*model.Coaster = make([]*model.Coaster, len(allCoaster))
	for i, c := range allCoaster {
		m := strings.Clone(c.Manufacture) // kopieren -- so richtig verstanden hab ich das nicht warum das muss !
		h := c.Height
		ret[i] = &model.Coaster{ID: c.ID, Name: c.Name, Manufacture: &m, Height: &h}
	}
	return ret, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
