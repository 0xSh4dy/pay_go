package main

import (
	"encoding/json"
	"fmt"
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

type changes struct{
	Email string `valid:"Required; MaxSize(50)"`
	Username string `valid:"Required; MaxSize(50)"`
}

type newAmount struct{
	Amount string
	Token string
	ID string
}

type newPassword struct{
	Password string `valid:"Required"`
	Token string `valid:"Required"`
}

type evtStruct struct{
	_id string
	Amount string
	Day int
	Description string
	Mode string
	Month int
	Title string
	Username string 
	Year int
}
func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	email_api_key := os.Getenv("email_api_key")
	fmt.Println(email_api_key)
	db_url := os.Getenv("db_url")
	jwt_secret := []byte(os.Getenv("jwt_secret"))
	client := ConnectDB(db_url)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
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
			fmt.Println(err)
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

	router.POST("/api/transaction/",func(c *gin.Context){
		var transaction transactions
		err := c.ShouldBindJSON(&transaction)
		if err!=nil{
			c.JSON(http.StatusBadRequest,gin.H{"result": err.Error()})
		}
		tkn := transaction.Authtoken
		chk := VerifyToken(tkn,jwt_secret)
		if chk=="Unauthorized access"{
			c.String(http.StatusOK,"Unauthorized access")
			return 
		}else if chk=="Bad request"{
			c.String(http.StatusOK,"Bad request")
			return
		}
		amnt,err := strconv.Atoi(transaction.Amount)
		if err!=nil{
			fmt.Println(err)
		}
		res := AddTranscation(chk,transaction.Mode,transaction.Title,amnt,transaction.Description,client)
		if res=="Error"{
			c.String(http.StatusOK,"Internal server error")
			return
		}

		c.String(http.StatusOK,"Saved")
	})
	
	router.GET("/api/transaction/",func(c *gin.Context){
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
			c.String(http.StatusUnauthorized,"Unauthorized access")
			return 
		}else if username=="Bad request"{
			c.String(http.StatusBadRequest,"Bad request")
			return 
		}
		stats := FetchTransactions(username,mode,year,month,client)
		c.String(http.StatusOK,stats)
	})
	
	router.GET("/api/transaction/:username",func(c *gin.Context){
		id := c.Query("_id")
		token := c.Query("token")
		username := VerifyToken(token,jwt_secret)
		if username=="Unauthorized access"{
			c.String(http.StatusUnauthorized,"Unauthorized access")
			return
		}else if username=="Bad request"{
			c.String(http.StatusBadRequest,"Bad request")
			return
		}
		fmt.Println(id)
		fmt.Println(username)
		fetchedData := FetchTransactionById(id,client)
		if fetchedData=="Error"{
			c.String(http.StatusNotFound,"Not found")
			return
		}
		c.String(http.StatusOK,fetchedData)
	})

	router.PATCH("/api/transaction/",func(c *gin.Context){
		var newAmnt newAmount
		err := c.ShouldBindJSON(&newAmnt)
		if err!=nil{
			c.String(http.StatusBadRequest,"Invalid format")
			return
		}
		token := newAmnt.Token
		amt := newAmnt.Amount
		id := newAmnt.ID
		username := VerifyToken(token,jwt_secret)
		if username=="Bad request"{
			c.String(http.StatusBadRequest,"Bad request")
			return
		}else if username=="Unauthorized access"{
			c.String(http.StatusUnauthorized,"Unauthorized access")
			return
		}
		var evStruct evtStruct
		er := json.Unmarshal([]byte(FetchTransactionById(id,client)),&evStruct)
		if er!=nil{
			fmt.Println(er)
		}

		if amt=="0"{
			
			if evStruct.Mode == "credit"{
				AddEvent(evStruct.Username,"creditClear",evStruct.Description,client)
			}else if evStruct.Mode == "debt"{
				AddEvent(evStruct.Username,"debtClear",evStruct.Description,client)
			}

		}else{
			if evStruct.Mode == "credit"{
				AddEvent(evStruct.Username,"receivedSome",evStruct.Description,client)
			}else if evStruct.Mode=="debt"{
				AddEvent(evStruct.Username,"paidSome",evStruct.Description,client)
			}
		}
		r1 := UpdateAmount(amt,id,client)
		if r1=="Error"{
			c.String(http.StatusInternalServerError,"Internal server error")
			return
		}
		c.String(http.StatusOK,"Updated")
	})

	router.GET("/api/dashboard",func(c* gin.Context){
		token := c.Query("token")
		username := VerifyToken(token,jwt_secret)
		if username=="Unauthorized access"{
			c.String(http.StatusUnauthorized,"Unauthorized access")
			return
		}else if username=="Bad request"{
			c.String(http.StatusBadRequest,"Bad request")
			return
		}
		fetchedData := FetchEvents(username,client)
		if fetchedData=="Error" || fetchedData=="Internal server error"{
			c.String(http.StatusInternalServerError,"Internal server error")
			return
		}
		c.String(http.StatusOK,fetchedData)
	})

	router.POST("/api/change/",func(c* gin.Context){
		var change changes
		err := c.ShouldBindJSON(&change)
		if err!=nil{
			c.String(http.StatusOK,"Invalid JSON format")
			return
		}
		api_key := os.Getenv("email_api_key")
		token := createToken(change.Email,jwt_secret)
		if token=="Failed"{
			c.String(http.StatusInternalServerError,"Internal server error")
			return
		}
		url := fmt.Sprintf("http://127.0.0.1:7000/api/change?token=%s",token)
		PasswordResetEmail(change.Username,change.Email,url, api_key)
		fmt.Println(change.Email)
		c.String(http.StatusOK,"alright")
	})

	router.GET("/api/change",func(c *gin.Context){
		token := c.Query("token")
		email := VerifyToken(token,jwt_secret)
		if email=="Unauthorized access"{
			c.String(http.StatusUnauthorized,"Unauthorized")
			return
		}else if email=="Bad request"{
			c.String(http.StatusBadRequest,"Bad request")
			return
		}
		c.HTML(http.StatusOK,"changePass.html",gin.H{"title":"Change Password"})
	})

	router.POST("/api/newpass/",func(c* gin.Context){
		var newpass newPassword
		err := c.ShouldBindJSON(&newpass)
		if err!=nil{
			c.String(http.StatusBadRequest,"Bad request")
			return
		}
		email := VerifyToken(newpass.Token,jwt_secret)
		if email=="Unauthorized access"{
			c.String(http.StatusUnauthorized,"Unauthorized")
		}else if email=="Bad request"{
			c.String(http.StatusBadRequest,"Bad Request")
		}
		res := ChangePassword(newpass.Password,email,client)
		if res=="Error"{
			c.String(http.StatusInternalServerError,"Internal server error")
			return
		}else if res=="Invalid email"{
			c.String(http.StatusNotFound,"Invalid email")
			return
		}
		c.String(http.StatusOK,"Done")
	})
	fmt.Println("server running at port 7000")
	router.Run(":7000")
}
