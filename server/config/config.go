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

	configureDatabase()

	var user models.User
	if err := DB.Where("username = ?", "admin").First(&user).Error; err != nil {
		configureTables()
		DB.Where("username = ?", "admin").First(&user)
		fmt.Println("Test user created.")
	}
	fmt.Println("Use username: admin, password: 12345 for testing.")
}

func configureDatabase() {
	DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		username VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		salt VARCHAR(255) NOT NULL DEFAULT 'ThisIsARandomSalt123!@#',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)

	DB.Exec(`CREATE TABLE IF NOT EXISTS clients (
		id VARCHAR(36) PRIMARY KEY,
		secret VARCHAR NOT NULL,
		name VARCHAR(255) NOT NULL,
		redirect_uri VARCHAR(255) NOT NULL
	);`)

	DB.Exec(`CREATE TABLE IF NOT EXISTS user_metadata (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		key VARCHAR(255) NOT NULL,
		value VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`)
}

func configureTables() {
	DB.Exec(`INSERT INTO users (email, username, password, first_name, last_name) VALUES ('admin@gcitizen.com', 'admin', '12345', 'Admin', 'User');`)

	DB.Exec(`INSERT INTO user_metadata (user_id, key, value) VALUES (1, 'email_verified', 'true');`)
	DB.Exec(`INSERT INTO user_metadata (user_id, key, value) VALUES (1, 'country', 'Globe');`)
	DB.Exec(`INSERT INTO user_metadata (user_id, key, value) VALUES (1, 'display_name', 'Admin User');`)

	DB.Exec(`INSERT INTO clients (id, secret, name, redirect_uri) VALUES ('notanid', 'absolutelynotasecret!', 'Ex-C', 'http://localhost:5173/oauth2/callback');`)
}
