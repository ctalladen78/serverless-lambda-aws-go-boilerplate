package utility

import (
	// "fmt"
	// "io"
	// "net/http"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

// github.com/nordcloud/cognito-go-auth
// github.com/golang-cognito-example
// github.com/serinth/serverless-cognito-auth
type CognitoService struct {
	CognitoClient *cognito.CognitoIdentityProvider
	UsernameFlow  *UsernameFlow
	AuthFlow      string
	UserPoolID    string
	AppClientID   string
}

type UsernameFlow struct {
	Username string
}

func CognitoServiceSetup() (*CognitoService, error) {
	conf := &aws.Config{Region: aws.String("us-east-1")}
	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}
	c := &CognitoService{
		CognitoClient: cognito.New(sess),
		UsernameFlow:  &UsernameFlow{},
		AuthFlow:      "USER_PASSWORD_AUTH",
	}
	return c, nil
}

func (c *CognitoService) GetUserDetails(u string) {

	_, err := c.CognitoClient.AdminGetUser(&cognito.AdminGetUserInput{
		UserPoolId: aws.String(c.UserPoolID),
		Username:   aws.String(u),
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok {
			if awsErr.Code() == cognito.ErrCodeUserNotFoundException {
				// safe to use username ok
				fmt.Sprintf("USERNAME NOT FOUND ", u)
				return
			}
		} else {
			// registration error
			fmt.Sprintf("NAME EXISTS, REGISTER WITH ANOTHER")
		}
	}

	fmt.Sprintf("NAME EXISTS, REGISTER WITH ANOTHER")
	return
}

func (c *CognitoService) SignUp(email string, password string) (*cognito.SignUpOutput, error) {

	signUpInput := &cognito.SignUpInput{
		Username: aws.String(""),
		Password: aws.String(""),
		ClientId: aws.String(""),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("phone_number"),
				Value: aws.String(""),
			},
		},
	}

	return c.CognitoClient.SignUp(signUpInput)

}

func (c *CognitoService) ConfirmSignUp(otp string) (*cognito.ConfirmSignUpOutput, error) {
	input := &cognito.ConfirmSignUpInput{
		ConfirmationCode: aws.String(otp),
		Username:         aws.String(c.UsernameFlow.Username),
		ClientId:         aws.String(c.AppClientID),
	}

	return c.CognitoClient.ConfirmSignUp(input)
}

func (c *CognitoService) SignIn(email string, password string) (*cognito.AuthenticationResultType, error) {
	// get credentials
	// get refresh token
	params := map[string]*string{
		"USERNAME": aws.String(email),
		"PASSWORD": aws.String(password),
	}

	authInput := &cognito.InitiateAuthInput{
		AuthFlow:       aws.String(c.AuthFlow),
		AuthParameters: params,
		ClientId:       aws.String(c.AppClientID),
	}
	res, err := c.CognitoClient.InitiateAuth(authInput)
	if err != nil {
		return nil, err
	}

	return res.AuthenticationResult, nil
}
