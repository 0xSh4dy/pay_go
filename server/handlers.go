package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"fmt"
)

type Claims struct{
	Username string `valid:"Required; MaxSize(50)"`
	jwt.StandardClaims
}

func createToken(username string, jwt_key []byte)string{
	claims := &Claims{
		Username :username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenStr, err := token.SignedString(jwt_key)
	if err!=nil{
		return "Failed"
	}
	return tokenStr
}

func VerifyToken(token string,jwt_key []byte) string{
	claims := &Claims{}
	tok,err := jwt.ParseWithClaims(token, claims,func(t *jwt.Token)(interface{},error){
		return jwt_key,nil
	})
	if err!=nil{
		if err==jwt.ErrSignatureInvalid{
			return "Unauthorized access"
		}
		return "Bad request"
	}
	if !tok.Valid{
		return "Unauthorized access"
	}
	return claims.Username
}

func PasswordResetEmail(username string, email string,url string, api_key string){
	from := mail.NewEmail("Rakshit Awasthi","rakshitawasthi14@gmail.com")
	subject := "Reset Password"
	to := mail.NewEmail(username,email)
	htmlContent := fmt.Sprintf("<div><p>Visit the following link to change your password</p>%s</div>",url)
	plainContent := "Email for testing"
	message := mail.NewSingleEmail(from,subject,to,plainContent,htmlContent)
	client := sendgrid.NewSendClient(api_key)
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}