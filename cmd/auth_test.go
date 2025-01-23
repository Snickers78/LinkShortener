package main

import (
	"GoAdvanced/internal/auth"
	"GoAdvanced/internal/user"
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func InitDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func InitData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    "Test@mail.ru",
		Password: "$2a$10$VFDHfdtBCJhdY7kQeL0acOIRis/ogOHsMlGlC2FPdUETZokk8xKSS",
		Name:     "Petya",
	})
}

func deleteData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", "Test@mail.ru").
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	db := InitDb()
	InitData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, err := json.Marshal(&auth.LoginRequest{
		Email:    "Test@mail.ru",
		Password: "wrer",
	})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var loginResp auth.LoginResponce
	err = json.Unmarshal(body, &loginResp)
	if err != nil {
		t.Fatal(err)
	}
	if loginResp.Token == "" {
		t.Fatalf("Received empty response\n")
	}
	deleteData(db)
}
