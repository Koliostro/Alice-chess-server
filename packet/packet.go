package packet

type BoardState struct {
	IsYourTurn    bool
	IsGameStarted bool
	Left          string
	Right         string
}

func (self *BoardState) SetWhichTurn() bool {
	if string(self.Left[len(self.Left)-1]) == "w" {
		return true
	} else {
		return false
	}
}
