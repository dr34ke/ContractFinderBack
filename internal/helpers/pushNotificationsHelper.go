package helper

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func SendPushNotification(token string, title string, body string) (err error) {
	// Initialize Firebase app
	opt := option.WithCredentialsFile("../../firebase-token.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return err
	}

	// Initialize Messaging client
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
		return err
	}

	// Define the message
	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Send the message
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalf("error sending message: %v\n", err)
		return err
	}
	log.Printf("Successfully sent message: %s\n", response)
	return nil
}
