package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

// https://github.com/jpcedenog/blog-api-security-authentication/blob/master/createnote/main.go
// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse


// TODO handler for /todo_new

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	id, err := uuid.NewUUID()
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	// https://serverless-stack.com/chapters/invoke-api-gateway-endpoints-locally.html
	cognitoIdentityID := request.RequestContext.Identity.CognitoIdentityID
	myNote := Todo{
		UserID:  cognitoIdentityID,
		TodoID:  id.String(),
		Content: request.Body,
	}

	err = putTodo(myNote)
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "Note created successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      201,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "notes-handler",
		},
	}

	return resp, nil
}

func putTodo(todo Todo) error {
	dattr, err := dynamodbattribute.MarshalMap(todo)
	if err != nil {
		return err
	}

	svc := dynamodb.New(session.Must(session.NewSession()))
	input := &dynamodb.PutItemInput{
		Item:      dattr,
		TableName: aws.String(os.Getenv("tableName")),
	}
	_, err = svc.PutItem(input)
	return nil

}

func main() {
	lambda.Start(Handler)
}
