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
