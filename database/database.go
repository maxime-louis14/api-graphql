package database

import (
	"context"
	"fmt"
	"graphql-tuto/graph/model"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatalf("Erreur lors de la connexion à MongoDB: %v", err)
	}

	// Vérifiez la connexion
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Impossible de se connecter à MongoDB: %v", err)
	}

	fmt.Println("Connexion à MongoDB réussie")
	return &DB{client: client}
}

func (db *DB) Save(input *model.NewDog) (*model.Dog, error) {
	collection := db.client.Database("animals").Collection("dogs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de l'insertion du chien: %w", err)
	}

	return &model.Dog{
		ID:       res.InsertedID.(primitive.ObjectID).Hex(),
		Name:     input.Name,
		IsGoodBoi: input.IsGoodBoi,
	}, nil
}

func (db *DB) FindByID(id string) (*model.Dog, error) {
	collection := db.client.Database("animals").Collection("dogs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dog model.Dog
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("ID invalide: %w", err)
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&dog)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du chien: %w", err)
	}

	return &dog, nil
}

func (db *DB) ALL() ([]*model.Dog, error) {
	collection := db.client.Database("animals").Collection("dogs")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des chiens: %w", err)
	}
	defer cursor.Close(ctx)

	var dogs []*model.Dog
	for cursor.Next(ctx) {
		var dog model.Dog
		if err := cursor.Decode(&dog); err != nil {
			return nil, fmt.Errorf("erreur lors du décodage du chien: %w", err)
		}
		dogs = append(dogs, &dog)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("erreur du curseur: %w", err)
	}

	return dogs, nil
}
