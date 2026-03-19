package sessionCookie

import (
	base64 "AliceChessServer/encoding"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

const CookieName string = "session_id"

type SessionCookie struct {
	Username   string `json:"Username"`
	Session_id string `json:"Session_id"`
	IsLogged   bool   `json:"IsLogged"`
}

func (self *SessionCookie) WriteCookie(value *[]byte) *http.Cookie {
	httpCookie := http.Cookie{
		Name:     CookieName,
		Value:    base64.Encode(*value),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}

	log.Println("Value: " + httpCookie.Value)

	return &httpCookie
}

func (self *SessionCookie) ReadCookie(context *echo.Context, name string) (*http.Cookie, error) {
	cookie, err := context.Cookie(name)

	if err != nil {
		return nil, err
	}

	return cookie, nil
}

func (self *SessionCookie) DecodeCookie(httoCookie *http.Cookie) (*[]byte, error) {
	raw, err := base64.Decode(httoCookie.Value)

	if err != nil {
		return nil, err
	}
	return &raw, nil
}
