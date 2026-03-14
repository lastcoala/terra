package app

import (
	"github.com/lastcoala/terra/config"
	"github.com/lastcoala/terra/internal/app/repo"
	"github.com/lastcoala/terra/internal/app/rest"
	"github.com/lastcoala/terra/internal/app/service"
	"github.com/lastcoala/terra/pkg/handler"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) CreateRestServer() handler.IHandler {
	gorm, err := repo.NewGormRepo(a.cfg.Rest.Db.DataStore, a.cfg.Rest.Db.NumberConn)
	if err != nil {
		panic(err)
	}

	userRepo := repo.NewUserGormRepo(gorm)
	userService := service.NewUserService(userRepo)

	restHandler := rest.NewRest(a.cfg.Rest.Server.Host, userService)

	return restHandler
}
