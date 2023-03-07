package main 

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
//Function to connect to mongo database instance
func Init(ctx context.Context, URI string) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}
    //Ping the database to check connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	fmt.Println("Successfully Connected to The Database")
	return client, nil
}