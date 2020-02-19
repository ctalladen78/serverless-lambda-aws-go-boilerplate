package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
)

type Identity struct {
	CognitoIdentityId     string
	CognitoIdentityPoolId string
}

// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-authentication-flow.html
// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
// https://serverless-stack.com/chapters/mapping-cognito-identity-id-and-user-pool-id.html
// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
// https://stackoverflow.com/questions/29928401/how-to-get-the-cognito-identity-id-in-aws-lambda
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentityprovider/#CognitoIdentityProvider.SignUp
// https://docs.aws.amazon.com/sdk-for-go/api/service/cognitoidentity/
// https://aws.amazon.com/blogs/mobile/understanding-amazon-cognito-user-pool-oauth-2-0-grants/
// https://dev.to/piczmar_0/serverless-authorizers---custom-rest-authorizer-16
// https://theiconic.tech/authentication-an-in-depth-look-at-aws-cognito-658336515169
// https://aws.amazon.com/blogs/mobile/understanding-amazon-cognito-authentication-part-3-roles-and-policies/
// https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-lambda-authorizer-output.html
// https://stackoverflow.com/questions/45783057/aws-api-gateway-custom-authorizer-how-to-access-principalid-in-lambda
// https://www.npmjs.com/package/serverless-offline#token-authorizers
// https://aws.amazon.com/blogs/compute/simply-serverless-using-aws-lambda-to-expose-custom-cookies-with-api-gateway/
// https://medium.com/faun/securing-api-gateway-with-lambda-authorizers-62845032bc7d
// https://github.com/seedboxtech/aws-go-lib/blob/master/cognito/requestidentity.go
// use this to verify cognito jwt token from api gateway
// https://medium.com/@chandupriya93/providing-authorization-to-api-gateway-with-cognito-identity-pools-af451fd3b532
// auth flow: verify token, then verify expiry, then process claims
func ValidateUser(token string, id string, poolId string) (bool, error) {
	pid := &Identity{
		CognitoIdentityId:     id,
		CognitoIdentityPoolId: poolId,
	}

	// event.requestContext.identity.cognitoIdentityId
	validateUserIdentity(token, pid)
	return false, nil
}

// this is for users with admin permission
// arg1 string: cognito token, arg2 cognito Identity
// returns: bool
func validateUserIdentity(token string, identity *Identity) bool {
	// if identity.CognitoIdentityId == "" {
	// 	return true
	// }

	// temp_app backend credentials
	awsConfig := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(awsConfig)

	if err != nil {
		return false
	}

	// event.requestContext.identity.cognitoIdentityId
	cognitoClient := cognitoidentity.New(sess, awsConfig)

	input := cognitoidentity.LookupDeveloperIdentityInput{
		DeveloperUserIdentifier: aws.String(token), // aws app id token
		IdentityId:              aws.String(identity.CognitoIdentityId),
		IdentityPoolId:          aws.String(identity.CognitoIdentityPoolId),
		MaxResults:              aws.Int64(1),
	}

	results, err := cognitoClient.LookupDeveloperIdentity(&input)
	if err != nil {
		return false
	}

	log.Println("IS COGNITO TOKEN VALID", results.DeveloperUserIdentifierList)

	return len(results.DeveloperUserIdentifierList) >= 1
}
