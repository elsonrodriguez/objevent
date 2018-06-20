package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
        "github.com/aws/aws-sdk-go/service/s3"

	"cloud.google.com/go/storage"
	"cloud.google.com/go/pubsub"
	"golang.org/x/net/context"
)


func awsHandler(bucketName string, endpointURL string) bool {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess)
	s3svc := s3.New(sess)


	topicParams := &sns.CreateTopicInput{
		Name: aws.String(bucketName),
	}

        topicResponse, err := svc.CreateTopic(topicParams)

	if err != nil { 
		fmt.Println(err.Error())
	}

	fmt.Println(topicResponse)

	endpointURLParsed, err := url.Parse(endpointURL)

	if err != nil {
		fmt.Println(err)
	}


        subParams := &sns.SubscribeInput{
		Endpoint: aws.String(endpointURL),
		Protocol: aws.String(endpointURLParsed.Scheme),
		TopicArn: topicResponse.TopicArn,
	}

	subscriptionResponse, err := svc.Subscribe(subParams)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(subscriptionResponse)

	snsPermParams := &sns.AddPermissionInput{
		AWSAccountId: aws.StringSlice([]string{"*",}),
		ActionName: aws.StringSlice([]string{"Subscribe","Receive","Publish"}),
		Label: aws.String("objevent default policy"),
		TopicArn: topicResponse.TopicArn,
	}

	snsPermResponse, err := svc.AddPermission(snsPermParams)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Permissions be like:")
	fmt.Println(snsPermResponse)

	//need to confirm subscription with token... might need some app logic, the app needs to have a token endpoint, maybe use the ARN as a unique path. endpointURL/tokens/{arn}/


	topicConfigParams := []*s3.TopicConfiguration{
		&s3.TopicConfiguration {
			Events: []*string{aws.String("s3:ObjectCreated:*")},
			TopicArn: topicResponse.TopicArn,
		},
	}

	notificationParams := &s3.NotificationConfiguration{
		TopicConfigurations: topicConfigParams,
	}

	bucketnotificationParams := &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketName),
		NotificationConfiguration: notificationParams,
	}

	notificationResponse, err := s3svc.PutBucketNotificationConfiguration(bucketnotificationParams)

	if err != nil {
 		fmt.Println(err.Error())
        }     

	fmt.Println(notificationResponse)

	fmt.Printf("Added event notifications for bucket %s to endpoint %s", bucketName, endpointURL)
	return true
}

func gcpHandler(bucketName string, endpointURL string) bool {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		out, err := exec.Command("gcloud", "config", "get-value", "project").Output()
		if err != nil {
			fmt.Println(err.Error())
		}
		projectID = strings.TrimSpace(string(out))
	}
	
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
	return true
}

func main () {
        bucket := os.Getenv("OBJEVENT_BUCKET")
	endpointURL := os.Getenv("OBJEVENT_ENDPOINT_URL")

	bucketURL, err := url.Parse(bucket)

	if err != nil {
		fmt.Println(err)
	}

	bucketName := bucketURL.Host
	bucketScheme := bucketURL.Scheme

	result := false

	if bucketScheme == "s3" {
		result = awsHandler(bucketName, endpointURL)
	} else if bucketScheme == "gs"{
		result = gcpHandler(bucketName, endpointURL)
	} else { //TODO figure out minio
		fmt.Println("Unsupported bucket protocol!")
	}

	fmt.Println(result)
}

