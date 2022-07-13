package graph

import "kom.com/m/v2/src/kom.com/graph/model"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos map[string]*model.Todo
}
