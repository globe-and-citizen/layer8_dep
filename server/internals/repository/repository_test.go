package repository_test

import (
	"globe-and-citizen/layer8/server/internals/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestLoginUserPrecheck(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := repository.NewPostgresRepository(db)

	// Make a mock loginPrecheck input
	username := "test_user"

	mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE username = (.+)").WithArgs(username).WillReturnRows(sqlmock.NewRows([]string{"salt"}).AddRow("test_salt"))

	// Call the function and check the result
	salt, err := repo.LoginUserPrecheck(username)
	if err != nil {
		t.Fatal("Failed to call LoginUserPrecheck:", err)
	}

	if salt != "test_salt" {
		t.Fatal("LoginUserPrecheck returned an unexpected salt:", salt)
	}

	// Check if all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal("Unmet expectations:", err)
	}
}
