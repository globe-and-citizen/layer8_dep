// Depreciated code and to be removed in future
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	ServerPORT int
}

func dep_SetupDatabaseConnection() *gorm.DB {
	errEnv := godotenv.Load()
	if errEnv != nil {
		logrus.Error("loading env vars", errEnv.Error())
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("SSL_MODE")
	sslRootCert := os.Getenv("SSL_ROOT_CERT")

	// DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s sslrootcert=%s", dbHost, dbUser, dbPass, dbName, dbPort,
		sslMode, sslRootCert)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		logrus.Error("error occured:", err.Error())
		panic("Failed to create a connection to database")
	}
	return db
}

func dep_CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		logrus.Error(err.Error())
		panic("Failed to close connection from Database")
	}
	dbSQL.Close()
}
