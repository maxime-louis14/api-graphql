package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.46

import (
	"context"
	"graphql-tuto/database"
	"graphql-tuto/graph/model"
)

var db = database.Connect()

// CreateDog is the resolver for the createDog field.
func (r *mutationResolver) CreateDog(ctx context.Context, input *model.NewDog) (*model.Dog, error) {
	return db.Save(input)
}

// Dog is the resolver for the dog field.
func (r *queryResolver) Dog(ctx context.Context, id string) (*model.Dog, error) {
	return db.FindByID(id)
}

// Dogs is the resolver for the dogs field.
func (r *queryResolver) Dogs(ctx context.Context) ([]*model.Dog, error) {
	return db.ALL()

}

// Modifiez vos résolveurs pour inclure la nouvelle mutation createDogs :
// func (r *mutationResolver) CreateDogs(ctx context.Context, inputs []*model.CreateDogInput) ([]*model.Dog, error) {
//     var dogs []*model.Dog
//     for _, input := range inputs {
//         newDog, err := r.DB.CreateDog(ctx, input)
//         if err != nil {
//             return nil, err
//         }
//         dogs = append(dogs, newDog)
//     }
//     return dogs, nil
// }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
