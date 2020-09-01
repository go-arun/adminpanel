package database

import (
    "context"
    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "crypto/rand"
    "io"
    "github.com/go-arun/adminpanel/modules/securepwd"
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

//RemoveSessionID ... 
func RemoveSessionID(sid string) {
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    filter := bson.D{{"sess_id", sid}}
        update := bson.D{{"$set", bson.D{
            {"sess_id", ""},// Made it empty
        }}}
    
	updateResult, err := Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
    	log.Fatal(err)
    }
    fmt.Println("SessID Removed-",updateResult)
}

//AddSessionID ... 
func AddSessionID(usrName string)(string) {
    sessID,_ := generateNewUUID()
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    filter := bson.D{{"usrnm", usrName}}
        update := bson.D{{"$set", bson.D{
            {"sess_id", sessID},//UUID
        }}}
    
	updateResult, err := Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
    	log.Fatal(err)
    }
    fmt.Println("SessID Added-",updateResult)
    return sessID
}
//AddAdminSessionID ... to inset admin sesid to db
func AddAdminSessionID()(string) {
    sessID,_ := generateNewUUID()
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    filter := bson.D{{"isadmn", true}}
        update := bson.D{{"$set", bson.D{
            {"sess_id", sessID},//UUID
        }}}
    
	updateResult, err := Collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
    	log.Fatal(err)
    }
    fmt.Println("SessID Added-",updateResult)
    return sessID
}

//DelUser .. 
func DelUser(uname string){
	
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    
    filter := bson.M{"usrnm" : uname}
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    Collection.DeleteOne(context.TODO(), filter)
    // Collection.DeleteOne(context.TODO(), filter).Decode(&result)
}
//FindAllUsers ...
func FindAllUsers(name string)([]bson.M){

    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    fmt.Println("Inside FindAllUsers")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    filter := bson.D{{"name", primitive.Regex{Pattern: "^"+name}}}
    Collection := client.Database("admnpanel").Collection("users")  
  
    cursor, err := Collection.Find(context.TODO(), filter)
    if err != nil {
        log.Fatal(err)
    }

    // get a list of all returned documents and print them out
    // see the mongo.Cursor documentation for more examples of using cursors
    var results []bson.M
    if err = cursor.All(context.TODO(), &results); err != nil {
        log.Fatal(err)
    }
    // for _, result := range results {
    //     fmt.Println(result)
    // }

    return results
}

//GetUser ... to get user , to pick only one , used while updating only
func GetUser(usrName string)(bool,User){

    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    var result User
    filter := bson.M{"usrnm" : usrName}
    
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
//TraceUserWithSID ... 
func TraceUserWithSID(receivedCookie string)(bool,User){

    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    var result User
    filter := bson.M{"sess_id" : receivedCookie }
    
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")
    err = Collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil  {
        fmt.Printf("No mactching record found for this SID !!")
        return false,result // no rec found with this sec-id
    }
    fmt.Printf("Found a single document: %+v\n", result)
    return true,result // found one 
}
//UserValidaiton in LoginPage
func UserValidaiton(uname,pwd string)(bool,User){
    var result User
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    
    filter := bson.M{"usrnm" : uname}
    
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
    //User exists , so will check pwd with hash too 
    match := securepwd.CheckPasswordHash(pwd, result.Pwd) // will reutn true / false
    fmt.Println("-------------------------------------------->Password match result",match,pwd)

    return match,result // yes credetial exists 
}

//InsertRec to insert into database
func InsertRec(name,email,usrName,pwd string,adminStatus bool)(error){
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
    	return err
    }
fmt.Println("Inserted a single document: ", insertResult.InsertedID)
return nil
}

//UpdateRec to insert into database
func UpdateRec(name,email,usrName,pwd string,adminStatus bool){
    ClientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(context.TODO(), ClientOptions)
    if err != nil {
        log.Fatal(err)
    }
    Collection := client.Database("admnpanel").Collection("users")

    filter := bson.D{{"usrnm", usrName}}
    if (pwd != ""){
        update := bson.D{{"$set", bson.D{
            {"email", email},
            {"name",name},
            {"pwd",pwd},
        }}}
        _, err := Collection.UpdateOne(context.TODO(), filter, update)
        if err != nil {
            log.Fatal(err)
        }
    }else{ // if passwd is empty means, keep it same
        update := bson.D{{"$set", bson.D{
            {"email", email},
            {"name",name},
        }}}
        _, err := Collection.UpdateOne(context.TODO(), filter, update) //Update result ignored 
        if err != nil {
            log.Fatal(err)
        }

    }

}

func generateNewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


