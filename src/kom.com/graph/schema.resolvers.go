package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"kom.com/m/v2/src/kom.com/graph/generated"
	"kom.com/m/v2/src/kom.com/graph/model"
	"errors"
	"math/rand"
)

// CreateTodo is the resolver for the createTodo field.
func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text: input.Text,
		ID:   fmt.Sprintf("T%d", rand.Intn(100000000)),
		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}
	if r.todos == nil {
		r.todos = make(map[string]*model.Todo)
	}
	r.todos[todo.ID] = todo
	return todo, nil
}

// CloseTodo is the resolver for the closeTodo field.
func (r *mutationResolver) CloseTodo(ctx context.Context, id string) (*model.Todo, error) {
	td, ok := r.todos[id]
	if ok {
		td.Done = true
		return td, nil
	}
	return nil, errors.New("nicht gefunden")
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	values := []*model.Todo{}
	for _, value := range r.todos {
		values = append(values, value)
	}
	return values, nil
}

// User is the resolver for the user field.
func (r *todoResolver) User(ctx context.Context, obj *model.Todo) (*model.User, error) {
	return &model.User{ID: obj.User.ID, Name: "user " + obj.User.ID}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Todo returns generated.TodoResolver implementation.
func (r *Resolver) Todo() generated.TodoResolver { return &todoResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type todoResolver struct{ *Resolver }
