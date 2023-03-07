package main

import (
	"context"

	"firebase.google.com/go/v4/messaging"
)

func SendNotification(
	fcmClient *messaging.Client,
	ctx context.Context,
	tokens []string,
	userId, message string,
) error {
	//Send to One Token
	_, err := fcmClient.Send(ctx, &messaging.Message{
		Token: tokens[0],
		Data: map[string]string{
			message: message,
		},
	})
	if err != nil {
		return err
	}

	//Send to Multiple Tokens
	_, err = fcmClient.SendMulticast(ctx, &messaging.MulticastMessage{
		Data: map[string]string{
			message: message,
		},
		Tokens: tokens,
	})
	return err
}
