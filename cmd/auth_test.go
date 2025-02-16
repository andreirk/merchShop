package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/avito-test/config"
	"go/avito-test/internal/auth"
	"go/avito-test/internal/models"
	"io"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func initDB() *gorm.DB {
	err := godotenv.Load("test.env")
	if err != nil {
		panic("Error loading .env file")
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}
	return db
}

func initData(db *gorm.DB, name string) {
	passHash := "$2a$10$J5YzQCdVbIP1i1purfjYEOEEtlsr4UgX7VCOSkDkWoJIldoJ69Vaa"
	id := db.Create(&models.User{
		Password: passHash,
		Username: name,
	})
	fmt.Println("created id", id)
}
func cleanData(db *gorm.DB, name any) {
	db.Unscoped().
		Where("username = ?", name).
		Delete(&models.User{})
}

func TestAuthSuccess(t *testing.T) {
	// Prepare
	db := initDB()
	testName := "VasiaSuccessTest"
	initData(db, testName)
	defer cleanData(db, testName)
	app := NewApp(config.LoadConfig("test"))
	ts := httptest.NewServer(app.Handler)
	defer ts.Close()
	data, _ := json.Marshal(&auth.AuthRequestDTO{
		Username: testName,
		Password: "123456",
	})

	res, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != 200 {
		t.Fatalf("Expected %d got %d", 200, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.AuthResponseDTO
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.AccessToken == "" {
		t.Fatal("access token empty")
	}
}

func TestAuthFailed(t *testing.T) {
	// Prepare
	db := initDB()
	testUserName := "VasiaFailedTest"
	initData(db, testUserName)
	defer cleanData(db, testUserName)

	// Test
	ts := httptest.NewServer(NewApp(config.LoadConfig("test")).Handler)
	defer ts.Close()
	data, _ := json.Marshal(&auth.LoginRequestDTO{
		Username: "Bob",
		// no password ! should fail
	})

	res, err := http.Post(ts.URL+"/api/auth", "application/json", bytes.NewReader(data))

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Printing Body", res.Body)
	if res.StatusCode != 400 {
		t.Fatalf("Expected %d got %d", 401, res.StatusCode)
	}
}
