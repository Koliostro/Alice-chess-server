package handlers

import (
	"AliceChessServer/cookies"
	"AliceChessServer/cookies/sessionCookie"
	"AliceChessServer/database"
	"AliceChessServer/database/database_models"
	"AliceChessServer/packet"
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
		return context.Redirect(http.StatusFound, "/login")
	}

	byteCookie, err := Cookie.DecodeCookie(rawCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	err = json.Unmarshal(*byteCookie, &Cookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	newGame := database_models.GAMES{
		ID:              self.generateSessionId(Cookie.Username),
		Game_date:       time.Now(),
		State:           database_models.WAITING,
		White_nick:      Cookie.Username,
		Black_nick:      "",
		Winner:          "",
		Last_turn_right: "R7/8/8/8/8/8/8/8 w",
		Last_turn_left:  "8/8/8/8/8/8/8/8 w",
	}

	err = database.CreateGame(self.DB, &newGame)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	return context.Redirect(http.StatusFound, "/games/"+newGame.ID)
}

func (self *Handler) GetwaitingRoom(context *echo.Context) error {
	var Cookie sessionCookie.SessionCookie
	Cookie = *cookies.NewSessionCookie()

	httpCookie, err := Cookie.ReadCookie(context)

	byteCookie, err := Cookie.DecodeCookie(httpCookie)

	json.Unmarshal(*byteCookie, &Cookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	username := Cookie.Username

	game, err := database.GetGameById(self.DB, context.Param("id"))

	// Only once set second player
	if game.White_nick != username && game.Black_nick == "" {
		game.Black_nick = username

		model, err := database.GetGameById(self.DB, context.Param("id"))

		if err != nil {
			return context.NoContent(http.StatusInternalServerError)
		}

		database.UpdateGamePlayer(self.DB, model, username)
		database.UpdateGameState(self.DB, model, "PROGRESS")

		log.Printf("Black player are connected\n")
	}

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

// TODO: change
func (self *Handler) PostNewState(context *echo.Context) error {
	var board packet.BoardState
	var sessionCookie sessionCookie.SessionCookie

	sessionCookie = *cookies.NewSessionCookie()

	err := context.Bind(&board)

	httpCookie, err := sessionCookie.ReadCookie(context)
	res, err := sessionCookie.DecodeCookie(httpCookie)
	json.Unmarshal(*res, &sessionCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	log.Printf("Getted packet : %v\n", board)

	// TODO: Need to fix sending of new state.

	ID := context.Param("id")

	game, err := database.GetGameById(self.DB, ID)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	if board.SetWhichTurn() {
		if sessionCookie.Username == game.White_nick {
			board.IsYourTurn = true
		} else {
			board.IsYourTurn = false
		}
	} else {
		if sessionCookie.Username == game.Black_nick {
			board.IsYourTurn = true
		} else {
			board.IsYourTurn = false
		}
	}

	// Do not work.
	err = database.UpdateGameTurns(self.DB, game, board.Left, board.Right)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	newstate, _ := database.GetGameById(self.DB, ID)
	log.Printf("New State : %v", newstate)

	return context.NoContent(http.StatusOK)
}

func (self *Handler) GetGameState(context *echo.Context) error {
	var Cookie sessionCookie.SessionCookie
	Cookie = *cookies.NewSessionCookie()

	httpCookie, err := Cookie.ReadCookie(context)
	byteCookie, err := Cookie.DecodeCookie(httpCookie)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	json.Unmarshal(*byteCookie, &Cookie)

	// NOTE: We presume that id is valid
	ID := context.Param("id")

	game, err := database.GetGameById(self.DB, ID)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	blackPlayer := game.Black_nick
	whitePlayer := game.White_nick

	// 1) Check if game is started
	// 2) If started check who ask and which turn are now

	var GamePacket packet.BoardState

	if whitePlayer == "" || blackPlayer == "" {
		log.Println("White player are waiting")

		GamePacket.IsGameStarted = false
		GamePacket.IsYourTurn = false
		GamePacket.Left = ""
		GamePacket.Right = ""

		res, err := json.Marshal(GamePacket)

		if err != nil {
			return context.NoContent(http.StatusInternalServerError)
		}

		return context.String(http.StatusOK, string(res))
	}

	// Should work. Need to test it.
	GamePacket.IsGameStarted = true
	GamePacket.Left = game.Last_turn_left
	GamePacket.Right = game.Last_turn_right

	if GamePacket.SetWhichTurn() {
		if Cookie.Username == game.White_nick {
			GamePacket.IsYourTurn = true
		} else {
			GamePacket.IsYourTurn = false
		}
	} else {
		if Cookie.Username == game.Black_nick {
			GamePacket.IsYourTurn = true
		} else {
			GamePacket.IsYourTurn = false
		}
	}

	res, err := json.Marshal(GamePacket)

	if err != nil {
		return context.NoContent(http.StatusInternalServerError)
	}

	// NOTE: Return through STRING not JSON!!!!
	return context.String(http.StatusOK, string(res))
}
