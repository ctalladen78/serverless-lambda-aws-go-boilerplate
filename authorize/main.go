package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"golang-projects/aws_auth_cognito/utility"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// use postman to sign in, get access jwt token
// use token to sign request, pass in policy document with principalId
type Response events.APIGatewayProxyResponse

// https://github.com/serverless/examples/blob/master/aws-golang-auth-examples/functions/auth/main.go
// https://www.npmjs.com/package/serverless-offline#token-authorizers
// https://gist.github.com/ctalladen78/753395c5bc49de019c55f6495e901398
// https://github.com/sohamkamani/jwt-go-example

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler0(ctx context.Context, req events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer
	auth := req.RequestContext.Authorizer
	p := req.RequestContext.Identity.CognitoAuthenticationProvider
	i := req.RequestContext.Identity.CognitoIdentityID
	fmt.Println("authorizer", auth)
	fmt.Println("provider", p)
	fmt.Println("cognito user", i)

	// https://github.com/GavL89/lambdaauthorizer-dotnet/blob/master/Authorizer/CheckToken.cs
	type PolicyDocument struct {
		PrincipalID    string            `json:"principalId"`
		PolicyDocument map[string]string `json:"policyDocument"`
	}
	email := "test@test.com"
	pdoc := make(map[string]string)
	pdoc["Version"] = ""
	pdoc["Statement"] = ""

	doc := &PolicyDocument{
		PrincipalID:    email,
		PolicyDocument: pdoc,
	}
	j, err := json.Marshal(doc)

	// validate event.requestContext.identity.cognitoIdentityId
	// isAccessTokenValid, err := ValidateUser(
	// 	req.RequestContext.Identity.AccessKey,
	// 	req.RequestContext.Identity.CognitoIdentityID,
	// 	req.RequestContext.Identity.CognitoIdentityPoolID,
	// )

	// if !isAccessTokenValid || err != nil {
	// 	return Response{StatusCode: 404}, err
	// }

	// body, err := json.Marshal(map[string]interface{}{
	// 	"message": "Cognito Authorizer successfully validated token",
	// })
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, j)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "world-handler",
		},
	}

	return resp, nil
}

func handler(request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	token := request.AuthorizationToken
	fmt.Println("authorizer", token)
	tokenSlice := strings.Split(token, " ")
	var bearerToken string
	if len(tokenSlice) > 1 {
		bearerToken = tokenSlice[len(tokenSlice)-1]
	}
	isValid, err := utility.ValidateToken(bearerToken)
	if err != nil {
		// redirect response
		// return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
		return generatePolicy("user", "Deny", ""), errors.New("Unauthorized")
	}
	// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html
	// https://docs.aws.amazon.com/apigateway/latest/developerguide/apigateway-use-lambda-authorizer.html
	if !isValid {
		return generatePolicy("user", "Deny", ""), errors.New("Unauthorized")
	}

	return generatePolicy("user", "Allow", request.MethodArn), nil
}

func main() {
	lambda.Start(handler)
}

func generatePolicy(principalID, effect, resource string) events.APIGatewayCustomAuthorizerResponse {
	authResponse := events.APIGatewayCustomAuthorizerResponse{PrincipalID: principalID}

	if effect != "" && resource != "" {
		authResponse.PolicyDocument = events.APIGatewayCustomAuthorizerPolicy{
			Version: "2012-10-17",
			Statement: []events.IAMPolicyStatement{
				{
					Action:   []string{"execute-api:Invoke"},
					Effect:   effect,
					Resource: []string{resource},
				},
			},
		}
	}
	return authResponse
}
