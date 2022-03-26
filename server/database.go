package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"
	// "reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// func Insert(params ...string) {
// 	prmLength := len(params)

// }

func ConnectDB(db_url string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(db_url))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		fmt.Println("Cannot connect to the mongo server")
		log.Fatal(err)
	}
	fmt.Println("Connected to the database")
	return client
}

func checkEmail(email string, usersCollection *mongo.Collection) string {
	filter := bson.D{{"email", email}}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if fmt.Sprint(err) == "mongo: no documents in result" {
			return "No email"
		}
	}
	return "Email exists"
}

func checkUsername(username string, usersCollection *mongo.Collection) string {
	filter := bson.D{{"username", username}}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if fmt.Sprint(err) == "mongo: no documents in result" {
			return "No user"
		}
	}
	return "User exists"
}

func checkCreds(username string, password string, usersCollection *mongo.Collection) string {
	filter := bson.D{{"username", username}, {"password", password}}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if fmt.Sprint(err) == "mongo: no documents in result" {
			return "Invalid creds"
		}
	}
	return "Correct creds"
}

func HandleLogin(username string, password string, client *mongo.Client) string {
	usersCollection := client.Database("NarutoDB").Collection("users")
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hash)
	credStatus := checkCreds(username, hashedPassword, usersCollection)
	if credStatus == "Invalid creds" {
		return "Invalid credentials"
	}
	return "Logged in"
}

func HandleSignup(username string, password string, email string, client *mongo.Client) string {
	usersCollection := client.Database("NarutoDB").Collection("users")
	userStatus := checkUsername(username, usersCollection)
	emailStatus := checkEmail(email, usersCollection)
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x", hash)
	if emailStatus == "Email exists" {
		return "Email taken"
	}
	if userStatus == "User exists" {
		return "Username taken"
	}
	if userStatus == "No user" && emailStatus == "No email" {
		oneUser := bson.D{{"username", username}, {"password", hashedPassword}, {"email", email}}
		insertResult, err := usersCollection.InsertOne(context.TODO(), oneUser)
		if err != nil {
			return "Error"
		}
		fmt.Println(insertResult)
	}
	return "Done"
}

func AddExpense(username string, amount int, title string, description string, client *mongo.Client) string {
	usersCollection := client.Database("NarutoDB").Collection("expenses")
	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()
	oneUser := bson.D{{"username", username}, {"amount", amount}, {"title", title}, {"description", description}, {"day", day}, {"month", month}, {"year", year}}
	insertResult, err := usersCollection.InsertOne(context.TODO(), oneUser)
	if err != nil {
		return "Error"
	}
	fmt.Println(insertResult)
	return "Done"
}

func FetchExpenses(username string, year int, month int, client *mongo.Client)string {
	usersCollection := client.Database("NarutoDB").Collection("expenses")
	var filter bson.D
	if year == 0 && month == 0 {
		mth := time.Now().Month()
		filter = bson.D{{"username", username},{"month",mth}}
	} else if year == 0 && month != 0 {
		filter = bson.D{{"username", username}, {"month", month}}
	} else if year != 0 && month == 0 {
		filter = bson.D{{"username", username}, {"year", year}}
	} else {
		filter = bson.D{{"username", username}, {"year", year}, {"month", month}}
	}
	opts := options.Find().SetSort(bson.D{{"day",1}})
	data, err := usersCollection.Find(context.TODO(), filter,opts)
	var results []bson.M
	if err = data.All(context.TODO(), &results); err != nil {
		return "Internal server error"
	}
	if err != nil {
		return "Internal server error"
	}
	r1,err := json.Marshal(results)
	return string(r1)
}

func AddTranscation(username string, mode string, amount int, description string,client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("transactions")
	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()
	filter := bson.D{{"username",username},{"mode",mode},{"amount",amount},{"description",description},{"day",day},{"month",month},{"year",year}}
	insertResult,err := collcn.InsertOne(context.TODO(),filter)
	if err!=nil{
		return "Error"
	}
	fmt.Println(insertResult)
	return "Done"
}

func FetchTranscations(username string,mode string, year int, month int,client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("transactions")
	var filter bson.D
	if mode=="none"{
		if year==0&&month==0{
			filter = bson.D{{"username",username}}
		}else if year==0&&month!=0{
			filter = bson.D{{"username",username},{"month",month}}
		}else if year!=0&& month==0{
			filter = bson.D{{"username",username},{"year",year}}
		}else{
			filter = bson.D{{"username",username},{"month",month},{"year",year}}
		}
	}else if mode=="loan" || mode=="credit"{
		if year==0&&month==0{
			filter = bson.D{{"username",username},{"mode",mode}}
		}else if year==0&&month!=0{
			filter = bson.D{{"username",username},{"month",month},{"mode",mode}}
		}else if year!=0&& month==0{
			filter = bson.D{{"username",username},{"year",year},{"mode",mode}}
		}else{
			filter = bson.D{{"username",username},{"month",month},{"year",year},{"mode",mode}}
		}
	}
	
	opts := options.Find().SetSort(bson.D{{"day",1},{"month",1},{"year",1}})
	data,err := collcn.Find(context.TODO(),filter,opts)
	if err!=nil{
		log.Fatal(err)
	}
	var results [] bson.M
	if err = data.All(context.TODO(),&results); err!=nil{
		return "Internal server error"
	}
	r1,err := json.Marshal(results)
	return string(r1)
}