package main

import (
	"crypto/sha256"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func checkCreds(username string, password string, usersCollection *mongo.Collection)string{
	filter := bson.D{{"username",username},{"password",password}}
	var result bson.M
	err := usersCollection.FindOne(context.TODO(),filter).Decode(&result)
	if err!=nil{
		if fmt.Sprint(err)=="mongo: no documents in result"{
			return "Invalid creds"
		}
	}
	return "Correct creds"
}


func HandleLogin(username string, password string, client *mongo.Client)string{
	usersCollection := client.Database("NarutoDB").Collection("users")
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x",hash)
	credStatus := checkCreds(username,hashedPassword,usersCollection)
	if credStatus == "Invalid creds"{
		return "Invalid credentials"
	}
	return "Logged in"
}

func AddExpense(username string, amount int, title string, description string, client *mongo.Client)string{
	usersCollection := client.Database("NarutoDB").Collection("expenses")
	oneUser := bson.D{{"username",username},{"amount",amount},{"title",title},{"description",description}}
	insertResult,err := usersCollection.InsertOne(context.TODO(),oneUser)
	if err!=nil{
		return "Error"
	}
	fmt.Println(insertResult)
	return "Done"
}

func HandleSignup(username string, password string, email string, client *mongo.Client) string {
	usersCollection := client.Database("NarutoDB").Collection("users")
	userStatus := checkUsername(username,usersCollection)
	emailStatus := checkEmail(email,usersCollection)
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x",hash)
	if emailStatus=="Email exists"{
		return "Email taken"
	}
	if userStatus=="User exists"{
		return "Username taken"
	}
	if userStatus=="No user" && emailStatus=="No email"{
		oneUser := bson.D{{"username", username}, {"password", hashedPassword}, {"email", email}}
		insertResult, err := usersCollection.InsertOne(context.TODO(), oneUser)
		if err != nil {
			return "Error"
		}
		fmt.Println(insertResult)
			}
	return "Done"
}
