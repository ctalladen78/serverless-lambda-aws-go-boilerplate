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

type Response events.APIGatewayProxyResponse

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	token := request.Headers["x-access-token"]
	fmt.Println("AUTH ", token)

	// tokenSlice := strings.Split(token, " ")
	// var bearerToken string
	// if len(tokenSlice) > 1 {
	// 	bearerToken = tokenSlice[len(tokenSlice)-1]
	// }

	c := &Credentials{
		Username: request.Headers["x-username"],
		Password: request.Headers["x-password"],
		Email:    request.Headers["x-email"],
	}

	// https://aws.amazon.com/blogs/compute/simply-serverless-using-aws-lambda-to-expose-custom-cookies-with-api-gateway/
	newtoken, err := utility.Signin(token, c.Email, c.Password)
	if err != nil {
		e, _ := json.Marshal(map[string]interface{}{"error": "user does not exists"})
		return Response{
			StatusCode: 404,
			Body:       string(e),
		}, err
		if newtoken == "" {
			e, _ := json.Marshal(map[string]interface{}{"error": "invalid token"})
			return Response{
				StatusCode: 404,
				Body:       string(e),
			}, err
		}
	}
	fmt.Println("NEW TOKEN ", newtoken)

	body, err := json.Marshal(map[string]interface{}{
		"message": "successfully signed in",
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
			"X-Access-Token":         newtoken,
		},
	}

	return resp, nil
}
