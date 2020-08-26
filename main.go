package main


import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	usrExists := database.UserValidaiton(usrName,usrPwd)
	if (!usrExists){ // Login Error
		c.HTML(
			http.StatusOK,
			"login_err.html",
			gin.H{"title": "User Login"},
		)
	}else{ //Login Success
		c.HTML(
			http.StatusOK,
			"home.html",
			gin.H{"title": "User Login"},
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
func main(){
	//database.InsertRec("Arun","ar@ar2.com","kumarcok1","pwd2",true)
	router := gin.Default()
	router.LoadHTMLGlob("htmls/*")

	router.GET("/", LoginPageGet)
	router.POST("/", LoginPagePost)
	router.POST("/home",HomepagePost)
	router.GET("/signup",SignupGet)
	router.POST("/signup",SignupPost)


     router.Run()
}
