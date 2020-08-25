package main


import (
	"fmt"
	"github.com/go-arun/adminpanel/modules/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

var usrName,usrPwd string 
//LoginPage for Web
func LoginPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index_login.html",
		gin.H{"title": "User Login"},
	)
}

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
	fmt.Println("-----------------------")
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

func main(){
	//database.InsertRec("Arun","ar@ar2.com","kumarcok1","pwd2",true)
	router := gin.Default()
	router.LoadHTMLGlob("htmls/*")

	router.GET("/", LoginPage)
	router.POST("/", LoginPagePost)
	router.POST("/home",HomepagePost)

     router.Run()
}
