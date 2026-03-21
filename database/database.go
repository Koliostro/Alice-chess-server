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

const dsn string = "ChessServer:1212@tcp(127.0.0.1:3306)/CHESS?parseTime=true"

func createSQLERrorHandler(res *gorm.DB) error {
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

func Create_user(db *gorm.DB, player *models.PLAYERS) error {
	res := db.Create(player)

	if res.Error != nil {
		return createSQLERrorHandler(res)
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

func UpdSessionId(db *gorm.DB, player *models.PLAYERS, value string) error {
	res := db.Model(player).Update("session_id", value)

	if res.Error != nil {
		return res.Error
	}

	log.Println("Updated  user's session_id")

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
