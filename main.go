package main

import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
)

//LoggedUserDetail .. to store details of users 
var LoggedUserDetail database.User // to Get values from DB
var zeroLoggedUserDetail database.User // to make above strcut empty sometimes 
var username string // is usfull while moving to update page , to store username before redirecting 

type values struct {
	Name,AdmnButonVisibility string
}
//HomePageValues ... To pass values to Home Page 
var HomePageValues values

var usrName,usrPwd string 
//LoginPageGet ...
func LoginPageGet(c *gin.Context) {
	isLoggedIN,_ := c.Cookie ("sid_cookie")
	if ( isLoggedIN != "" ){ // User not yet logged out so ,redirect to home page
		_,LoggedUserDetail = database.TraceUserWithSID(isLoggedIN)
		HomePageValues.Name = strings.ToUpper(LoggedUserDetail.Name)
		isAdmin := LoggedUserDetail.IsAdmn
		if (isAdmin){ // if admin make admin button visible
			HomePageValues.AdmnButonVisibility = "visible"
		}else {
			HomePageValues.AdmnButonVisibility = "hidden"
		}
		c.HTML(
			http.StatusOK,
			"home.html",
			HomePageValues,
		)
	}else{
	c.HTML(
		http.StatusOK,
		"index_login.html",
		gin.H{"title": "User Login"},
	)
	}
}
//HomepagePost ...

  //LoginPagePost for Web
 func LoginPagePost(c *gin.Context)  {
	c.Request.ParseForm()
	for key, value := range c.Request.PostForm {
		fmt.Println(key,value)
		if (key == "usrname") { // getting value from form
			usrName = value[0]
		}else{
			usrPwd = value[0]
		}
	}
	var usrExists bool
	usrExists,LoggedUserDetail = database.UserValidaiton(usrName,usrPwd)

	if (!usrExists){ // Login Error
		c.HTML(
			http.StatusOK,
			"login_err.html",
			gin.H{"title": "User Login"},
		)
	}else{ //Login Success
		//generating session ID using UUID 
		sessionID := database.AddSessionID(usrName) 
		c.SetCookie("sid_cookie",
		sessionID,
		3600*12, // 12hrs
		"/",
		"",false,false, //domain excluded 
		)
		//c.SetCookie("cookieName", "testCookie", 100000, "/", "", false, false)
		HomePageValues.Name = strings.ToUpper(LoggedUserDetail.Name)
		isAdmin := LoggedUserDetail.IsAdmn
		if (isAdmin){ // if admin make admin button visible
			HomePageValues.AdmnButonVisibility = "visible"
		}else {
			HomePageValues.AdmnButonVisibility = "hidden"
		}
		 
		c.HTML(
			http.StatusOK,
			"home.html",
			HomePageValues,
		)
	}
 }
 //AdmnpanelPost ...
 func AdmnpanelPost(c *gin.Context){
	c.Request.ParseForm()
	var searchResults []bson.M

	fmt.Println("Actiio->",c.Request.PostForm["action"][0])
	operation :=  c.Request.PostForm["action"][0] 
	
	switch operation{ // based on action value
	case "find":
		searchKey := c.Request.PostForm["searchkey"][0] // Value in Find textbox
		// _,LoggedUserDetail = database.GetUsers(searchKey)
		searchResults = database.FindAllUsers(searchKey)
	case "del":
		fmt.Println("selcted to del ->",c.Request.PostForm["select"][0])
		database.DelUser(c.Request.PostForm["select"][0]) // uname of selected one 
		LoggedUserDetail = zeroLoggedUserDetail // if not made th
	case "modi":
		username = c.Request.PostForm["select"][0]
		c.Redirect(http.StatusMovedPermanently, "/update")
		c.Abort()
	}

	if (operation != "modi"){ // if it is modification code is redirected to update page, so this lines need to be skipped
		c.HTML(
			http.StatusOK,
			"admnpanel.html",gin.H{
			"CollectedUserDetail": searchResults,
		})
	}	

}

  //SignupGet ...
func SignupGet( c *gin.Context){
	c.HTML(
		http.StatusOK,
		"signup.html",
		gin.H{"title": "User Login"},
	)

}
//SignupPost ...
func SignupPost(c *gin.Context){
	c.Request.ParseForm()
	// for key,value := range c.Request.PostForm{
		name := c.Request.PostForm["name"][0]
		username := c.Request.PostForm["username"][0]
		email := c.Request.PostForm["email"][0]
		passwd1 := c.Request.PostForm["pwd1"][0]
		passwd2  := c.Request.PostForm["pwd2"][0]
		fmt.Println(name,username,email,passwd1,passwd2)
		database.InsertRec(name,email,username,passwd1,false)
	// }

}
//AdmnpanelGet ...
func AdmnpanelGet(c *gin.Context){
	fmt.Println("OKOKOKOK",LoggedUserDetail.Usrnm)
	c.HTML(
		http.StatusOK,
		"admnpanel.html",gin.H{
		"LoggedUserDetail": LoggedUserDetail,
	})
}
//UpdateGet ... 
func UpdateGet(c *gin.Context){
	_,LoggedUserDetail = database.GetUser(username)
	c.HTML(
		http.StatusOK,
		"update.html",gin.H{
		"CollectedUserDetail": LoggedUserDetail,
	})
}
//UpdatePost ... 
func UpdatePost(c *gin.Context){
	c.Request.ParseForm()
	// for key,value := range c.Request.PostForm{
		name := c.Request.PostForm["name"][0]
		// username := c.Request.PostForm["username"][0]
		email := c.Request.PostForm["email"][0]
		passwd1 := c.Request.PostForm["pwd1"][0]
		database.UpdateRec(name,email,username,passwd1,false)
		
		c.HTML( // after updation loading admin page
			http.StatusOK,
			"admnpanel.html",gin.H{
			"LoggedUserDetail": LoggedUserDetail,
		})
}
//HomepagePost ...
func HomepagePost(c *gin.Context){
	//Handling Log Out
	sidFromBrwser,_ := c.Cookie ("sid_cookie")
	database.RemoveSessionID(sidFromBrwser)
	c.Redirect(http.StatusMovedPermanently, "/") // redirecting to loging page
	c.Abort()
	fmt.Println(" Home Page post triggered !!")
	c.SetCookie("sid_cookie", // Deleting cookie
	"",
	-1, // delete now !!
	"/",
	"",false,false,
	)
}

func main(){
	//database.InsertRec("Arun","ar@ar2.com","kumarcok1","pwd2",true)
	router := gin.Default()
	router.LoadHTMLGlob("htmls/*")

	router.GET("/", LoginPageGet)
	router.POST("/", LoginPagePost)
	router.POST("/home",HomepagePost)
	router.GET("/signup",SignupGet)
	router.POST("/signup",SignupPost)
	router.GET("/admnpanel",AdmnpanelGet)
	router.POST("/admnpanel",AdmnpanelPost)
	router.GET("/update",UpdateGet)
	router.POST("/update",UpdatePost)


     router.Run()
}
