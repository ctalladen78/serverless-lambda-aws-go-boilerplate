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

	// validate signature
	isValid, err := ValidateToken(token)
	if err != nil {
		if !isValid {
			return "", errors.New("INVALID TOKEN")
		}
		return "", errors.New("INVALID TOKEN")
	}
	_, err = ValidateCredentials(email, password)
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
	unsigned := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := unsigned.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["objectid"] = shortuuid.New()
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

// parse token, check if valid
// get claim info
// check claim info against db
func ValidateToken(jwtkn string) (bool, error) {
	fmt.Println("VALIDATE TOKEN ", jwtkn)
	// validate against local secret token file (pem file)
	claims := &jwtCustomClaims{}
	// https://github.com/juusechec/jwt-beego/blob/master/jwt.go
	token, err := jwt.Parse(jwtkn, func(token *jwt.Token) (interface{}, error) {
		return private_key, nil
	})
	// https://godoc.org/github.com/dgrijalva/jwt-go#ParseRSAPublicKeyFromPEM
	// token, err := jwt.ParseWithClaims(jwtkn, claims, func(t *jwt.Token) (interface{}, error) {
	// 	return private_key, nil
	// })
	// token, err := jwt.ParseWithClaims(
	// 	jwtkn, claims,
	// 	func(token *jwt.Token) (interface{}, error) {
	// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
	// 			return nil, errors.New("invalid token")
	// 		}
	// 		return private_key, nil
	// 	})

	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		fmt.Printf("CLAIMS %v %v", claims, claims.StandardClaims)
	}

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Println("SIGNATURE INVALID")
			return false, err
		}
		return false, err
	}
	fmt.Println("TOKEN VALID", token.Valid)
	fmt.Println("TOKEN", token)
	fmt.Println("CLAIMS", claims)
	// check if expired
	if time.Unix(claims.StandardClaims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return false, errors.New("EXPIRED TOKEN")
	}

	return token.Valid, nil
}
