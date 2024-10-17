package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
)

// RequestIDKey ...
const (
	RequestIDKey       = "request_id"
	RequestIDHeaderKey = "X-Request-Id"

	ClientIDKey = "client_id"

	ErrorIDKey = "err"

	// NeverExpiresUnixTime 永久有效期，使用2100.01.01 00:00:00 的unix time作为永久有效期的表示，单位秒
	// time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	NeverExpiresUnixTime = 4102444800
)

// ErrNilRequestBody ...
var ErrNilRequestBody = errors.New("request Body is nil")

// ReadRequestBody will return the body in []byte, without change the origin body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, ErrNilRequestBody
	}

	body, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, err
}

// GetRequestID ...
func GetRequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// SetRequestID ...
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(RequestIDKey, requestID)
}

// GetClientID ...
func GetClientID(c *gin.Context) string {
	return c.GetString(ClientIDKey)
}

// SetClientID ...
func SetClientID(c *gin.Context, clientID string) {
	c.Set(ClientIDKey, clientID)
}

// GetError ...
func GetError(c *gin.Context) (interface{}, bool) {
	return c.Get(ErrorIDKey)
}

// SetError ...
func SetError(c *gin.Context, err error) {
	c.Set(ErrorIDKey, err)
}

// BasicAuthAuthorizationHeader ...
func BasicAuthAuthorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(StringToBytes(base))
}

func BindJson(r *http.Request, ptr any) error {
	targetValue := reflect.ValueOf(ptr)
	if targetValue.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	if targetValue.IsNil() {
		return errors.New("target cannot be a nil pointer")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, ptr); err != nil {
		return err
	}

	return nil
}

func ResponseJson(w *http.ResponseWriter, code int, data any) {
	(*w).Header().Set("Content-Type", "application/json")
	if code != 0 {
		(*w).WriteHeader(code)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		_, _ = (*w).Write([]byte(fmt.Sprintf("%v", data)))
		return
	}

	_, _ = (*w).Write(jsonBytes)
}
