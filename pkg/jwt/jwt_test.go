package jwt_test

import (
	"GoAdvanced/pkg/jwt"
	"testing"
)

func TestJWTCreate(t *testing.T) {
	const email = "Test@mail.ru"
	JWT := jwt.NewJWT("+XZFh214gfd-1hx")
	token, err := JWT.Create(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	IsValid, data := JWT.Parse(token)
	if !IsValid {
		t.Fatalf("Email %s is not equal to %s", data, email)
	}
}
