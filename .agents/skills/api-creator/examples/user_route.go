package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/internal/app/service"
)

type UserRoute struct {
	userService service.IUserService
}

func NewUserRoute(userService service.IUserService) *UserRoute {
	return &UserRoute{userService: userService}
}

// InsertUser godoc
//
// @Summary Insert user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body InsertUserReq true "User data"
// @Success 200	{object}	ResponseDoc{data=DataDoc{user=UserResp}}	"Success"
// @Router /user [post]
func (r *UserRoute) InsertUser(c echo.Context) error {
	var req InsertUserReq
	if err := c.Bind(&req); err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusBadRequest, resp)
	}
	user, err := r.userService.InsertUser(c.Request().Context(),
		req.ToDomain(), "admin")
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusInternalServerError, resp)
	}
	resp := NewResponseDto(MSG_SUCCESS, userToResp(user), "user")
	return c.JSON(http.StatusOK, resp)
}

// GetUsers godoc
//
// @Summary Get users
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param page  query int false "Page number (1-based)"
// @Param limit query int false "Page size"
// @Success 200 {object} ResponseDoc{data=DataDoc{users=[]UserResp}} "Success"
// @Router /user [get]
func (r *UserRoute) GetUsers(c echo.Context) error {
	offset, limit := queryParamToOffsetLimit(c, true)

	users, err := r.userService.GetUsers(c.Request().Context(), offset, limit)
	if err != nil {
		resp := NewResponsesDto[UserResp](err.Error(), nil, "users")
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp := NewResponsesDto(MSG_SUCCESS, usersToResp(users), "users")
	return c.JSON(http.StatusOK, resp)
}

// GetUser godoc
//
// @Summary Get user by ID
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} ResponseDoc{data=DataDoc{user=UserResp}} "Success"
// @Router /user/{userId} [get]
func (r *UserRoute) GetUser(c echo.Context) error {
	id, err := getUserId(c)
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusBadRequest, resp)
	}

	user, err := r.userService.GetUser(c.Request().Context(), id)
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp := NewResponseDto(MSG_SUCCESS, userToResp(user), "user")
	return c.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
//
// @Summary Update user
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param user body UpdateUserReq true "User data"
// @Success 200 {object} ResponseDoc{data=DataDoc{user=UserResp}} "Success"
// @Router /user/{userId} [put]
func (r *UserRoute) UpdateUser(c echo.Context) error {
	id, err := getUserId(c)
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusBadRequest, resp)
	}

	var req UpdateUserReq
	if err := c.Bind(&req); err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusBadRequest, resp)
	}

	user, err := r.userService.UpdateUser(c.Request().Context(), domain.User{
		Id:     id,
		Name:   req.Name,
		Gender: req.Gender,
		Email:  req.Email,
	}, "admin")
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp := NewResponseDto(MSG_SUCCESS, userToResp(user), "user")
	return c.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
//
// @Summary Delete user
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} ResponseDoc{} "Success"
// @Router /user/{userId} [delete]
func (r *UserRoute) DeleteUser(c echo.Context) error {
	id, err := getUserId(c)
	if err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := r.userService.DeleteUser(c.Request().Context(), id, "admin"); err != nil {
		resp := NewResponseDto(err.Error(), nil, "user")
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp := NewResponseDto(MSG_SUCCESS, nil, "user")
	return c.JSON(http.StatusOK, resp)
}
