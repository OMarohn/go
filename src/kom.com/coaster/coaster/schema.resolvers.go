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

// DeleteCoaster is the resolver for the deleteCoaster field.
func (r *mutationResolver) DeleteCoaster(ctx context.Context, id *string) (*model.Coaster, error) {
	err := r.service.deleteCoaster(*id)
	if err != nil {
		return nil, err
	}
	return &model.Coaster{ID: *id}, nil
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

// CoasterByID is the resolver for the coasterById field.
func (r *queryResolver) CoasterByID(ctx context.Context, id *string) (*model.Coaster, error) {
	theCoaster, err := r.service.getCoaster(*id)
	if err != nil {
		return &model.Coaster{ID: *id}, err
	}
	mc := &model.Coaster{ID: theCoaster.ID, Name: theCoaster.Name, Manufacture: &theCoaster.Manufacture, Height: &theCoaster.Height}
	return mc, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
