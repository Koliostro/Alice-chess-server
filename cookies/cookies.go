package cookies

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type SessionCookie struct {
	Username   string `json:"Username"`
	Session_id string `json:"Session_id"`
	IsLogged   bool   `json:"IsLogged"`
}

func WriteCookie(name string, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	return &cookie
}

func ReadCookie(context *echo.Context, name string) (*http.Cookie, error) {
	cookie, err := context.Cookie(name)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func EncodeBase64(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func DecodeBase64(value string) ([]byte, error) {
	res, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return nil, err
	}

	return res, nil
}
