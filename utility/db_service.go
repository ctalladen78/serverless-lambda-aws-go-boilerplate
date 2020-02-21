package utility

import (
	"errors"
	"fmt"
)

func ValidateCredentials(email string, password string) (bool, error) {
	dbCtrl := InitLocalDbConnection("http://localhost:8000")
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
	dbCtrl := InitLocalDbConnection("http://localhost:8000")
	u := &UserObject{
		Email:    email,
		Password: password,
	}
	dbCtrl.PutItem("usertable", u)
	return false, nil
}
