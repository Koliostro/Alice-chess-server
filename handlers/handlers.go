package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/cookies/sessionCookie"
	"AliceChessServer/database"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewGenericHandler() (*Handler, error) {
	db, err := database.Db_init()

	if err != nil {
		return nil, err
	}

	return &Handler{DB: db}, nil
}

func (self *Handler) GetGame(context *echo.Context) error {
	return context.File("./templates/game.html")
}

func (self *Handler) GetMain(context *echo.Context) error {
	if !self.checkIfAuthorised(context) {
		return context.Render(http.StatusOK, "main.html", map[string]string{
			"username": "Not logged",
		})
	}

	SesCookie := cookies.NewSessionCookie()

	RawCookie, err := SesCookie.ReadCookie(context, "session_id")

	if err != nil {
		return context.Render(http.StatusOK, "main.html", map[string]string{
			"username": "Not logged",
		})
	}

	cookie, err := SesCookie.DecodeCookie(RawCookie)

	var ReadCookie sessionCookie.SessionCookie

	err = json.Unmarshal(*cookie, &ReadCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.Render(http.StatusOK, "main.html", map[string]string{
		"username": ReadCookie.Username,
	})
}

func (self *Handler) generateSessionId(username string) string {
	raw_session_id := username + strconv.Itoa(int(time.Now().Unix()))
	bytes_session_id := sha256.Sum256([]byte(raw_session_id))
	return hex.EncodeToString(bytes_session_id[:])
}

func (self *Handler) checkIfAuthorised(context *echo.Context) bool {
	Cookie := cookies.NewSessionCookie()

	rawCookie, err := Cookie.ReadCookie(context, sessionCookie.CookieName)

	if err != nil {
		return false
	}

	var SessionCookie sessionCookie.SessionCookie

	byteCookie, err := Cookie.DecodeCookie(rawCookie)

	if err != nil {
		return false
	}

	err = json.Unmarshal(*byteCookie, &SessionCookie)

	if err != nil {
		return false
	}

	user := database.Find_user(self.DB, SessionCookie.Username)

	if user.Session_id != SessionCookie.Session_id {
		return false
	}

	if SessionCookie.IsLogged {
		return true
	}

	return false
}

func (self *Handler) NotImplemented(context *echo.Context) error {
	return context.NoContent(http.StatusNotImplemented)
}
