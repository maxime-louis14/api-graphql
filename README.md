# Guide Simple pour GraphQL

## Introduction

GraphQL est un langage de requête pour les APIs et un environnement d'exécution pour les exécuter avec vos données existantes. Il vous permet de demander exactement ce dont vous avez besoin et rien de plus.

## Structure de base d'une API GraphQL

Une API GraphQL est composée de :

- **Schéma** : Définit les types d'objets et les relations entre eux.
- **Query** : Permet de lire ou récupérer des valeurs.
- **Mutation** : Permet de créer, mettre à jour ou supprimer des valeurs.
- **Resolvers** : Fournissent les fonctions pour remplir les données de votre schéma.

## Exemples

### Schéma GraphQL

Voici un exemple de schéma pour un type `Dog` et les opérations associées :

```graphql
type Dog {
  _id: ID!
  name: String!
  isGoodBoi: Boolean!
}

type Query {
  dog(_id: ID!): Dog
  dogs: [Dog!]!
}

input CreateDogInput {
  name: String!
  isGoodBoi: Boolean!
}

type Mutation {
  createDog(input: CreateDogInput!): Dog
  createDogs(inputs: [CreateDogInput!]!): [Dog!]!
}
```

### Query : Récupérer les données

1. Récupérer un chien par son ID :

```graphql
query GetDog {
  dog(_id: "some-dog-id") {
    _id
    name
    isGoodBoi
  }
}
```

2. Récupérer tous les chiens :

```graphql
query GetDogs {
  dogs {
    _id
    name
    isGoodBoi
  }
}
```

### Mutation : Modifier les données

1. Créer un chien :

```graphql
mutation AddDog {
  createDog(input: { name: "Rex", isGoodBoi: true }) {
    _id
    name
    isGoodBoi
  }
}
```

2. Créer plusieurs chiens :

```graphql
mutation AddDogs {
  createDogs(inputs: [
    { name: "Buddy", isGoodBoi: true },
    { name: "Bella", isGoodBoi: false },
    { name: "Max", isGoodBoi: true }
  ]) {
    _id
    name
    isGoodBoi
  }
}
```

### Resolvers : Implémenter les fonctions

En Go, vous pouvez implémenter les résolveurs comme ceci :

```go
package graph

import (
    "context"
    "graphql-tuto/database"
    "graphql-tuto/graph/model"
)

var db = database.Connect()

func (r *mutationResolver) CreateDog(ctx context.Context, input *model.CreateDogInput) (*model.Dog, error) {
    return db.CreateDog(ctx, input)
}

func (r *mutationResolver) CreateDogs(ctx context.Context, inputs []*model.CreateDogInput) ([]*model.Dog, error) {
    var dogs []*model.Dog
    for _, input := range inputs {
        newDog, err := db.CreateDog(ctx, input)
        if err != nil {
            return nil, err
        }
        dogs = append(dogs, newDog)
    }
    return dogs, nil
}

func (r *queryResolver) Dog(ctx context.Context, id string) (*model.Dog, error) {
    return db.FindByID(ctx, id)
}

func (r *queryResolver) Dogs(ctx context.Context) ([]*model.Dog, error) {
    return db.All(ctx)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
```

### Exemple de Serveur GraphQL en Go

Pour configurer un serveur GraphQL en Go, vous pouvez utiliser le package `gqlgen` :

1. **Installation** :

```sh
go get -u github.com/99designs/gqlgen
```

2. **Initialisation** :

```sh
go run github.com/99designs/gqlgen init
```

3. **Configuration du Serveur** :

Voici un exemple de serveur simple :

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/playground"
    "graphql-tuto/graph"
    "graphql-tuto/graph/generated"
)

const defaultPort = "8080"

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = defaultPort
    }

    srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

    http.Handle("/", playground.Handler("GraphQL playground", "/query"))
    http.Handle("/query", srv)

    log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
```

### Lancer le Serveur

Pour démarrer votre serveur GraphQL, exécutez :

```sh
go run main.go
```

Accédez à `http://localhost:8080/` pour utiliser le playground GraphQL.

---

Ce guide fournit une introduction simple à GraphQL avec des exemples pratiques de requêtes, de mutations, et de résolveurs. Pour approfondir, consultez la [documentation officielle de GraphQL](https://graphql.org/learn/) et celle de [gqlgen](https://gqlgen.com/).