package database

import (
	"AliceChessServer/database/database_models"
	models "AliceChessServer/database/database_models"

	"gorm.io/gorm"
)

func CreateGame(db *gorm.DB, game *models.GAMES) error {
	res := db.Create(game)

	if res.Error != nil {
		return createSQLERrorHandler(res)
	}

	return nil
}

func GetSelectedGames(db *gorm.DB, state models.GameState) *[]database_models.GAMES {
	var GAME []database_models.GAMES

	db.Find(&GAME, "State = ?", state).Limit(10)

	return &GAME
}

func GetGameById(db *gorm.DB, id string) (*models.GAMES, error) {
	GAME := models.GAMES{
		ID: id,
	}
	res := db.First(&GAME)

	if res.Error != nil {
		return nil, res.Error
	}

	return &GAME, nil
}

func UpdateGameState(db *gorm.DB, game *models.GAMES, state string) error {
	res := db.Model(game).Update("State", state)
	return res.Error
}

func EndGame(db *gorm.DB, game *models.GAMES, WinnerStr string, gameID string) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		record := db.Model(game).Where("ID = ?", gameID)
		res := record.Update("State", models.ENDED)
		if res.Error != nil {
			return res.Error
		}

		res = record.Update("Winner", WinnerStr)
		if res.Error != nil {
			return res.Error
		}
		return nil
	})

	return err
}

func UpdateGamePlayer(db *gorm.DB, game *models.GAMES, nick string) error {
	res := db.Model(game).Update("Black_nick", nick)
	return res.Error
}

func UpdateGameTurns(db *gorm.DB, game *models.GAMES, Left string, Right string) error {
	res := db.Model(game).Updates(models.GAMES{Last_turn_left: Left, Last_turn_right: Right})
	return res.Error
}

func DeleteGame(db *gorm.DB, id string) error {
	GAME := models.GAMES{
		ID: id,
	}

	res := db.Delete(&GAME)

	return res.Error
}
