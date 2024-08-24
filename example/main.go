package main

import (
	"log"
	"net/http"

	"github.com/csiwek/cisession"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	sessMiddleware, err := cisession.NewMiddleware(cisession.MiddlewareConfig{
		SessionDir:       "/tmp",
		SessionExpirySec: 1800,
		UnauthorizedFunc: Unauthorized,
		AuthorizerFunc:   Authorizer,
	})
	if err != nil {
		log.Fatal("Unable to create Sess Middleware: ", err)
	}
	r.GET("/loggedin", sessMiddleware.Middleware(), LoggedIn)
	r.POST("/logout", sessMiddleware.Middleware(), Logout)
	r.GET("/login", LoginPage)
	r.POST("/login", Login)
	r.Run(":8081")
}

// Function to handle unauthorized requests - can be used to redirect to login page
func Unauthorized(c *gin.Context, code int, message string) {
	c.Abort()
	c.Redirect(http.StatusFound, "/login")
}

// Authorizer function performes extra checks and can populate Gin's context with extra information - for example user id of logged user.i
// If err not nil request will not be authorized
func Authorizer(c *gin.Context) error {
	c.Set("userid", int64(123))
	return nil
}

func LoggedIn(c *gin.Context) {
	ci_session, exist := c.Get("ci_session")
	if !exist {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session, err := cisession.NewSession(ci_session.(string), cisession.SessionManagerConfig{})
	if err != nil {
		c.String(200, "could not create session error")
		return

	}
	defer session.Write()
	//userid it's the context's variable which was previously populated by the Authorizer func
	userId := c.GetInt64("userid")
	errorFlash := session.GetFlash("info")
	c.Header("Content-Type", "text/html")
	c.String(200, "<html><body><form action=\"/logout\" method=\"POST\"><input type=\"submit\" value=\"logout\"></form> <br> UserId: %d <br>Flash: %s <BR> Flash message will disappear on reload", userId, errorFlash)
}

func Logout(c *gin.Context) {
	ci_session, exist := c.Get("ci_session")
	if !exist {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	session, err := cisession.NewSession(ci_session.(string), cisession.SessionManagerConfig{})
	if err != nil {
		c.String(200, "could not create session error")
		return

	}
	defer session.Write()
	c.SetCookie("ci_session", "", 0, "/", "", false, true)
	c.Set("ci_session", nil)

	c.Redirect(http.StatusFound, "/login")
}

func LoginPage(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(200, "<html><body><form method=\"POST\"><input type=\"submit\" value=\"login\"></form> ")
}

func Login(c *gin.Context) {
	session, _ := cisession.CreateSession(cisession.SessionManagerConfig{})
	defer session.Write()
	session.SetUserData("my_login_session", "logged_in")
	session.SetFlash("info", "THIS IS A FLASH MESSAGE")
	c.SetCookie("ci_session", session.SessionId(), 1800, "/", "", false, true)
	c.Redirect(http.StatusFound, "/loggedin")
}
