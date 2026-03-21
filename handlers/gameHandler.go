package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/database"
	"AliceChessServer/database/database_models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

type ConnectMenuData struct {
	Title []string
}

func (self *Handler) GetCreateRoom(context *echo.Context) error {
	Cookie := cookies.NewSessionCookie()

	rawCookie, err := Cookie.ReadCookie(context)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	byteCookie, err := Cookie.DecodeCookie(rawCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	err = json.Unmarshal(*byteCookie, &Cookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	gameObj := database_models.GAMES{
		ID:         Cookie.Username,
		Game_date:  time.Now(),
		State:      database_models.WAITING,
		White_nick: Cookie.Username,
		Black_nick: "",
		Winner:     "",
	}

	err = database.CreateGame(self.DB, &gameObj)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.Redirect(http.StatusFound, "/games/waiting/"+gameObj.ID)
}

func (self *Handler) GetwaitingRoom(context *echo.Context) error {
	return context.Render(http.StatusOK, "waitingRoom.html", map[string]string{
		"id": context.Param("id"),
	})
}

func (self *Handler) GetConnectionMenu(context *echo.Context) error {
	var DBData []database_models.GAMES
	var Data ConnectMenuData

	var buffer []string = make([]string, 10)

	DBData = *database.GetSelectedGames(self.DB, database_models.WAITING)

	counter := 0
	for i := 0; i < len(DBData); i++ {
		if i == 2 {
			break
		}
		buffer[i] = DBData[i].ID
		counter = i
	}

	Data.Title = buffer[:counter+1]

	return context.Render(http.StatusOK, "connectionMenu.html", Data)
}

func (self *Handler) PostCloseGame(context *echo.Context) error {
	return context.String(http.StatusOK, context.Param("id"))
}
