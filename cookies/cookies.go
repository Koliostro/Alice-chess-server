package cookies

import (
	"AliceChessServer/cookies/sessionCookie"
	"net/http"

	"github.com/labstack/echo/v5"
)

type GenericCookie interface {
	WriteCookie(value *[]byte) *http.Cookie
	ReadCookier(context *echo.Context) (*http.Cookie, error)
	DecodeCookie(cookie *http.Cookie) (*[]byte, error)
}

func NewSessionCookie() *sessionCookie.SessionCookie {
	cookie := sessionCookie.SessionCookie{
		Username:   "",
		Session_id: "",
		IsLogged:   false,
	}

	return &cookie
}
