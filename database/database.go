package database

import (
	"AliceChessServer/database/database_errors"
	models "AliceChessServer/database/database_models"
	"log"

	"strconv"

	offSql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dsn string = "ChessServer:1212@tcp(127.0.0.1:3306)/CHESS"

func Create_user(db *gorm.DB, player *models.PLAYERS) error {
	res := db.Create(player)

	if res.Error != nil {
		var mysqlErr *offSql.MySQLError
		mysqlErr = res.Error.(*offSql.MySQLError)

		switch mysqlErr.Number {
		case 1062:
			log.Println(uint16(mysqlErr.Number))
			return database_errors.SQLErrObjDup
		default:
			log.Println(strconv.FormatUint(uint64(mysqlErr.Number), 10))
			return database_errors.SQLErrUnexp
		}
	}

	log.Println("Created user")
	return nil
}

func Find_user(db *gorm.DB, username string) *models.PLAYERS {
	PLAYER := models.PLAYERS{
		Nick: username,
	}

	db.First(&PLAYER, "Nick = ?", username)

	return &PLAYER
}

func Upd_usr_passw(db *gorm.DB, player *models.PLAYERS, value string) error {
	res := db.Model(player).Update("Passw", value)

	if res.Error != nil {
		return res.Error
	}

	log.Println("Updated  user")

	return nil
}

func Db_init() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("Initialized database")

	return db, nil
}
