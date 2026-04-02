package atomicMap

type BoardState struct {
	Left  string `json:"Left"`
	Right string `json:"Right"`
}

type GameState struct {
	GameID      string
	WhitePlayer string
	BlackPLayer string
	State       SmallState
}

type SmallState struct {
	Board     BoardState
	IsUpdated bool
	IsReaded  bool
}

type Store struct {
	Games map[string]*GameState
}

func NewStore() *Store {
	return &Store{Games: make(map[string]*GameState)}
}
