package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang-projects/aws_auth_cognito/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// https://github.com/jpcedenog/blog-api-security-authentication/blob/master/createnote/main.go
// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse

// TODO handler for /todo_new

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	// id, err := uuid.NewUUID()
	// if err != nil {
	// 	return Response{StatusCode: 404}, err
	// }

	// https://serverless-stack.com/chapters/invoke-api-gateway-endpoints-locally.html
	cognitoIdentityID := request.RequestContext.Identity.CognitoIdentityID
	// cognitoSessionToken := request.RequestContext.Authorizer

	fmt.Println("COGNITO IDENTITY ID", cognitoIdentityID)
	type PostBody struct {
		Data string `json:"data"`
	}
	data := &PostBody{}
	// TODO use velocity template to parse body
	err := json.Unmarshal([]byte(request.Body), data)
	fmt.Println("REQUEST BODY RAW", request.Body)
	fmt.Println("REQUEST DATA", data)
	// decoder := json.NewDecoder(request.Body)
	// err = decoder.Decode(c)
	// if err != nil {
	// 	return Response{
	// 		StatusCode: 404,
	// 	}, err
	// }

	myNote := &utility.TodoObject{
		CreatedBy: cognitoIdentityID,
		Todo:      data.Data,
	}
	_, err = utility.PutTodo(myNote)
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

func main() {
	lambda.Start(Handler)
}
