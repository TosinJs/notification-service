package main

import (
	"context"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/bson/primitive"
)

// Get all the tokens registered for a user
func GetTokens(
  coll *mongo.Collection,
  ctx context.Context,
  userId string,
) ([]string, error) {
  filter := bson.D{{Key: "userId", Value: userId}}
	tokenCursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0)
	for tokenCursor.Next(ctx) {
		var token NotificationToken
		err = tokenCursor.Decode(&token)
		tokens = append(tokens, token.DeviceId)
	}

	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// Insert a token
func InsertToken(
  coll *mongo.Collection, 
  token NotificationToken,
  ctx context.Context,
) error {
  // Check if the token already exists
  filter := bson.D{{Key: "deviceId", Value: token.DeviceId}}
	res := coll.FindOne(ctx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
      // If token does not exist insert it
			token.ID = primitive.NewObjectID()
			_, err := coll.InsertOne(ctx, token)
			return err
		}
		return res.Err()
	}

  // If token exists update the timestamp to now
	_, err := coll.UpdateOne(ctx, filter, bson.M{"$set": bson.M{"timestamp": time.Now().UTC()}})
	return err
}