package graph

type Service interface {
}

type Resolver struct {
	service Service
}

func NewResolver(service Service) *Resolver {
	return &Resolver{service}
}
