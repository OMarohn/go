package coaster

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	service *CoasterService
}

func NewCoasterResolver(theService *CoasterService) Resolver {
	return Resolver{
		service: theService,
	}
}
