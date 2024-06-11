package aws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// CreateAWSSession initializes an AWS session with LocalStack.
func CreateAWSSession() *session.Session {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(os.Getenv("AWS_REGION")),
		Endpoint: aws.String("http://localhost:4566"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),     // Access Key ID
			os.Getenv("AWS_SECRET_ACCESS_KEY"), // Secret Access Key
			""),                                // Token can be left blank for LocalStack
		S3ForcePathStyle: aws.Bool(true), // Use path-style addressing for S3
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	return sess
}

// CreateS3BucketIfNotExists CreateS3Bucket creates an S3 bucket.
func CreateS3BucketIfNotExists(svc *s3.S3, bucketName string) {
	_, err := svc.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err == nil {
		fmt.Printf("Bucket %s already exists\n", bucketName)
		return
	}

	// Create the bucket if it does not exist
	_, err = svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Fatalf("Failed to create S3 bucket: %v", err)
	}
	fmt.Printf("Bucket %s created successfully\n", bucketName)
}

// CreateIAMRole creates an IAM role for Lambda functions.
func CreateIAMRole(svc *iam.IAM) string {
	assumeRolePolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "lambda.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`

	result, err := svc.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String("lambda-ex"),
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
	})

	if err != nil {
		log.Fatalf("Failed to create role: %v", err)
	}

	return *result.Role.Arn
}

// CreateLambdaFunction creates a Lambda function.
func CreateLambdaFunction(svc *lambda.Lambda, roleArn string) {
	zipFile := []byte{} // Load your zip file content here

	_, err := svc.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String("MyLambdaFunction"),
		Runtime:      aws.String("go1.x"),
		Role:         aws.String(roleArn),
		Handler:      aws.String("main"),
		Code: &lambda.FunctionCode{
			ZipFile: zipFile,
		},
		Description: aws.String("A simple Lambda function"),
	})

	if err != nil {
		log.Fatalf("Failed to create Lambda function: %v", err)
	}
}

// CreateAPIGateway sets up an API Gateway that triggers the Lambda function.
func CreateAPIGateway(svc *apigateway.APIGateway, lambdaArn string) {
	restApiOutput, err := svc.CreateRestApi(&apigateway.CreateRestApiInput{
		Name: aws.String("MyAPI"),
	})
	if err != nil {
		log.Fatalf("Failed to create REST API: %v", err)
	}

	apiId := *restApiOutput.Id
	resourcesOutput, err := svc.GetResources(&apigateway.GetResourcesInput{
		RestApiId: aws.String(apiId),
	})
	if err != nil {
		log.Fatalf("Failed to get resources: %v", err)
	}

	rootId := *resourcesOutput.Items[0].Id
	resourceOutput, err := svc.CreateResource(&apigateway.CreateResourceInput{
		RestApiId: aws.String(apiId),
		ParentId:  aws.String(rootId),
		PathPart:  aws.String("mylambda"),
	})
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	resourceId := *resourceOutput.Id
	_, err = svc.PutMethod(&apigateway.PutMethodInput{
		RestApiId:         aws.String(apiId),
		ResourceId:        aws.String(resourceId),
		HttpMethod:        aws.String("GET"),
		AuthorizationType: aws.String("NONE"),
	})
	if err != nil {
		log.Fatalf("Failed to put method: %v", err)
	}

	uri := fmt.Sprintf(os.Getenv("API_GATEWAY_ARN"), lambdaArn)
	_, err = svc.PutIntegration(&apigateway.PutIntegrationInput{
		RestApiId:             aws.String(apiId),
		ResourceId:            aws.String(resourceId),
		HttpMethod:            aws.String("GET"),
		Type:                  aws.String("AWS_PROXY"),
		IntegrationHttpMethod: aws.String("POST"),
		Uri:                   aws.String(uri),
	})
	if err != nil {
		log.Fatalf("Failed to put integration: %v", err)
	}

	_, err = svc.CreateDeployment(&apigateway.CreateDeploymentInput{
		RestApiId: aws.String(apiId),
		StageName: aws.String("test"),
	})
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}

	fmt.Printf("API Gateway created successfully\n")
}

// CreateSQSQueue creates an SQS queue.
func CreateSQSQueue(svc *sqs.SQS) {
	_, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("MyQueue"),
	})
	if err != nil {
		log.Fatalf("Failed to create SQS queue: %v", err)
	}
}

// CreateDynamoDBTable creates a DynamoDB table.
func CreateDynamoDBTable(svc *dynamodb.DynamoDB) {
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String("MyTable"),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		},
	})

	if err != nil {
		log.Fatalf("Failed to create DynamoDB table: %v", err)
	}
}
