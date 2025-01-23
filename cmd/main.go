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

func App() http.Handler {
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

	go statService.AddClick()

	//Middleware
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8080",
		Handler: app,
	}

	server.ListenAndServe()
}
