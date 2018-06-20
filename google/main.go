package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"cloud.google.com/go/storage"
	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)


func main() {
	ctx := context.Background()

	//TODO: check for GOOGLE_CLOUD_PROJECT before trying gcloud config
	out, err := exec.Command("gcloud", "config", "get-value", "project").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	projectID := strings.TrimSpace(string(out))

        bucketName := os.Getenv("OBJEVENT_BUCKET_NAME")
	endpointURL := os.Getenv("OBJEVENT_ENDPOINT_URL")

	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		fmt.Println(err.Error())
	}
	
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
                fmt.Println(err.Error())
	}
	
	topicResponse, err := client.CreateTopic(ctx, bucketName)

	if err != nil {
		fmt.Println(err.Error())
	}
	
	fmt.Println(topicResponse)

	//The endpoint, if not hosted on Appengine, needs to be verified on google's webmaster tools AND added to your project. I don't think I can automate the webmaster tools.
	pushParams := pubsub.PushConfig{
		Endpoint: endpointURL,
	}

	subParams := pubsub.SubscriptionConfig{
		Topic: client.Topic(bucketName),
		PushConfig: pushParams,
	}

	subResponse, err := client.CreateSubscription(ctx, bucketName, subParams)

	if err != nil {
		fmt.Println(err.Error())
	}
	
	fmt.Println(subResponse)
	
	bucket := storageClient.Bucket(bucketName)

	notificationParams := &storage.Notification {
		TopicProjectID: projectID,
		TopicID: bucketName,
		PayloadFormat: storage.JSONPayload,
	}

	//need to sleep here until the topic and subscription are available.

	notificationResponse, err := bucket.AddNotification(ctx, notificationParams)

	if err != nil {
		fmt.Println(err.Error())

	}

	fmt.Println(notificationResponse)
	fmt.Printf("Added event notifications for bucket %s to endpoint %s", bucketName, endpointURL)
}

