# ObjEvent

Small utility for adding events to multiple object stores.

For use in conjunction with Event Gateway.

## Usage

First, your environment must have valid credentials set for Google Cloud Platform or Amazon Web Services depending on where your bucket is.

### AWS Setup

Set the following variables, or configure ~/.aws/credentials:

```
export AWS_ACCESS_KEY_ID=xxxxx
export AWS_SECRET_ACCESS_KEY=xxxxx
export AWS_REGION=us-west-2 #this must be set even if your config file has a region
```

### GCP Setup

Define a service account and a project

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/serviceaccount.key
export GOOGLE_CLOUD_PROJECT=project_id #if this is not set we will scrape your current project id from the gcloud client
```
### Invocation 

Next, run the tool

```
OBJEVENT_BUCKET=s3://bucketname
OBJEVENT_ENDPOINT_URL=https://youreventhandler.com/
./objevent
```

## Eratta

GCP and AWS have mechanisms for authenticating an endpoint URL.

For GCP, you must first add a domain to their [Webmaster Tools](https://www.google.com/webmasters/verification/home?hl=en), and then add the Domain to your project under the [Google API Console](https://console.developers.google.com/apis/credentials/domainverification?project=objectevent&folder&organizationId). Due the the sheer number of registrars, this tool will not be able to automate the verification, but we might be able to automate the addition to the Project.

For AWS, there is a token sent to your endpoint which must be [Confirmed](https://docs.aws.amazon.com/sdk-for-go/api/service/sns/#SNS.ConfirmSubscription). This can be automated if a convention is arrived at for retreiving tokens from endpoints.
