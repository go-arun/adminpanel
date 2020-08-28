package main


import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)
//LoggedUserDetail .. to store details of users 
var LoggedUserDetail database.User // to Get values from DB
var zeroLoggedUserDetail database.User // to make above strcut empty sometimes 

type values struct {
	Name,AdmnButonVisibility string
}
//HomePageValues ... To pass values to Home Page 
var HomePageValues values

var usrName,usrPwd string 
//LoginPageGet ...
func LoginPageGet(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index_login.html",
		gin.H{"title": "User Login"},
	)
}
//HomepagePost ...
func HomepagePost(c *gin.Context){

	fmt.Println("sdsdsdsds")
}


  //LoginPagePost for Web
 func LoginPagePost(c *gin.Context)  {
	c.Request.ParseForm()
	for key, value := range c.Request.PostForm {
		fmt.Println(key,value)
		if (key == "usrname") {
			usrName = value[0]
		}else{
			usrPwd = value[0]
		}
	}
	var usrExists bool
	usrExists,LoggedUserDetail = database.UserValidaiton(usrName,usrPwd)
	//fmt.Println("LoggedUserDetail-->",LoggedUserDetail.Name,usrExists)
	if (!usrExists){ // Login Error
		c.HTML(
			http.StatusOK,
			"login_err.html",
			gin.H{"title": "User Login"},
		)
	}else{ //Login Success
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

	fmt.Println("Actiio->",c.Request.PostForm["action"][0])
	operation :=  c.Request.PostForm["action"][0] 
	
	switch operation{ // based on action value

	case "find":
		fmt.Println("Inside Fine")
		searchKey := c.Request.PostForm["searchkey"][0] // Value in Find textbox
		_,LoggedUserDetail = database.GetUsers(searchKey)
	case "del":
		database.DelUser(c.Request.PostForm["select"][0]) // uname of selected one 
		LoggedUserDetail = zeroLoggedUserDetail // if not made th
	
	}
	
	
	c.HTML(
		http.StatusOK,
		"admnpanel.html",gin.H{
		"LoggedUserDetail": LoggedUserDetail,
})

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


     router.Run()
}
