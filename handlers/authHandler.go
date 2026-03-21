package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/cookies/sessionCookie"
	"AliceChessServer/database"
	"AliceChessServer/database/database_errors"
	"AliceChessServer/database/database_models"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

func (self *Handler) GetReg(context *echo.Context) error {
	return context.Render(http.StatusOK, "register.html", map[string]bool{
		"log":  true,
		"pass": true,
	})
}

func (self *Handler) GetLogin(context *echo.Context) error {
	if self.checkIfAuthorised(context) {
		return context.Render(http.StatusOK, "login.html", map[string]bool{
			"logout": true,
		})
	}
	return context.Render(http.StatusOK, "login.html", nil)
}

func (self *Handler) PostLogin(context *echo.Context) error {
	var player *database_models.PLAYERS
	username := context.FormValue("login")
	password := context.FormValue("passw")

	player = database.Find_user(self.DB, username)

	if player.Passw != password {
		return context.Render(http.StatusOK, "login.html", map[string]bool{
			"Notlog": false,
			"logout": false,
		})
	}

	Cookie := cookies.NewSessionCookie()

	Cookie = &sessionCookie.SessionCookie{
		Username:   player.Nick,
		Session_id: self.generateSessionId(player.Nick),
		IsLogged:   true,
	}

	byteCookie, err := json.Marshal(Cookie)

	if err != nil {
		log.Println(err.Error())
		return context.NoContent(http.StatusInternalServerError)
	}

	player.Session_id = Cookie.Session_id
	err = database.UpdSessionId(self.DB, player, Cookie.Session_id)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	context.SetCookie(Cookie.WriteCookie(&byteCookie))

	return context.Redirect(http.StatusFound, "/")
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

	user.Session_id = self.generateSessionId(user.Nick)

	Cookie := sessionCookie.SessionCookie{
		Session_id: user.Session_id,
		Username:   user.Nick,
		IsLogged:   true,
	}

	byteCookie, err := json.Marshal(Cookie)

	if err != nil {
		log.Println("JSON error: " + err.Error())
		return context.NoContent(http.StatusInternalServerError)
	}

	cookie := Cookie.WriteCookie(&byteCookie)

	err = database.Create_user(self.DB, &user)

	if err != nil {
		if errors.Is(err, database_errors.SQLErrObjDup) {
			return context.Render(http.StatusOK, "register.html", map[string]bool{
				"Notlog": false,
				"pass":   true,
			})
		} else {
			return context.NoContent(http.StatusInternalServerError)
		}
	}

	context.SetCookie(cookie)

	return context.Redirect(http.StatusFound, "/")
}

func (self *Handler) GetLogOut(context *echo.Context) error {
	Cookie := cookies.NewSessionCookie()

	httpCookie, err := Cookie.ReadCookie(context)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	httpCookie.MaxAge = -1
	httpCookie.Path = "/"

	context.SetCookie(httpCookie)

	return context.Redirect(http.StatusFound, "/")
}
