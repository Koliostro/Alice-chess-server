package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/database"
	"AliceChessServer/database/database_errors"
	"AliceChessServer/database/database_models"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

type SessionCookie struct {
	Username   string `json:"Username"`
	Session_id string `json:"Session_id"`
	IsLogged   bool   `json:"IsLogged"`
}

func NewGenericHandler() (*Handler, error) {
	db, err := database.Db_init()

	if err != nil {
		return nil, err
	}

	return &Handler{DB: db}, nil
}

func (self *Handler) GetReg(context *echo.Context) error {
	if self.checkIfAuthorised(context) {
		return context.Redirect(http.StatusFound, "/")
	}

	return context.Render(http.StatusOK, "register.html", map[string]bool{
		"log":  true,
		"pass": true,
	})
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

	RawCookie, err := cookies.ReadCookie(context, "session_id")

	if err != nil {
		return context.Render(http.StatusOK, "main.html", map[string]string{
			"username": "Not logged",
		})
	}

	ByteCookie, err := cookies.DecodeBase64(RawCookie.Value)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	var ReadSessionCookie SessionCookie

	err = json.Unmarshal(ByteCookie, &ReadSessionCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.Render(http.StatusOK, "main.html", map[string]string{
		"username": ReadSessionCookie.Username,
	})
}

func (self *Handler) generateSessionId(username string) string {
	raw_session_id := username + strconv.Itoa(int(time.Now().Unix()))
	bytes_session_id := sha256.Sum256([]byte(raw_session_id))
	return hex.EncodeToString(bytes_session_id[:])
}

func (self *Handler) checkIfAuthorised(context *echo.Context) bool {
	res, err := cookies.ReadCookie(context, "session_id")

	if err != nil {
		return false
	}
	decodedCookie, err := cookies.DecodeBase64(res.Value)

	if err != nil {
		return false
	}

	var SessionCookie SessionCookie
	err = json.Unmarshal([]byte(decodedCookie), &SessionCookie)

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

func (self *Handler) PostReg(context *echo.Context) error {
	var user database_models.PLAYERS

	username := context.FormValue("login")
	password := context.FormValue("passw")
	password_rep := context.FormValue("passw_rep")

	if password != password_rep {
		return context.Render(http.StatusOK, "register.html", map[string]bool{
			"log":  true,
			"pass": false,
		})
	}
	user.Nick = username
	user.Passw = password

	// Generate session id
	user.Session_id = self.generateSessionId(user.Nick)

	Cookie := SessionCookie{
		Session_id: user.Session_id,
		Username:   user.Nick,
		IsLogged:   true,
	}

	byteCookie, err := json.Marshal(Cookie)

	if err != nil {
		log.Println("JSON error: " + err.Error())
		return context.NoContent(http.StatusInternalServerError)
	}

	cookie := cookies.WriteCookie("session_id", cookies.EncodeBase64(byteCookie))

	err = database.Create_user(self.DB, &user)

	if err != nil {
		if errors.Is(err, database_errors.SQLErrObjDup) {
			return context.Render(http.StatusOK, "register.html", map[string]bool{
				"log":  false,
				"pass": true,
			})
		} else {
			return context.NoContent(http.StatusInternalServerError)
		}
	}

	context.SetCookie(cookie)

	return context.Redirect(http.StatusFound, "/")
}

func (self *Handler) NotImplemented(context *echo.Context) error {
	return context.NoContent(http.StatusNotImplemented)
}
