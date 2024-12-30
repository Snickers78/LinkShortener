package main

import (
	"GoAdvanced/configs"
	"GoAdvanced/internal/auth"
	"GoAdvanced/internal/link"
	"GoAdvanced/internal/user"
	"GoAdvanced/pkg/db"
	"GoAdvanced/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	DB := db.NewDb(conf)
	router := http.NewServeMux()

	//Repositories
	linkRepository := link.NewLinkRepository(*DB)
	userRepository := user.NewUserRepository(*DB)

	//Services
	authService := auth.NewAuthService(userRepository)

	// handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config:         conf,
	})

	//Middleware
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}
	server.ListenAndServe()
}
