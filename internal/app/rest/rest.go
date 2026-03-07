package rest

import (
	"github.com/labstack/echo/v4"
	v1 "github.com/lastcoala/terra/internal/app/rest/v1"
	"github.com/lastcoala/terra/internal/app/service"
	"github.com/lastcoala/terra/pkg/handler"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title TERRA API
// @version 1.0
// @description This is a sample REST API Server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /v1
type RestHandler struct {
	*handler.ServerHandler
	userService service.IUserService
}

func NewRest(host string, userService service.IUserService) *RestHandler {

	server := handler.NewServerHandler(host)

	return &RestHandler{
		ServerHandler: server,
		userService:   userService,
	}
}

func (h *RestHandler) Start() error {
	h.Echo = echo.New()
	h.Echo.HideBanner = true
	h.Route()
	h.StartServer()

	return nil
}

func (h *RestHandler) Route() {
	h.ServerHandler.Route()

	v1Route := h.Echo.Group("/v1")

	v1Route.GET("/swagger/*", echoSwagger.WrapHandler)
	v1Route.GET("/ping", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	v1.Route(v1Route, h.userService)
}
