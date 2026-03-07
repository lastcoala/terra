package v1

import (
	echo "github.com/labstack/echo/v4"
	"github.com/lastcoala/terra/internal/app/service"
)

const (
	USER_ID = "userId"
)

func Route(group *echo.Group, userService service.IUserService) {

	userRoute := NewUserRoute(userService)
	group.POST("/user", userRoute.InsertUser)
	group.GET("/user", userRoute.GetUsers)
	group.GET("/user/:"+USER_ID, userRoute.GetUser)
	group.PUT("/user/:"+USER_ID, userRoute.UpdateUser)
	group.DELETE("/user/:"+USER_ID, userRoute.DeleteUser)
}
