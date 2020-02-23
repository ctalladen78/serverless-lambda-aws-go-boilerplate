package utility

import (
	"errors"
	"fmt"
)

func ValidateCredentials(email string, password string) (bool, error) {
	// dbCtrl := InitLocalDbConnection("http://localhost:8000")
	// see google: github issue setting ifconfig alias dynamodb local
	// ifconfig lo0 alias 172.16.123.1
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	// maybe returns list
	userList, err := dbCtrl.QueryUser(EMAIL, email)
	if err != nil {
		return false, err
	}
	u := &UserObject{}
	fmt.Println("QUERY RESULT", userList)
	for _, v := range userList {
		if password == v.Password {
			return false, err
		}
		u = v
		fmt.Println("USER FOUND", u)
		return true, nil
	}
	// not found
	return false, errors.New("USER NOT FOUND")
}

func SaveCredentials(email string, password string) (bool, error) {
	// open dynamodb connection
	// dbCtrl := InitLocalDbConnection("http://localhost:8000")
	dbCtrl := InitLocalDbConnection("http://172.16.123.1:8000")
	fmt.Println("DB PATH %s", dbCtrl)
	u := &UserObject{
		Email:    email,
		Password: password,
	}
	_, err := dbCtrl.PutItem("usertable", u)
	if err != nil {
		return false, err
	}
	return true, nil
}
