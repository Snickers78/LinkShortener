package main

import (
	"GoAdvanced/configs"
	"GoAdvanced/internal/auth"
	"GoAdvanced/internal/link"
	"GoAdvanced/internal/stat"
	"GoAdvanced/internal/user"
	"GoAdvanced/pkg/db"
	"GoAdvanced/pkg/event"
	"GoAdvanced/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	DB := db.NewDb(conf)
	router := http.NewServeMux()
	eventbus := event.NewEventBus()

	//Repositories
	linkRepository := link.NewLinkRepository(DB)
	userRepository := user.NewUserRepository(DB)
	statRepository := stat.NewStatRepository(DB)

	//Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventbus,
		StatRepository: statRepository,
	})

	// handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventbus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
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

	go statService.AddClick()

	server.ListenAndServe()
}
