package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMySQL() (*gorm.DB, error) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbname := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, dbname)

	var db *gorm.DB
	var err error

	// retry 10 times, waiting for MySQL
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Connected to MySQL")
			return db, nil
		}
		log.Printf("MySQL not ready (attempt %d/10): %v", i+1, err)
		time.Sleep(3 * time.Second)
	}
	return nil, fmt.Errorf("failed to connect to MySQL after retries: %w", err)
}
