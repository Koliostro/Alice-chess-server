package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/database"
	"AliceChessServer/database/database_models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

type ConnectMenuData struct {
	Title []string
	Path  []string
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: checkWSOrigin,
	}
)

func checkWSOrigin(request *http.Request) bool {
	Cookie := cookies.NewSessionCookie()

	httpCookie, err := request.Cookie("session_id")

	if err != nil {
		return false
	}

	res, err := Cookie.DecodeCookie(httpCookie)

	if err != nil {
		return false
	}

	err = json.Unmarshal(*res, &Cookie)

	if err != nil {
		return false
	}

	if Cookie.IsLogged {
		return true
	}

	return false
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
		ID:         self.generateSessionId(Cookie.Username),
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
	res, err := database.GetGameById(self.DB, context.Param("id"))

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.Render(http.StatusOK, "game.html", map[string]string{
		"URL": res.ID,
	})
}

func (self *Handler) GetConnectionMenu(context *echo.Context) error {
	var DBData []database_models.GAMES
	var Data ConnectMenuData

	var bufferLink []string = make([]string, 10)
	var bufferTitle []string = make([]string, 10)

	DBData = *database.GetSelectedGames(self.DB, database_models.WAITING)

	counter := 0
	for i := 0; i < len(DBData); i++ {
		if i == 2 {
			break
		}
		bufferLink[i] = DBData[i].ID
		bufferTitle[i] = DBData[i].White_nick
		counter = i
	}

	if len(DBData) > 0 {
		Data.Title = bufferTitle[:counter+1]
		Data.Path = bufferLink[:counter+1]
	} else {
		Data.Title = nil
		Data.Path = nil
	}

	return context.Render(http.StatusOK, "connectionMenu.html", Data)
}

func (self *Handler) PostCloseGame(context *echo.Context) error {
	id := context.Param("id")

	clearedId := strings.TrimRight(id, "/")

	res := database.DeleteGame(self.DB, clearedId)

	if res != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return nil
}

func (self *Handler) WSConnection(context *echo.Context) error {
	type MESSAGE struct {
		Header string `json:"header"`
		Data   string `json:"data"`
	}

	ws, err := upgrader.Upgrade(context.Response(), context.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	Working := true

	var readMessage MESSAGE

	for Working {
		_, data, err := ws.ReadMessage()

		if err != nil {
			return err
		}

		log.Println(string(data))

		err = json.Unmarshal(data, &readMessage)

		if err != nil {
			log.Println("Json decode error: " + err.Error())
		}

		switch readMessage.Header {
		case "END":
			ws.Close()
			Working = false
		case "START":
			newMessage := MESSAGE{
				Header: "RECIVED",
				Data:   "",
			}

			res, err := json.Marshal(&newMessage)

			if err != nil {
				log.Println("JSON encode error" + err.Error())
			}

			ws.WriteMessage(websocket.TextMessage, res)
		}
	}

	return nil
}
