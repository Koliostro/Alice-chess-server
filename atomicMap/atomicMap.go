package atomicMap

import (
	"errors"
	"sync"
)

type BoardState struct {
	IsWhiteTurn bool   `json:"IsWhiteTurn"`
	Left        string `json:"Left"`
	Right       string `json:"Right"`
}

type GameState struct {
	GameID      string
	WhitePlayer string
	BlackPLayer string
	State       *SmallState
}

type SmallState struct {
	Board     BoardState
	IsUpdated bool
	IsReaded  bool
	mutex     sync.Mutex
}

type Store struct {
	Games map[string]*GameState
}

func NewStore() *Store {
	return &Store{Games: make(map[string]*GameState)}
}

func (self *Store) Get(id string, res *GameState) error {
	item := self.Games[id]

	if item == nil {
		return errors.New("Item is not exist")
	}

	res.BlackPLayer = item.BlackPLayer
	res.GameID = item.GameID
	res.State = item.State
	res.WhitePlayer = item.WhitePlayer

	return nil
}

func (self *SmallState) Write(newState *BoardState) {
	self.mutex.Lock()
	self.IsReaded = false
	self.IsUpdated = true
	self.Board.IsWhiteTurn = newState.IsWhiteTurn
	self.Board = *newState
	self.mutex.Unlock()
}

func (self *SmallState) Read(result *SmallState) {
	self.mutex.Lock()
	result.Board.Left = self.Board.Left
	result.Board.Right = self.Board.Right
	result.Board.IsWhiteTurn = self.Board.IsWhiteTurn
	result.IsReaded = self.IsReaded
	result.IsUpdated = self.IsUpdated

	// Save that already readed.
	self.IsReaded = true
	self.IsUpdated = false
	self.mutex.Unlock()
}
