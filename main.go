package main


import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var loggedUserDetail database.User // to Get values from DB
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
	usrExists,loggedUserDetail = database.UserValidaiton(usrName,usrPwd)
	//fmt.Println("loggedUserDetail-->",loggedUserDetail.Name,usrExists)
	if (!usrExists){ // Login Error
		c.HTML(
			http.StatusOK,
			"login_err.html",
			gin.H{"title": "User Login"},
		)
	}else{ //Login Success
		HomePageValues.Name = strings.ToUpper(loggedUserDetail.Name)
		isAdmin := loggedUserDetail.IsAdmn
		if (isAdmin){
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
	c.HTML(
		http.StatusOK,
		"admnpanel.html",
		gin.H{"title": "Admin Panel"},
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


     router.Run()
}
