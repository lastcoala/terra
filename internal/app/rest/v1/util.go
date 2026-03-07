package v1

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func queryParamToOffsetLimit(c echo.Context, useDefault bool) (int, int) {
	var page, limit int

	err := echo.QueryParamsBinder(c).Int("page", &page).
		Int("limit", &limit).BindError()

	if err != nil {
		return 0, 10
	}

	if page < 1 {
		if useDefault {
			return 0, 10
		} else {
			return 0, 0
		}
	}

	return (page - 1) * limit, limit
}

func getUserId(c echo.Context) (int, error) {
	id, err := strconv.Atoi(c.Param(USER_ID))
	if err != nil {
		return 0, err
	}
	return id, nil
}
