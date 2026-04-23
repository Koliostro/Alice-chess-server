package database_models

import "time"

type GameState string

const (
	ENDED    GameState = "ENDED"
	WAITING  GameState = "WAITING"
	PROGRESS GameState = "PROGRESS"
)

type PLAYERS struct {
	Nick          string `gorm:"type:varchar(32);primaryKey;not null"`
	Passw         string `gorm:"type:varchar(256);not null"`
	Amount_played uint   `gorm:"default:0;not null"`
	Amount_won    uint   `gorm:"default:0;not null"`
	Session_id    string `gorm:"type:varchar(128)"`
	IsAdmin       bool   `gorm:"default:false;not null"`
}

type GAMES struct {
	ID              string    `gorm:"type:varchar(68);primaryKey;not null"`
	Game_date       time.Time `gorm:"type:date;not null"`
	State           GameState `gorm:"type:enum('ENDED', 'PROGRESS', 'WAITING'); not null"`
	White_nick      string    `gorm:"type:varchar(32), not null"`
	Black_nick      string    `gorm:"type:varchar(32), not null"`
	Winner          string    `gorm:"type:varchar(32), not null"`
	Last_turn_left  string    `gorm:"type:varchar(75), not null, default:''"`
	Last_turn_right string    `gorm:"type:varchar(75), not null, default:''"`
}
