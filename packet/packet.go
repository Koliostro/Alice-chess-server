package packet

type BoardState struct {
	IsYourTurn    bool
	IsGameStarted bool
	Left          string
	Right         string
}

func (self *BoardState) SetWhichTurn() bool {
	for i := 0; i < len(string(self.Left)); i++ {
		if string(self.Left[i]) == " " {
			if string(self.Left[i+1]) == "w" {
				return true
			} else {
				return false
			}
		}
	}
	return false
}
