package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/redbonzai/user-management-api/internal/config"
	"github.com/redbonzai/user-management-api/internal/db"
	"github.com/redbonzai/user-management-api/internal/infrastructure"
	"github.com/redbonzai/user-management-api/pkg/logger"
	"go.uber.org/zap"
)

// @title User Management API
// @version 1.0
// @description This is a sample User Management server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /users
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize logger
	logger.InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("could not load config:", zap.Error(err))
	}

	logger.Info("-- Config loaded successfully --")
	db.InitDB(cfg)

	router := infrastructure.NewRouter()

	// Serve static files for Swagger
	router.Static("/swagger", "docs/swagger.yaml")

	// sess := aws.CreateAWSSession()

	//s3svc := s3.New(sess)
	//iamClient := iam.New(sess)
	//lambdaClient := lambda.New(sess)
	//apiGateWayClient := apigateway.New(sess)
	//sqsClient := sqs.New(sess)
	//dynamodbClient := dynamodb.New(sess)

	//aws.CreateS3BucketIfNotExists(s3svc, "user-management-api-bucket")
	//roleArn := aws.CreateIAMRole(iamClient)
	//aws.CreateLambdaFunction(lambdaClient, roleArn)
	//
	//lambdaArn := fmt.Sprintf(os.Getenv("LAMBDA_ARN"))
	//aws.CreateAPIGateway(apiGateWayClient, lambdaArn)
	//aws.CreateSQSQueue(sqsClient)
	//aws.CreateDynamoDBTable(dynamodbClient)

	if err := router.Start(cfg.ServerAddress); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
