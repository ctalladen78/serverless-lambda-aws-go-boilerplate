package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// https://github.com/jpcedenog/blog-api-security-authentication/blob/master/createnote/main.go
// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
type Response events.APIGatewayProxyResponse

type Todo struct {
	UserID  string `json:"userId"`
	TodoID  string `json:"noteId"`
	Content string `json:"content"`
}

// TODO handler for /todo_list

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var buf bytes.Buffer

	todoList, err := getTodoList()
	if err != nil {
		return Response{StatusCode: 404}, err
	}

	body, err := json.Marshal(map[string]interface{}{
		"message": todoList,
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

func getTodoList() ([]Todo, error) {
	out := []Todo{}
	// l := []map[string]*dynamodb.AttributeValue{}

	// svc := dynamodb.New(session.Must(session.NewSession()))
	// TODO dynamodb query or scan
	// input := &dynamodb.QueryInput{

	// }
	// todoList := svc.Query(input)
	// input := &dynamodb.ScanInput{
	// 	TableName: aws.String("Todo"),
	// }
	// todoList, err := svc.Scan(input)
	// l = todoList.Items
	t1 := &Todo{
		UserID:  "",
		TodoID:  "",
		Content: "",
	}
	t2 := &Todo{
		UserID:  "",
		TodoID:  "",
		Content: "",
	}
	out = append(out, *t1)
	out = append(out, *t2)

	// err = dynamodbattribute.UnmarshalListOfMaps(l, out)
	// if err != nil {
	// 	// return Response{StatusCode: 404}, err
	// 	return nil, err
	// }
	return out, nil

}

func main() {
	lambda.Start(Handler)
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	fmt.Println(err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, err
}
