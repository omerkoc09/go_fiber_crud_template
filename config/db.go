package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	mysqlerr "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)

	var err error

	level := os.Getenv("DB_LOG_LEVEL")
	var logLevel logger.LogLevel

	switch strings.ToLower(level) {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Silent
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		mysqlErr, ok := err.(*mysqlerr.MySQLError)
		if ok && mysqlErr.Number == 1049 {
			log.Printf("Error connecting to database: Database does not exist")

			if err := createDatabase(dbUser, dbPassword, dbHost, dbPort, dbName); err != nil {
				log.Fatalf("Error creating database: %v", err)
			}

			DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logLevel),
			})

			if err != nil {
				log.Fatalf("Error connecting to new database: %v", err)
			}

			log.Println("Database created successfully")

			if err := runSqlScript(DB); err != nil {
				log.Fatalf("Error running SQL script: %v", err)
			}

			log.Println("SQL script executed successfully")

		} else {
			log.Fatalf("critical error connecting to database: %v", err)
		}

	} else {
		log.Println("successfully connected to database")
	}

}

func createDatabase(dbUser, dbPassword, dbHost, dbPort, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	query := fmt.Sprintf(
		"CREATE DATABASE `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;",
		dbName,
	)

	if err := db.Exec(query).Error; err != nil {
		return err
	}

	sqlDB, _ := db.DB()
	sqlDB.Close()

	return nil
}

func runSqlScript(db *gorm.DB) error {
	filepath := "database/script.sql"

	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("error reading SQL script file: %v", err)
	}

	parts := strings.Split(string(content), "--||--")

	for i, query := range parts {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		if strings.HasPrefix(strings.ToUpper(query), "DELIMITER") ||
			strings.HasPrefix(query, "/*!") ||
			strings.HasPrefix(strings.ToUpper(query), "USE ") {
			continue
		}

		log.Printf("Executing SQL query %d: %s", i+1, query)

		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("error executing SQL query %d: %v", i+1, err)
		}
	}

	return nil
}
