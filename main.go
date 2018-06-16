package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
        "github.com/aws/aws-sdk-go/service/s3"
)


func main() {
	sess := session.Must(session.NewSession())
	svc := sns.New(sess)
	s3svc := s3.New(sess)

        bucketname := "superpuke"
	endpointurl := "http://google.com"

	topicparams := &sns.CreateTopicInput{
		Name: aws.String(bucketname),
	}

        topicresp, err := svc.CreateTopic(topicparams)

	if err != nil { 
		fmt.Println(err.Error())
		return
	}

	fmt.Println(topicresp)

        subparams := &sns.SubscribeInput{
		Endpoint: aws.String(endpointurl),
		Protocol: aws.String("http"), // Parse this out later		
		TopicArn: topicresp.TopicArn,
	}

	subresp, err := svc.Subscribe(subparams)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(subresp)

	//now need to confirm subscription with token... might need some app logic, the app needs to have a token endpoint, maybe use the ARN as a unique path. /tokens/{arn}/

        eventlist := "s3:ObjectCreated:*"

	topicconfigparams := []*s3.TopicConfiguration{
		&s3.TopicConfiguration {
			Events: []*string{aws.String(eventlist)},
			TopicArn: topicresp.TopicArn,
		},
	}

	notificationparams := &s3.NotificationConfiguration{
		TopicConfigurations: topicconfigparams,
	}

	bucketnotificationparams := &s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(bucketname),
		NotificationConfiguration: notificationparams,
	}

	notificationresp, err := s3svc.PutBucketNotificationConfiguration(bucketnotificationparams)

	if err != nil {
 		fmt.Println(err.Error())
		return
        }     

	fmt.Println(notificationresp)

}
