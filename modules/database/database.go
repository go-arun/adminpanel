package database

import (
    "context"
    "fmt"
    "log"

   "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

//User ...
type User struct {
    Name,Email,Pwd,Usrnm string
    IsAdmn bool
}
//Credentials ...
type Credentials struct{
    Unm,Pwd string
}

var Collection *mongo.Collection

var ClientOptions *options.ClientOptions

//UserValidaiton in LoginPage
func UserValidaiton(uname,pwd string)(bool){
    var result User
    filter := bson.M{"usrnm" : uname,"pwd" : pwd}
    
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    err = Collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
         return false // no such user 
    }
    fmt.Printf("Found a single document: %+v\n", result)
    return true // yes credetial exists 

}
//InsertRec to insert into database
func InsertRec(name,email,usrName,pwd string,adminStatus bool){
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    var usr User
    usr.Name = name
    usr.Usrnm = usrName
    usr.Email = email
    usr.Pwd = pwd 
    usr.IsAdmn = adminStatus

	insertResult, err := Collection.InsertOne(context.TODO(), usr)
	if err != nil {
    	log.Fatal(err)
    }
fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
