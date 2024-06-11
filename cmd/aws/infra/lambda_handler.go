// cmd/aws/infra/lambda_handler.go
package infra

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is the main function for the Lambda
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var body map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Invalid request body"}, nil
	}

	responseBody, err := json.Marshal(map[string]string{"message": "Hello, World!"})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Failed to marshal response"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 200, Body: string(responseBody)}, nil
}

func main() {
	lambda.Start(Handler)
}
