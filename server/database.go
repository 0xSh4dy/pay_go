package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

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
		fmt.Println(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		// Can't connect to Mongo server
		fmt.Println("Cannot connect to the mongo server")
		fmt.Println(err)
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
	eventResult := AddEvent(username,"expenditure",description,client)
	if eventResult == "Error"{
		return "Error"
	}
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
	
	r1,err := json.Marshal(results)
	return string(r1)
}

func AddTranscation(username string, mode string,title string, amount int, description string,client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("transactions")
	day := time.Now().Day()
	month := time.Now().Month()
	year := time.Now().Year()
	filter := bson.D{{"username",username},{"mode",mode},{"title",title},{"amount",amount},{"description",description},{"day",day},{"month",month},{"year",year}}
	insertResult,err := collcn.InsertOne(context.TODO(),filter)
	if err!=nil{
		return "Error"
	}
	fmt.Println(insertResult)
	var eventResult string
	if mode=="credit"{
		eventResult = AddEvent(username,"new","Credit: "+strconv.Itoa(amount)+"--> " +description,client)
	}else if mode=="debt"{
		eventResult = AddEvent(username,"new","Debt: "+strconv.Itoa(amount)+"--> "+description,client)
	}
	if eventResult=="Error"{
		return "Error"
	}
	return "Done"
}

func FetchTransactions(username string,mode string, year int, month int,client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("transactions")
	var filter bson.D
	if mode=="none"{
		if year==0&&month==0{
			filter = bson.D{{"username",username}}
		}else if year==0&&month!=0{
			filter = bson.D{{"username",username},{"month",month}}
		}else if year!=0&& month==0{
			filter = bson.D{{"username",username},{"year",year}}
		}else if year!=0 &&month!=0{
			filter = bson.D{{"username",username},{"month",month},{"year",year}}
		}else{
			filter = bson.D{{"username",username}}
		}
	}else if mode=="debt" || mode=="credit"{
		if year==0&&month==0{
			filter = bson.D{{"username",username},{"mode",mode}}
		}else if year==0&&month!=0{
			filter = bson.D{{"username",username},{"month",month},{"mode",mode}}
		}else if year!=0&& month==0{
			filter = bson.D{{"username",username},{"year",year},{"mode",mode}}
		}else{
			filter = bson.D{{"username",username},{"month",month},{"year",year},{"mode",mode}}
		}
	}else{
		return "Invalid mode"
	}
	
	opts := options.Find().SetSort(bson.D{{"year",1},{"month",1},{"day",1}})
	data,err := collcn.Find(context.TODO(),filter,opts)
	if err!=nil{
		fmt.Println(err)
	}
	var results [] bson.M
	if err = data.All(context.TODO(),&results); err!=nil{
		return "Internal server error"
	}
	r1,err := json.Marshal(results)
	return string(r1)
}

func FetchTransactionById(id string,client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("transactions")
	s1,err := primitive.ObjectIDFromHex(id)
	if err!=nil{
		return "Error"
	}
	var result bson.M
	er := collcn.FindOne(context.TODO(),bson.M{"_id":s1}).Decode(&result)
	if er!=nil{
		return "Error"
	}
	r1,e := json.Marshal(result)
	if e!=nil{
		return "Error"
	}
	return string(r1)
}
func AddEvent(username string, mode string, description string, client *mongo.Client)string{
	eventsCollection := client.Database("NarutoDB").Collection("events")
	var x bson.D
	if mode=="new"{
		x = bson.D{{"username",username},{"description",description},{"key",0}}
	}else if mode=="creditClear"{
		x = bson.D{{"username",username},{"description",description},{"key",1}}
	}else if mode=="debtClear"{
		x = bson.D{{"username",username},{"description",description},{"key",2}}
	}else if mode=="paidSome"{
		x = bson.D{{"username",username},{"description",description},{"key",3}}
	}else if mode=="receivedSome"{
		x = bson.D{{"username",username},{"description",description},{"key",4}}
	}else if mode=="expenditure"{
		x = bson.D{{"username",username},{"description",description},{"key",5}}
	}
	insertResult,err := eventsCollection.InsertOne(context.TODO(),x)
	if err!=nil{
		return "Error"
	}
	fmt.Println(insertResult)
	return "Saved"
}

func UpdateAmount(newAmount string,id string, client *mongo.Client)string{
	fmt.Println(newAmount)
	fmt.Println(id)
	collcn := client.Database("NarutoDB").Collection("transactions")
	s1,err := primitive.ObjectIDFromHex(id)
	if err!=nil{
		return "Error"
	}
	updateRes,er := collcn.UpdateOne(context.TODO(),bson.M{"_id":s1},bson.D{{"$set",bson.D{{"amount",newAmount}}}})
	if er!=nil{
		fmt.Println(er)
		return "Error"
	}
	fmt.Println(updateRes)
	return "DONE"
}

func FetchEvents(username string, client *mongo.Client)string{
	evtColl := client.Database("NarutoDB").Collection("events")
	opts := options.Find().SetLimit(20)
	result,err := evtColl.Find(context.TODO(),bson.M{"username":username},opts)
	if err!=nil{
		fmt.Println(err)
		return "Error"
	}
	var results [] bson.M
	if err = result.All(context.TODO(), &results); err != nil {
		return "Internal server error"
	}
	
	r1,err := json.Marshal(results)
	return string(r1)
}

func ChangePassword(password string, email string, client *mongo.Client)string{
	collcn := client.Database("NarutoDB").Collection("users")
	hash := sha256.Sum256([]byte(password))
	hashedPassword := fmt.Sprintf("%x",hash)
	fmt.Println(hashedPassword)
	var result bson.M
	err := collcn.FindOne(context.TODO(),bson.M{"email":email}).Decode(&result)
	fmt.Println(email)
	if err!=nil{
		fmt.Println(err)
		return "Invalid email"
	}

	updateResult,er := collcn.UpdateOne(context.TODO(),bson.M{"email":email},bson.D{{"$set",bson.M{"password":hashedPassword}}})
	if er!=nil{
		fmt.Println(er)
		return "Error"
	}
	fmt.Println(updateResult)
	return "Done"
}