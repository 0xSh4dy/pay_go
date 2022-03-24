package main

import (
	"github.com/dgrijalva/jwt-go"
)

type Claims struct{
	Username string `valid:"Required; MaxSize(50)`
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
