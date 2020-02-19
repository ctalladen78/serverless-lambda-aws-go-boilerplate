package utility

// github.com/sohamkamani/jwt-go-example
import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lithammer/shortuuid"
	// dynamodbclient
)

var private_key = []byte("PRIVATE_KEY")

type jwtCustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"email"`
}

// parse token for validity
// if token expired then issue new token with claims
// check claims.credentials with db
// if claims invalid then redirect to signup page
// on success issue new token
func Signin(token string, email string, password string) (string, error) {

	// isValid, err := ValidateToken(token)
	// if err != nil {
	// 	if !isValid {
	// 		return "", errors.New("INVALID TOKEN")
	// 	}
	// 	return "", errors.New("INVALID TOKEN")
	// }
	_, err := ValidateCredentials(email, password)
	if err != nil {
		return "", errors.New("INVALID CREDENTIALS")
	}

	// Set custom claims
	// claims := &jwtCustomClaims{
	// 	"Jon Snow",
	// 	true,
	// 	jwt.StandardClaims{
	// 		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
	// 	},
	// }

	// new token expiration for 5 mins
	expiration := time.Now().Add(5 * time.Minute)

	// Create token
	unsigned := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := unsigned.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["admin"] = true
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["exp"] = expiration
	signed, err := unsigned.SignedString(private_key)
	if err != nil {
		return "", err
	}
	fmt.Println("NEW TOKEN ", signed)
	return signed, nil
}

// save credentials to db
// issue new token with claims
func Signup(email string, password string) (string, error) {

	// new token expiration for 5 mins
	expiration := time.Now().Add(5 * time.Minute)

	// Create token
	unsigned := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := unsigned.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["objectid"] = shortuuid.New()
	claims["admin"] = true
	// claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["exp"] = expiration
	signed, err := unsigned.SignedString(private_key)

	// https://stackoverflow.com/questions/28204385/using-jwt-go-library-key-is-invalid-or-invalid-type
	// unsigned := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	// signed, err := unsigned.SignedString(private_key)
	if err != nil {
		return "", err
	}
	fmt.Println("NEW TOKEN ", signed)
	return signed, nil
}

// used for signin
// parse token for valid signature
// get claim info
// issue new token with claim
func Refresh(token string) (string, error) {
	claims := &jwtCustomClaims{}
	// https://godoc.org/github.com/dgrijalva/jwt-go
	// jwt.ParseRSAPublicKeyFromPEM()
	key, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return private_key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("SIGNATURE INVALID")
		}
		return "", err
	}
	if !key.Valid {
		fmt.Println("SIGNATURE INVALID")
		return "", err
	}
	fmt.Println("CLAIMS", claims)

	// isValidClaims(claims)

	expiration := time.Now().Add(5 * time.Minute)

	claims.Email = claims.Email
	claims.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
	}

	unsigned := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signed, err := unsigned.SignedString(private_key)
	if err != nil {
		return "", err
	}
	fmt.Println("NEW TOKEN ", signed)

	return signed, nil
}

// parse token, check if valid
// get claim info
// check claim info against db
func ValidateToken(jwtkn string) (bool, error) {
	fmt.Println("VALIDATE TOKEN ", jwtkn)
	// validate against local secret token file (pem file)
	claims := &jwtCustomClaims{}
	// https://godoc.org/github.com/dgrijalva/jwt-go#ParseRSAPublicKeyFromPEM
	token, err := jwt.ParseWithClaims(jwtkn, claims, func(t *jwt.Token) (interface{}, error) {
		return private_key, nil
	})
	// token, err := jwt.Parse(jwtkn, func(t *jwt.Token) (interface{}, error) {
	// 	return private_key, nil
	// })

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("SIGNATURE INVALID")
			// return false, err
		}
		// return false, err
	}
	fmt.Println("TOKEN VALID", token.Valid)
	if !token.Valid {
		// return false, err
	}
	fmt.Println("TOKEN", token)
	fmt.Println("CLAIMS", claims)
	// check if expired
	if time.Unix(claims.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return false, errors.New("EXPIRED TOKEN")
	}

	return token.Valid, nil
}

func ValidateCredentials(email string, password string) (bool, error) {
	dbCtrl := InitLocalDbConnection("http://localhost")
	// maybe returns list
	userList, err := dbCtrl.QueryUser(EMAIL, email)
	if err != nil {
		return false, err
	}
	fmt.Println("QUERY RESULT", userList)
	for _, v := range userList {
		if password == v.Password {
			return false, err
		}
	}
	return true, nil
}

func SaveCredentials(email string, password string) (bool, error) {
	// open dynamodb connection
	dbCtrl := InitLocalDbConnection("http://localhost")
	u := &UserObject{
		Email:    email,
		Password: password,
	}
	dbCtrl.PutItem("usertable", u)
	return false, nil
}
