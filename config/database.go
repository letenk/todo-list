package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/letenk/todo-list/models/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var counts int64

func openConnection() (*gorm.DB, error) {
	DBUser := os.Getenv("MYSQL_USER")
	DBPassword := os.Getenv("MYSQL_PASSWORD")
	DBHost := os.Getenv("MYSQL_HOST")
	DBPort := os.Getenv("MYSQL_PORT")
	DBName := os.Getenv("MYSQL_DBNAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DBUser, DBPassword, DBHost, DBPort, DBName)

	// Open connection to db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SetupDB() *gorm.DB {
	for {
		conn, err := openConnection()
		if err != nil {
			log.Println("MySQL not yet ready ...")
			// Increments var counts
			counts++
		} else {
			// Auto Migrate
			err = conn.AutoMigrate(&domain.Activity{}, &domain.Todo{})

			if err != nil {
				log.Fatalf("Failed to auto migration %v", err)
			}
			log.Println("Connected to MySQL!")
			return conn
		}

		// If var count is greater that 1-
		if counts > 10 {
			// Print log error from connection
			log.Printf("Database connection error: %s\n", err)
			return nil
		}

		// Print log for waiting two second each trying connection again
		log.Println("Backing off for two seconds ...")
		time.Sleep(2 * time.Second)
		continue
	}
}
