package main

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

func FirebaseInit(ctx context.Context) (*messaging.Client, error) {
	// Use the path to your service account credential json file
	opt := option.WithCredentialsFile("Path to Service Account.json File")

	// Create a new firebase app
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	// Get the FCM object
	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}

	return fcmClient, nil
}
