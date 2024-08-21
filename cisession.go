package cisession

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yvasiyarov/php_session_decoder"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

type SessionManager struct {
	sessionFileDir string
	sessId         string
	mu             sync.Mutex
	data           php_session_decoder.PhpSession
}

func NewSession(sessId string) (*SessionManager, error) {
	s := new(SessionManager)
	s.sessionFileDir = "/tmp"
	if len(sessId) < 1 {
		return s, fmt.Errorf("No Session Id provided")
	}
	s.sessId = sessId
	data, err := s.getSerializedDataFromFile()
	if err == nil {
		s.data, _ = s.decodeSerializedData(data)

	} else {
		return s, err
	}

	//sd.data = make(php_session_decoder.PhpSession)
	return s, nil
}

func (s *SessionManager) tokenFromCookie(c *gin.Context) string {
	cookie, err := c.Cookie("ci_session")
	if err != nil {
		return ""
	}
	return cookie
}

func (s *SessionManager) getSerializedDataFromFile() (string, error) {
	dat, err := os.ReadFile(s.sessionFileDir + "/ci_session" + s.sessId)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func (s *SessionManager) decodeSerializedData(data string) (php_session_decoder.PhpSession, error) {
	decoder := php_session_decoder.NewPhpDecoder(data)
	return decoder.Decode()
}

func (s *SessionManager) get(key string) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := s.getSerializedDataFromFile()
	if err == nil {
		s.data, _ = s.decodeSerializedData(data)
		if v, ok := (s.data)[key]; !ok {
			return nil, fmt.Errorf("session key not found")
		} else {
			return v, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func (s *SessionManager) SetFlash(key, value string) error {
	if key == "" {
		return fmt.Errorf("SetFlash: Key cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	s.updateCiVars(php_serialize.PhpArray{
		key: "new",
	})

	return nil
}

func (s *SessionManager) SetUserData(key, value string) error {
	if key == "" {
		return fmt.Errorf("SetUserData: Key cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return nil
}

func (s *SessionManager) GetUserData(key string) (string, error) {
	if key == "" {
		return "", fmt.Errorf("GetUserData: Key cannot be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[key]
	if ok {
		return val.(string), nil
	}
	return "", fmt.Errorf("Not Found")
}

func (s *SessionManager) updateCiVars(ciVars interface{}) {
	newCiVars := make(map[php_serialize.PhpValue]php_serialize.PhpValue)
	oldCiVars, ok := s.data["__ci_vars"]
	if ok {
		//copy current CI vars
		for k, v := range oldCiVars.(php_serialize.PhpArray) {
			newCiVars[k] = v.(php_serialize.PhpValue)
		}

	}
	for k, v := range ciVars.(php_serialize.PhpArray) {
		newCiVars[k] = v.(php_serialize.PhpValue)
	}
	s.data["__ci_vars"] = newCiVars
}

func (s *SessionManager) GetFlash(key string) string {
	if key == "" {
		return ""
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	flashVar, ok := s.data[key]
	if ok {
		s.updateCiVars(php_serialize.PhpArray{
			key: "old",
		})
		delete(s.data, key)
		return flashVar.(string)
	}
	return ""

}

func (s *SessionManager) Write() error {
	s.data["__ci_last_regenerate"] = time.Now().Unix()
	encoder := php_session_decoder.NewPhpEncoder(s.data)
	if result, err := encoder.Encode(); err == nil {
		err := os.WriteFile(s.sessionFileDir+"/ci_session"+s.sessId, []byte(result), 0666)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return err
	}
}
