// Configurations are specified here because this is just an example project.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"globe-and-citizen/layer8/server/models"
)

var (
	SECRET_KEY string = "46b28fb8-ef5e-4515-9b72-b557180043c6"
	DB         *gorm.DB
)

func InitDB() {
	errEnv := godotenv.Load()
	if errEnv != nil && os.Getenv("ENVIRONMENT") != "testing" {
		panic(errEnv)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	sslMode := os.Getenv("SSL_MODE")
	sslRootCert := os.Getenv("SSL_ROOT_CERT")
	var dsn string

	if sslRootCert != "" {
		// DSN for PostgreSQL
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s sslrootcert=%s", dbHost, dbUser, dbPass, dbName, dbPort,
			sslMode, sslRootCert)
	} else {
		// DSN for PostgreSQL
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPass, dbName, dbPort, sslMode)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = db

	// Configure tables:
	var user models.User
	if err := DB.Where("username = ?", "").First(&user).Error; err != nil {
		configureDatabase()
	}
}

func configureDatabase() {
	DB.Exec(`CREATE TABLE users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		-- phone_number VARCHAR(50) NOT NULL,
		-- address VARCHAR(255) NOT NULL,
		-- email_verified BOOLEAN NOT NULL DEFAULT FALSE,
		-- phone_number_verified BOOLEAN NOT NULL DEFAULT FALSE,
		-- location_verified BOOLEAN NOT NULL DEFAULT FALSE,
		-- national_id_verified BOOLEAN NOT NULL DEFAULT FALSE,
		salt VARCHAR(255) NOT NULL DEFAULT 'ThisIsARandomSalt123!@#',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)

	DB.Exec(`CREATE TABLE clients (
		id VARCHAR(36) PRIMARY KEY,
		secret VARCHAR NOT NULL,
		name VARCHAR(255) NOT NULL,
		redirect_uri VARCHAR(255) NOT NULL
	);`)

	DB.Exec(`CREATE TABLE user_metadata (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		key VARCHAR(255) NOT NULL,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
}
