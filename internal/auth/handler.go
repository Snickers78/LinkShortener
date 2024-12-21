package auth

import (
	"GoAdvanced/configs"
	"GoAdvanced/pkg/req"
	"GoAdvanced/pkg/res"
	"fmt"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}
type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, Deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      Deps.Config,
		AuthService: Deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		fmt.Println(body)
		resp := LoginResponce{
			Token: "123",
		}
		res.Json(w, http.StatusOK, resp)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRegister](&w, r)
		if err != nil {
			return
		}
		handler.AuthService.Register(body.Email, body.Password, body.Name)
	}

}
