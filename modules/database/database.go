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
// var result User // to pass to main
//Credentials ...
type Credentials struct{
    Unm,Pwd string
}

var Collection *mongo.Collection

var ClientOptions *options.ClientOptions

//DelUser .. 
func DelUser(uname string){
    
    filter := bson.M{"usrnm" : uname}
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    Collection.DeleteOne(context.TODO(), filter)
    // Collection.DeleteOne(context.TODO(), filter).Decode(&result)
}

//GetUsers ... to get users based on search criteria
func GetUsers(name string)(bool,User){
    var result User
    // name = "/^" + name + "/"
    // fmt.Println ( "uname modified to ->",name)
    filter := bson.M{"name" : name}
    
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    err = Collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
         return false,result // no such user 
    fmt.Println("No Match Found !!!")

    }
    fmt.Printf("Found a single document: %+v\n", result)
    return true,result // yes mach/es found
}


//UserValidaiton in LoginPage
func UserValidaiton(uname,pwd string)(bool,User){
    var result User
    
    filter := bson.M{"usrnm" : uname,"pwd" : pwd}
    
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    err = Collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
         return false,result // no such user 
    }
    fmt.Printf("Found a single document: %+v\n", result)
    return true,result // yes credetial exists 
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
