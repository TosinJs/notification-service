package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	mc, err := Init(ctx, "MONGO_URI")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mongoDB := mc.Database("notification_service")
	tokenCollection := mongoDB.Collection("notificationTokens")

	const hours_in_a_week = 24 * 7
	//create the index model with the field "timestamp"
	index := mongo.IndexModel{
		Keys: bson.M{"timestamp": 1},
		Options: options.Index().SetExpireAfterSeconds(
			int32((time.Hour * 3 * hours_in_a_week).Seconds()),
		),
	}
	//Create the index on the token collection
	_, err = tokenCollection.Indexes().CreateOne(ctx, index)
	if err != nil {
		fmt.Printf("mongo index error: %v", err)
		os.Exit(1)
	}

	fcmClient, err := FirebaseInit(ctx)
	if err != nil {
		fmt.Printf("error connecting to firebase: %v", err)
		os.Exit(1)
	}

	// Route to POST
	r := chi.NewRouter()
	r.Post("/tokens", func(w http.ResponseWriter, r *http.Request) {
		var token NotificationToken
		err := json.NewDecoder(r.Body).Decode(&token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = InsertToken(tokenCollection, token, ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	//Trigger a notification
	r.Post("/send-notifications", func(w http.ResponseWriter, r *http.Request) {
		var message Message
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		tokens, err := GetTokens(tokenCollection, ctx, message.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = SendNotification(fcmClient, ctx, tokens, "UserId", "Message")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Server Starting on Port 3000")
	http.ListenAndServe(":3000", r)
}
