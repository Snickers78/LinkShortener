package auth_test

import (
	"GoAdvanced/configs"
	"GoAdvanced/internal/auth"
	"GoAdvanced/internal/user"
	"GoAdvanced/pkg/db"
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bootstrap() (Handler *auth.AuthHandler, mock sqlmock.Sqlmock, err error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestLoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow("Test@mail.ru", "$2a$10$VFDHfdtBCJhdY7kQeL0acOIRis/ogOHsMlGlC2FPdUETZokk8xKSS")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, err := json.Marshal(&auth.LoginRequest{
		Email:    "Test@mail.ru",
		Password: "wrer",
	})
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Login failed. Expected %d. Got %d", http.StatusOK, w.Code)
	}
}
