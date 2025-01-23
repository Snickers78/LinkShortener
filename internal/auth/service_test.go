package auth_test

import (
	"GoAdvanced/internal/auth"
	"GoAdvanced/internal/user"
	"testing"
)

type MockUserRepository struct{}

func (repo *MockUserRepository) FindByEmail(email string) (*user.User, error) {
	return nil, nil
}

func (repo *MockUserRepository) Create(u *user.User) (*user.User, error) {
	return &user.User{
		Email: "Test@mail.ru",
	}, nil
}

func TestRegisterSuccess(t *testing.T) {
	const InitialEmail = "Test@mail.ru"
	AuthService := auth.NewAuthService(&MockUserRepository{})
	email, err := AuthService.Register("Test@mail.ru", "wrer", "Petya")
	if err != nil {
		t.Fatal(err)
	}
	if email != InitialEmail {
		t.Fatalf("Email does not match. Expected %s, got %s", InitialEmail, email)
	}
}
