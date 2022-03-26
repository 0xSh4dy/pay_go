package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	// "reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

type expenses struct {
	Amount      string `valid:"Required; MaxSize(10)"`
	Title       string `valid:"Required; MaxSize(50)"`
	Description string `valid:"Required; MaxSize(50)"`
	Authtoken   string `valid:"Required; MaxSize(50)"`
}

type transactions struct{
	Amount string `valid:"Required; MaxSize(10)"`
	Title string `valid:"Required; MaxSize(50)"`
	Mode string `valid:"Required; MaxSize(10)"`
	Description string `valid:"Required; MaxSize(50)"`
	Authtoken string `valid:"Required MaxSize(200)"`
}

type cookies struct {
	Authcookie string `valid:"Required; MaxSize(50)"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	db_url := os.Getenv("db_url")
	jwt_secret := []byte(os.Getenv("jwt_secret"))
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
		loginStatus := HandleLogin(creds.Username, creds.Password, client)
		if loginStatus == "Invalid credentials" {
			fmt.Println("Invalid credentials")
			c.String(http.StatusOK, "Invalid username or password")
			return
		}
		token := createToken(creds.Username, jwt_secret)
		if token == "Failed" {
			c.String(http.StatusOK, "Internal server error")
			fmt.Print("Internal server error")
			return
		}
		c.String(http.StatusOK, token)
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
		if res == "Email taken" {
			c.String(http.StatusOK, "Email is already taken")
			return
		} else if res == "Username taken" {
			c.String(http.StatusOK, "Username is already taken")
			return
		} else if res == "Error" {
			c.String(http.StatusOK, "Internal server error")
		}
		c.String(http.StatusOK, "Successfully registered")

	})

	router.POST("/api/expenses/", func(c *gin.Context) {
		var expense expenses
		if err := c.ShouldBindJSON(&expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
			return
		}
		fmt.Println(expense.Amount)
		fmt.Println(expense.Title)
		fmt.Println(expense.Description)
		fmt.Println(expense.Authtoken)
		amount, err := strconv.Atoi(expense.Amount)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(amount)
		username := VerifyToken(expense.Authtoken, jwt_secret)
		fmt.Println(username)
		if username == "Unauthorized access" {
			c.String(http.StatusOK, "Unauthorized access")
			return
		} else if username == "Bad request" {
			c.String(http.StatusOK, "Bad request")
			return
		}
		fmt.Println(amount)
		fmt.Println(username)
		result := AddExpense(username, amount, expense.Title, expense.Description, client)
		if result == "Error" {
			c.String(http.StatusOK, "Internal server error")
		}
		c.String(http.StatusOK, "Added to database")
	})


	router.GET("/api/expenses/", func(c *gin.Context) {
		token := c.Query("token")
		month,err := strconv.Atoi(c.Query("month"))
		if err!=nil{
			fmt.Println(err)
			return
		}
		year,err := strconv.Atoi(c.Query("year"))
		if(err!=nil){
			fmt.Println(err)
			return
		}
		username := VerifyToken(token,jwt_secret)
		if username=="Unauthorized access"{
			c.String(http.StatusOK,"Unauthorized access")
			return 
		}else if username=="Bad request"{
			c.String(http.StatusOK,"Bad request")
			return
		}
		stats := FetchExpenses(username,year,month,client)
		c.String(http.StatusOK, stats)
	})

	router.POST("/api/transactions/",func(c *gin.Context){
		var transaction transactions
		err := c.ShouldBindJSON(&transaction)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"result": err.Error()})
		}
		fmt.Println(transaction)
		c.String(http.StatusOK,"saved")
	})
	
	router.GET("/api/transcations/",func(c *gin.Context){
		token := c.Query("token")
		mode := c.Query("mode")
		month,err := strconv.Atoi(c.Query("month"))
		if err!=nil{
			fmt.Println(err)
			return
		}
		year,err := strconv.Atoi(c.Query("year"))
		if err!=nil{
			fmt.Println(err)
			return
		}
		username := VerifyToken(token,jwt_secret)
		if username=="Unauthorized access"{
			c.String(http.StatusOK,"Unauthorized access")
			return 
		}else if username=="Bad request"{
			c.String(http.StatusOK,"Bad request")
			return 
		}
		stats := FetchTranscations(username,mode,year,month,client)
		c.String(http.StatusOK,stats)
	})
	fmt.Println("server running at port 7000")
	router.Run(":7000")
}
