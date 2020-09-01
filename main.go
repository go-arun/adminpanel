package main

import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/go-arun/adminpanel/modules/securepwd"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"go.mongodb.org/mongo-driver/bson"
	
)

//LoggedUserDetail .. to store details of users 
var LoggedUserDetail database.User // to Get values from DB
var zeroLoggedUserDetail database.User // to make above strcut empty sometimes 
var username string // is usfull while moving to update page , to store username before redirecting 

type values struct { // mdofy this name to an appropriate one TODO
	Name,AdmnButonVisibility string
}
//HomePageValues ... To pass values to Home Page 
var HomePageValues values 

var usrName,usrPwd string 
//LoginPageGet ...
func LoginPageGet(c *gin.Context) {
	//var recordFound bool
	fmt.Println("inside LoginpageGET....")
	recordFound,HomePageValues := getHomePageIfsessionActive(c)
	if (recordFound){
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

// To Load homme page , if there is an existing valid Sesion ID in DB
func getHomePageIfsessionActive(c *gin.Context)(recordFound bool,HomePageValues values) { 
	sessionCookie,err := c.Cookie ("user_sid_cookie")
	fmt.Println("err Val -->",err)
	if sessionCookie == "" { // no cookie found 
		return recordFound,HomePageValues // by default 'recordFound' val will be false
	}
	recordFound,LoggedUserDetail := database.TraceUserWithSID(sessionCookie)
	HomePageValues.Name = strings.ToUpper(LoggedUserDetail.Name)
	isAdmin := LoggedUserDetail.IsAdmn

	if (isAdmin){ // if admin make admin button visible
		HomePageValues.AdmnButonVisibility = "visible"
	}else {
		HomePageValues.AdmnButonVisibility = "hidden"
	}
	return recordFound,HomePageValues // Details of user to show in Home Page
}

//HomepagePost ...

  //LoginPagePost for Web
 func LoginPagePost(c *gin.Context){
	fmt.Println("Login page Post ----------->")
	recordFound,HomePageValues := getHomePageIfsessionActive(c)
	if (recordFound){ // If user is loged in then no need to show sign-in page ( while click on back)
		c.HTML(
			http.StatusOK,
			"home.html",
			HomePageValues,
		)
	}else{
		c.Request.ParseForm()
		for key, value := range c.Request.PostForm {
			if (key == "usrname") { // getting value from form
				usrName = value[0]
			}else{
				usrPwd = value[0]
			}
		}
		//usrPwd = 
		var usrExists bool
		usrExists,LoggedUserDetail = database.UserValidaiton(usrName,usrPwd)

		if (!usrExists){ // Login Error
			//c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
			//c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
			c.HTML(
				http.StatusOK,
				"login_err.html",
				gin.H{"title": "User Login"},
			)
		}else{ //Login Success
			//generating session ID using UUID 
			sessionID := database.AddSessionID(usrName) 
			c.SetCookie("user_sid_cookie",
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
			//c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
			//c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")

			if (!isAdmin){
				c.HTML(
					http.StatusOK,
					"home.html",
					HomePageValues,
				)
			}else{
				searchResults := database.FindAllUsers("")
			c.HTML(											
			http.StatusOK,
			"admnpanel.html",gin.H{
			"CollectedUserDetail": searchResults,
			})
			}
		}
	}
 }
 //AdmnpanelPost ...
 func AdmnpanelPost(c *gin.Context){
	c.Request.ParseForm()
	var searchResults []bson.M
//	searchKey := c.Request.PostForm["searchkey"][0];  // searchKey is userName
	

	operation :=  c.Request.PostForm["action"][0] 
	searchResults = database.FindAllUsers("") // always need to show full list afte modi del or update
	switch operation{ // based on action value
	case "find":
		searchKey := c.Request.PostForm["searchkey"][0] // Value in Find textbox
		// _,LoggedUserDetail = database.GetUsers(searchKey)
		searchResults = database.FindAllUsers(searchKey)
		fmt.Println("SearchKey------>",searchKey)
	case "del":
		database.DelUser(c.Request.PostForm["select"][0]) // uname of selected one 
		LoggedUserDetail = zeroLoggedUserDetail // if not made th
	case "modi":
		username = c.Request.PostForm["select"][0]
		c.Redirect(http.StatusMovedPermanently, "/update")
		c.Abort()
	case "logout":
		sidFromBrwser,_ := c.Cookie ("admin_sid_cookie")
		database.RemoveSessionID(sidFromBrwser)
		c.Redirect(http.StatusMovedPermanently, "/admin") // redirecting to admin loging page
		c.Abort()
		c.SetCookie("admin_sid_cookie", // Deleting cookie
		"",
		-1, // delete now !!
		"/",
		"",false,false,
		)
	}
	fmt.Println( "val of operation -->",operation)
	if (operation != "modi" && operation != "logout" ){ // redirection should happen if the operatiom  is not for modification and Logout
		//c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		//c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		//searchResults = database.FindAllUsers("") // added at the time of interview ..
		c.HTML(											
			http.StatusOK,
			"admnpanel.html",gin.H{
			"CollectedUserDetail": searchResults,
		})
	}	


}

  //SignupGet ...
func SignupGet( c *gin.Context){
	recordFound,HomePageValues := getHomePageIfsessionActive(c)
	if (recordFound){ // If user is loged in then no need to show sign-in page ( while click on back)
		c.HTML(
			http.StatusOK,
			"home.html",
			HomePageValues,
		)
	}else{
		c.HTML(
			http.StatusOK,
			"signup.html",
			[]string{"hidden"},
		)
	}

}
//SignupPost ...
func SignupPost(c *gin.Context){
	recordFound,HomePageValues := getHomePageIfsessionActive(c)
	if (recordFound){ // If user is loged in then no need to show sign-in page ( while click on back)
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
		c.Request.ParseForm()
		// for key,value := range c.Request.PostForm{
			name := c.Request.PostForm["name"][0]
			username := c.Request.PostForm["username"][0]
			email := c.Request.PostForm["email"][0]
			passwd1 := c.Request.PostForm["pwd1"][0]
			passwd1,_ = securepwd.HashPassword(passwd1) //hashing
			err := database.InsertRec(name,email,username,passwd1,false)
			if (err != nil ){
				c.HTML( // after addding new user direct to login page
					http.StatusOK,
					"signup.html",
					[]string{"visible"},
				)
			}else{
				c.HTML( // after addding new user direct to login page
				http.StatusOK,
				"index_login.html",
				gin.H{"title": "User Login"},
			)
		}
	}

}
//AdmnpanelGet ...
func AdmnpanelGet(c *gin.Context){ 
	if (admnSessinStatus(c)) {
		searchResults := database.FindAllUsers("")
		c.HTML(											
			http.StatusOK,
			"admnpanel.html",gin.H{
			"CollectedUserDetail": searchResults,
			})
	}else{
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "User Login"},
		)
	}
}
//UpdateGet ... 
func UpdateGet(c *gin.Context){
	if (admnSessinStatus(c)) {
		_,LoggedUserDetail = database.GetUser(username)
		c.HTML(
			http.StatusOK,
			"update.html",gin.H{
			"CollectedUserDetail": LoggedUserDetail,
		})
	}else{
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "User Login"},
		)
	}
}
//UpdatePost ... 
func UpdatePost(c *gin.Context){
	c.Request.ParseForm()
	// for key,value := range c.Request.PostForm{
		name := c.Request.PostForm["name"][0]
		// username := c.Request.PostForm["username"][0]
		email := c.Request.PostForm["email"][0]
		passwd1 := c.Request.PostForm["pwd1"][0]
		if (passwd1 != ""){ // Only if admin is changing pwd then only need to hash
			passwd1,_ = securepwd.HashPassword(passwd1) //hashing
		}
		database.UpdateRec(name,email,username,passwd1,false)
		
		searchResults := database.FindAllUsers("")
		c.HTML(											
		http.StatusOK,
		"admnpanel.html",gin.H{
		"CollectedUserDetail": searchResults,
		})
}
//HomepagePost ...
func HomepagePost(c *gin.Context){
	//Handling Log Out
	sidFromBrwser,_ := c.Cookie ("user_sid_cookie")
	database.RemoveSessionID(sidFromBrwser)
	c.Redirect(http.StatusMovedPermanently, "/") // redirecting to loging page
	c.Abort()
	c.SetCookie("user_sid_cookie", // Deleting cookie
	"",
	-1, // delete now !!
	"/",
	"",false,false,
	)
}
//
func HomepageGet(c *gin.Context) {
	//var recordFound bool
	fmt.Println("inside LoginpageGET....")
	recordFound,HomePageValues := getHomePageIfsessionActive(c)
	if (recordFound){
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
//AdminLoginGet ...
func AdminLoginGet( c *gin.Context){
	if (admnSessinStatus(c)) {
		searchResults := database.FindAllUsers("")
		c.HTML(											
			http.StatusOK,
			"admnpanel.html",gin.H{
			"CollectedUserDetail": searchResults,
			})
	}else{
		c.HTML(
			http.StatusOK,
			"admin_login.html",
			gin.H{"title": "User Login"},
		)
	}
}
//AdminLoginPost ...
func AdminLoginPost( c *gin.Context){
	sessionID := database.AddAdminSessionID() 
	c.SetCookie("admin_sid_cookie",
	sessionID,
	3600*12, // 12hrs
	"/",
	"",false,false, //domain excluded 
	)
	searchResults := database.FindAllUsers("")
	c.HTML(											
		http.StatusOK,
		"admnpanel.html",gin.H{
		"CollectedUserDetail": searchResults,
		},
	)
}

func admnSessinStatus(c *gin.Context)(sesStatus bool) { 
	sessionCookie,_ := c.Cookie ("admin_sid_cookie")
	if sessionCookie == "" { // no cookie found 
		return sesStatus // by default 'recordFound' val will be false
	} //cookie received frim boeser still we need to ensure from database too
	sesStatus,_ = database.TraceUserWithSID(sessionCookie)
	fmt.Println("Admin Session Exissts ?",sesStatus)
	return sesStatus // 
}

func main(){
	router := gin.Default()
	router.LoadHTMLGlob("htmls/*")

	router.GET("/", LoginPageGet)
	router.POST("/", LoginPagePost)
	router.POST("/home",HomepagePost)
	router.GET("/home",HomepageGet)
	router.GET("/signup",SignupGet)
	router.POST("/signup",SignupPost)
	router.GET("/admnpanel",AdmnpanelGet)
	router.POST("/admnpanel",AdmnpanelPost)
	router.GET("/update",UpdateGet)
	router.POST("/update",UpdatePost)
	router.GET("/admin",AdminLoginGet)
	router.POST("/admin",AdminLoginPost)


     router.Run()
}
