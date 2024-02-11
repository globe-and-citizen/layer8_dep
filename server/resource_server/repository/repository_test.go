package repository

import (
	"globe-and-citizen/layer8/server/resource_server/dto"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRegisterUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test user DTO
	testUser := dto.RegisterUserDTO{
		Email:       "test@gcitizen.com",
		Username:    "test_user",
		FirstName:   "Test",
		LastName:    "User",
		Password:    "TestPass123",
		Country:     "Unknown",
		DisplayName: "user",
	}

	// Call the RegisterUser function
	repo.RegisterUser(testUser)
}

func TestRegisterClient(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test client DTO
	testClient := dto.RegisterClientDTO{
		Name:        "testclient",
		RedirectURI: "https://gcitizen.com/callback",
	}

	// Call the RegisterClient function
	repo.RegisterClient(testClient)
}

func TestGetClientData(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the client repository with the mock database connection
	repo := NewRepository(db)

	// Define a test client name
	testClientName := "testclient"

	// Call the GetClientData function
	repo.GetClientData(testClientName)

}

func TestLoginPreCheckUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test login precheck DTO
	testLoginPrecheck := dto.LoginPrecheckDTO{
		Username: "test_user",
	}

	repo.LoginPreCheckUser(testLoginPrecheck)

}

func TestProfileUser(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Call the ProfileUser function
	repo.ProfileUser(1)

}

func TestVerifyEmail(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Call the VerifyEmail function
	repo.VerifyEmail(1)

}

func TestUpdateDisplayName(t *testing.T) {
	// Create a new mock DB and a GORM database connection
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("Failed to create mock DB:", err)
	}
	defer mockDB.Close()

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: mockDB}), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to mock DB:", err)
	}

	// Create the user repository with the mock database connection
	repo := NewRepository(db)

	// Define a test update display name DTO
	testUpdateDisplayName := dto.UpdateDisplayNameDTO{
		DisplayName: "new_user",
	}

	// Call the UpdateDisplayName function
	repo.UpdateDisplayName(1, testUpdateDisplayName)

}
