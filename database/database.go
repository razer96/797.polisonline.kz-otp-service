package database

import (
	"os"

	"insurance-otp-service/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Database struct {
	DB         *gorm.DB
	dbHost     string
	dbUsername string
	dbPassword string
	dbPort     string
	dbName     string
}

var (
	log = logger.GetLogger()
	db  = &Database{}
)

func init() {
	db.dbHost = os.Getenv("DATABASE_HOST")
	db.dbUsername = os.Getenv("DATABASE_USERNAME")
	db.dbPassword = os.Getenv("DATABASE_PASSWORD")
	db.dbPort = os.Getenv("DATABASE_PORT")
	db.dbName = os.Getenv("DATABASE_NAME")
	connectionString := "host=" + db.dbHost + " port=" + db.dbPort + " user=" + db.dbUsername + " dbname=" + db.dbName + " password=" + db.dbPassword + " sslmode=disable"
	var err error
	db.DB, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Error("Connection to db was faild", err)
		os.Exit(1)
	}
}

func GetDB() *Database {
	return db
}

func (d *Database) Close() {
	d.DB.Close()
}
