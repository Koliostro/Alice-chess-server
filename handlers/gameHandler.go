package handlers

import (
	"AliceChessServer/atomicMap"
	"AliceChessServer/cookies"
	"AliceChessServer/cookies/sessionCookie"
	"AliceChessServer/database"
	"AliceChessServer/database/database_models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

type ConnectMenuData struct {
	Title []string
	Path  []string
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

	boardState := atomicMap.BoardState{
		IsWhiteTurn: true,
		Left:        "BKNQP3/8/8/8/8/8/8/7R",
		Right:       "8/8/8/8/8/8/8/6r1",
	}

	smallState := atomicMap.SmallState{
		Board:     boardState,
		IsReaded:  false,
		IsUpdated: false,
	}

	initialState := atomicMap.GameState{
		GameID:      self.generateSessionId(Cookie.Username),
		WhitePlayer: Cookie.Username,
		BlackPLayer: "",
		State:       &smallState,
	}

	newGame := database_models.GAMES{
		ID:         initialState.GameID,
		Game_date:  time.Now(),
		State:      database_models.WAITING,
		White_nick: initialState.WhitePlayer,
		Black_nick: "",
		Winner:     "",
	}

	err = database.CreateGame(self.DB, &newGame)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	//NOTE: I still doesn't know is my system for multythread writing are working but without testing
	//		I don't get any answers.

	self.Hub.Games[initialState.GameID] = &initialState

	return context.Redirect(http.StatusFound, "/games/"+initialState.GameID)
}

func (self *Handler) GetwaitingRoom(context *echo.Context) error {
	return context.Render(http.StatusOK, "game.html", map[string]string{
		"URL": context.Param("id"),
	})
}

func (self *Handler) GetConnectionMenu(context *echo.Context) error {
	var DBData []database_models.GAMES
	var Data ConnectMenuData

	var bufferLink []string = make([]string, 10)
	var bufferTitle []string = make([]string, 10)

	DBData = *database.GetSelectedGames(self.DB, database_models.WAITING)

	for i := 0; i < len(DBData); i++ {
		err := self.Hub.Get(DBData[i].ID, &atomicMap.GameState{})

		if err != nil {
			database.DeleteGame(self.DB, DBData[i].ID)
		}
	}

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

func (self *Handler) PostNewState(context *echo.Context) error {
	var board atomicMap.BoardState
	var item atomicMap.GameState
	var sessionCookie sessionCookie.SessionCookie

	sessionCookie = *cookies.NewSessionCookie()

	err := context.Bind(&board)

	httpCookie, err := sessionCookie.ReadCookie(context)
	res, err := sessionCookie.DecodeCookie(httpCookie)
	json.Unmarshal(*res, &sessionCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	id := context.Param("id")

	log.Println(board)

	err = self.Hub.Get(id, &item)

	if err != nil {
		return context.NoContent(http.StatusBadRequest)
	}

	item.State.Write(&board)

	return context.NoContent(http.StatusOK)
}

func (self *Handler) GetGameState(context *echo.Context) error {
	var gettedState atomicMap.SmallState
	var gettedGameState atomicMap.GameState

	// NOTE: We presume that id is valid
	ID := context.Param("id")

	err := self.Hub.Get(ID, &gettedGameState)

	if err != nil {
		return context.String(http.StatusBadRequest, "")
	}

	gettedGameState.State.Read(&gettedState)

	log.Print(gettedState.Board)

	res, err := json.Marshal(gettedState.Board)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	// NOTE: Return through STRING not JSON!!!!
	return context.String(http.StatusOK, string(res))
}
