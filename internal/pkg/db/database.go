package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/config"

	"time"
)

var (
	Connec *sql.DB
	err    error
)

func GetDB() *sql.DB {
	return Connec
}

type Connection struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// SetupDB opens a database and saves the reference to `Database` struct.
func SetupDB() {
	var db = Connec

	configuration := config.GetConfig()
	connInfo := Connection{
		Host:     configuration.Database.Host,
		Port:     configuration.Database.Port,
		Username: configuration.Database.Username,
		Password: configuration.Database.Password,
		Database: configuration.Database.Dbname,
	}
	db, err = sql.Open(configuration.Database.Driver, connToString(connInfo))
	if err != nil {
		fmt.Printf("Error connecting to the DB: %s\n", err.Error())
		return
	} else {
		fmt.Printf("BD conectada con exito!!\n")
	}

	db.SetMaxIdleConns(configuration.Database.MaxIdleConns)
	db.SetMaxOpenConns(configuration.Database.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(configuration.Database.MaxLifetime) * time.Second)
	Connec = db
}

func connToString(info Connection) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.Host, info.Port, info.Username, info.Password, info.Database)

}
