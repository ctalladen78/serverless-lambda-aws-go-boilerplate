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
	fmt.Println("QUERY RESULT", userList)
	for _, v := range userList {
		fmt.Println("USER ", v)
		if password == v.Password {
			return v.ObjectId, err
		}
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

func GetTodoListByUser(userid string) ([]*TodoObject, error) {

	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	todoList, err := dbCtrl.QueryTodo("todotable1", CREATED_BY, userid)
	if err != nil {
		// return Response{StatusCode: 404}, err
		return nil, err
	}

	// test object
	todoList = append(todoList, &TodoObject{Todo: "testtest"})
	fmt.Println("TODO LIST ", todoList)
	return todoList, nil

}

func PutTodo(table string, todo *TodoObject) (bool, error) {
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	todo.ObjectId = "TODO-" + shortuuid.New()
	todo.CreatedAt = time.Time.Format(time.Now(), time.RFC3339)
	_, err := dbCtrl.PutItem(table, todo)
	if err != nil {
		return false, err
	}
	return true, nil
}
