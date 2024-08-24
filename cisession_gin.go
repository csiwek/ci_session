package cisession

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yvasiyarov/php_session_decoder"
)

type MiddlewareConfig struct {
	SessionDir       string
	SessionExpirySec int
	UnauthorizedFunc func(*gin.Context, int, string)
	AuthorizerFunc   func(*gin.Context) error
}

type GinMiddleware struct {
	sessionFileDir   string
	sessionExpirySec int
	authorizerFunc   func(*gin.Context) error
	unauthorizedFunc func(*gin.Context, int, string)
}

func NewMiddleware(conf MiddlewareConfig) (*GinMiddleware, error) {
	mw := new(GinMiddleware)
	if len(conf.SessionDir) < 1 {
		mw.sessionFileDir = "/tmp"
	} else {
		mw.sessionFileDir = conf.SessionDir
	}

	if conf.SessionExpirySec < 1 {
		mw.sessionExpirySec = 1800
	}
	mw.unauthorizedFunc = conf.UnauthorizedFunc
	mw.authorizerFunc = conf.AuthorizerFunc

	return mw, nil
}

func (mw *GinMiddleware) MiddlewareInit() error {

	return nil
}
func (mw *GinMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *GinMiddleware) middlewareImpl(c *gin.Context) {

	cookie := mw.tokenFromCookie(c)
	if cookie == "" {
		mw.unauthorizedFunc(c, http.StatusUnauthorized, "No cookie")
		return
	}
	c.Set("ci_session", cookie)
	err := mw.authorizerFunc(c)
	if err != nil {
		mw.unauthorizedFunc(c, http.StatusUnauthorized, fmt.Sprintf("Authorizer did not succeed %v", err))
		return

	}
	c.SetCookie("ci_session", cookie, mw.sessionExpirySec, "/", "", false, true)
	c.Next()
}

func (mw *GinMiddleware) tokenFromCookie(c *gin.Context) string {
	cookie, err := c.Cookie("ci_session")
	if err != nil {
		return ""
	}
	return cookie
}

func (mw *GinMiddleware) getSerializedData(cookie string) (string, error) {
	dat, err := os.ReadFile(mw.sessionFileDir + "/ci_session" + cookie)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func (mw *GinMiddleware) decodeSerializedData(data string) (string, error) {
	decoder := php_session_decoder.NewPhpDecoder(data)
	sessionDataDecoded, err := decoder.Decode()
	if err != nil {
		return "", err
	}
	if v, ok := (sessionDataDecoded)["user_session"]; !ok {
		return "", fmt.Errorf("session not found")
	} else {
		return v.(string), nil
	}
}
