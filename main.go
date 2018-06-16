package main

import (
	"fmt"
	"os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
        "github.com/aws/aws-sdk-go/service/s3"
)


func main() {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess)
	s3svc := s3.New(sess)

        bucketName := os.Getenv("OBJEVENT_BUCKET_NAME")
	endpointURL := os.Getenv("OBJEVENT_ENDPOINT_URL")

	topicParams := &sns.CreateTopicInput{
		Name: aws.String(bucketName),
	}

        topicResponse, err := svc.CreateTopic(topicParams)

	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	fmt.Println(topicResponse)

        subParams := &sns.SubscribeInput{
		Endpoint: aws.String(endpointURL),
		Protocol: aws.String("http"), // Parse this out later		
		TopicArn: topicResponse.TopicArn,
	}

	subscriptionResponse, err := svc.Subscribe(subParams)

	if err != nil {
		fmt.Println(err.Error())
		return
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
		return
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
		return
        }     

	fmt.Println(notificationResponse)

	fmt.Printf("Added event notifications for bucket %s to endpoint %s", bucketName, endpointURL)
}
