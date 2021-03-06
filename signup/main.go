package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"golang-projects/aws_auth_cognito/utility"
)

type Response events.APIGatewayProxyResponse

func main() {
	lambda.Start(Handler)
}

func CreateNewUser(c *utility.Credentials) (string, error) {
	fmt.Println("NEW USER CREDS", c)

	// check if user exists
	userId, err := utility.ValidateCredentials(c.Email, c.Password)
	if userId != "" {
		fmt.Println("USER EXISTS ", userId)
		return "", errors.New("user exists ")
	}
	// spawn new goroutine to save credentials
	_, err = utility.SaveCredentials(c.Email, c.Password)
	// if err != nil {
	// 	return "", err
	// }
	newtoken, err := utility.Signup(c.Email, c.Password)
	if err != nil {
		return "", err
	}
	// newtoken = "testtoken"

	return newtoken, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	// use SSL
	// https://stackoverflow.com/questions/14409930/how-to-safely-include-password-in-query-string
	fmt.Println("REQUEST BODY", request.Body)
	c := &utility.Credentials{
		Username: request.Headers["x-username"],
		Password: request.Headers["x-password"],
		Email:    request.Headers["x-email"],
	}
	// use velocity template to parse body
	// err := json.Unmarshal([]byte(request.Body), c)
	// decoder := json.NewDecoder(request.Body)
	// err = decoder.Decode(c)
	// if err != nil {
	// 	return Response{
	// 		StatusCode: 404,
	// 	}, err
	// }
	for k, v := range request.Headers {
		fmt.Println("REQUEST QUERY PARAMS", k, v)
	}

	// https://aws.amazon.com/blogs/compute/simply-serverless-using-aws-lambda-to-expose-custom-cookies-with-api-gateway/
	token, err := CreateNewUser(c)
	if err != nil {
		return Response{
			StatusCode: 404,
			Body:       "error creating new user",
		}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": "successfully signed up new user",
	})
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      201,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "notes-handler",
			"X-Access-Token":         token,
		},
	}

	return resp, nil
}
