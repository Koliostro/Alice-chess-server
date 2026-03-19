package handlers

import (
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
			"log": true,
		})
	}
	return context.Render(http.StatusOK, "login.html", nil)
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
