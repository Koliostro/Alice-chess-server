package database_models

type PLAYERS struct {
	Nick          string `gorm:"type:varchar(32);primaryKey;not null"`
	Passw         string `gorm:"type:varchar(256);not null"`
	Amount_played uint   `gorm:"default:0;not null"`
	Amount_won    uint   `gorm:"default:0;not null"`
	Session_id    string `gorm:"type:varchar(128)"`
	IsAdmin       bool   `gorm:"default:false;not null"`
}
