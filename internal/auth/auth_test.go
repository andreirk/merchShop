package auth_test

import (
	"bytes"
	"encoding/json"
	"go/avito-test/config"
	"go/avito-test/internal/auth"
	"go/avito-test/internal/models"
	"go/avito-test/internal/repositories"
	"go/avito-test/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) Create(u *models.User) (*models.User, error) {
	return &models.User{
		Username: "bob",
	}, nil
}
func (repo *MockUserRepository) FindUserByName(email string) (*models.User, error) {
	return nil, nil
}

func bootstrap() (*auth.Handler, sqlmock.Sqlmock, error) {
	dataBase, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: dataBase,
	}))
	if err != nil {

		return nil, nil, err
	}
	userRepo := repositories.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.Handler{
		Config: &config.Config{
			Auth: config.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil

}

func TestAuthHandlerFailed(t *testing.T) {
	const testUserName = "bob"
	correctPassHash := "$2a$10$J5YzQCdVbIP1i1purfjYEOEEtlsr4UgX7VCOSkDkWoJIldoJ69Vaa"
	wrongPassword := "123456"
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow(testUserName, correctPassHash)
	mock.ExpectQuery("Select").WillReturnRows(rows)

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub Gorm database connection", err)
		return
	}
	data, err := json.Marshal(&auth.LoginRequestDTO{
		Username: testUserName,
		Password: wrongPassword,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Code == http.StatusOK {
		t.Errorf("got %d, expected %d ", w.Code, http.StatusOK)
	}

}

func TestAuthHandlerSuccess(t *testing.T) {
	const testUserName = "bob"
	password := "123456"
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"password", "name"})
	mock.ExpectQuery("Select").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub Gorm database connection", err)
		return
	}
	data, err := json.Marshal(&auth.RegisterRequest{
		Username: testUserName,
		Password: password,
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/auth", reader)
	handler.Register()(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("got %d, expected %d ", w.Code, http.StatusCreated)
	}
}
