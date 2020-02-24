package utility

import (
	"errors"
	"fmt"
	"time"

	"github.com/lithammer/shortuuid"
)

// return userId if exists
func ValidateCredentials(email string, password string) (string, error) {
	// dbCtrl := InitLocalDbConnection("http://localhost:8000")
	// see google: github issue setting ifconfig alias dynamodb local
	// ifconfig lo0 alias 172.16.123.1
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	// maybe returns list
	userList, err := dbCtrl.QueryUser(EMAIL, email)
	if err != nil {
		return "", err
	}
	u := &UserObject{}
	fmt.Println("QUERY RESULT", userList)
	for _, v := range userList {
		if password == v.Password {
			return "", err
		}
		u = v
		fmt.Println("USER FOUND", u)
		return v.ObjectId, nil
	}
	// not found
	return "", errors.New("USER NOT FOUND")
}

func SaveCredentials(email string, password string) (bool, error) {
	// open dynamodb connection
	// dbCtrl := InitLocalDbConnection("http://localhost:8000")
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	u := &UserObject{
		Email:    email,
		Password: password,
		ObjectId: shortuuid.New(),
	}
	_, err := dbCtrl.PutItem("usertable", u)
	if err != nil {
		return false, err
	}
	return true, nil
}

func GetTodoList() ([]Todo, error) {
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

func PutTodo(todo *TodoObject) (bool, error) {
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	todo.ObjectId = shortuuid.New()
	todo.CreatedAt = time.Now()
	_, err := dbCtrl.PutItem("todotable", todo)
	if err != nil {
		return false, err
	}
	return true, nil
}
