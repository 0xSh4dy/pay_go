package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)


func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

type login struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

type signup struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
	Email    string `valid:"Required; MaxSize(50)"`
}

func main() {
	err := godotenv.Load()
	if err!=nil{
		log.Fatal(err)
	}
	db_url := os.Getenv("db_url")
	client := ConnectDB(db_url)
	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Server is up and running")
	})

	router.POST("/api/login/", func(c *gin.Context) {
		var creds login
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		loginStatus := HandleLogin(creds.Username,creds.Password,client)
		fmt.Println(loginStatus)
		c.JSON(http.StatusOK, gin.H{"result": "success"})
	})
	router.POST("/api/register/", func(c *gin.Context) {
		var creds signup
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		fmt.Println(creds.Username)
		fmt.Println(creds.Password)
		fmt.Println(creds.Email)
		res := HandleSignup(creds.Username, creds.Password, creds.Email, client)
		fmt.Println(res)
		c.JSON(http.StatusOK, gin.H{"result": "success"})

	})
	fmt.Println("server running at port 7000")
	router.Run(":7000")
}
