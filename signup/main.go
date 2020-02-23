package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"golang-projects/aws_auth_cognito/utility"
)

type Response events.APIGatewayProxyResponse

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// TODO return auth access token
func main() {
	lambda.Start(Handler)
}

func checkUserExists(email string, pwd string) (bool, error) {
	return utility.ValidateCredentials(email, pwd)
}

func CreateNewUser(c *Credentials) (string, error) {
	fmt.Println("NEW USER CREDS", c)

	// check if user exists
	// didUserExist, err := checkUserExists(c.Email, c.Password)
	// if didUserExist {
	// 	return "", errors.New("user exists")
	// }
	// TODO spawn new goroutine to save credentials
	_, err := utility.SaveCredentials(c.Email, c.Password)
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

	// TODO use SSL
	// https://stackoverflow.com/questions/14409930/how-to-safely-include-password-in-query-string
	fmt.Println("REQUEST BODY", request.Body)
	c := &Credentials{
		Username: request.Headers["x-username"],
		Password: request.Headers["x-password"],
		Email:    request.Headers["x-email"],
	}
	// TODO use velocity template to parse body
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
