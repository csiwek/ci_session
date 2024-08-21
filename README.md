# CI Session

## Purpose

This library provides a way access and manipulate Sessions created by Codeigniter PHP framework.
It can help if you are planning to move to golang so part of the website can still operate on Codeigniter and you can gradually migrate your code to Golang.
The library also provides Middleware for Gin framework

## Supported Functions

- Has been tested with Codeigniter 3.0
- Provides Flash functions
- Provides Get/Set Userdata
- Support only sessions stored in files. If you are currently storing sessions in database you either need to switch to files or add functionality to store sessions in required storage.
- Provides Middleware for Gin

## Examples

### Gin Middleware

```
import (
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
        })
        if err != nil {
                log.Fatal("Unable to create Sess Middleware: ", err)
        }
        r.Use(sessMiddleware.Middleware())
	r.GET("/test", TestFunc)
	r.Run(":80")
}

func Unauthorized(c *gin.Context, code int, message string) {
        c.Abort()
        c.Redirect(http.StatusFound, "/login")
}

```


### Handler func

```
func TestFunc(c *gin.Context) {
        ci_session, exist := c.Get("ci_session")
        if !exist {
                c.Redirect(http.StatusFound, "/login") 
                return
        }
        session, err := cisession.NewSession(ci_session.(string))
        if err != nil {
                c.String(200, "could not create session error")
                return

        }
	defer session.Write()	
        session.SetFlash("info", "Session found")
	errorFlash := session.GetFlash("error")
        c.String(200, "Session found, error flash" + errorFlash)
}

```
