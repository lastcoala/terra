package v1

import (
	"errors"
	"net/http"
	"testing"

	"github.com/lastcoala/terra/internal/app/domain"
	"github.com/lastcoala/terra/internal/mocks"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// InsertUser
// ---------------------------------------------------------------------------

func TestUserRoute_InsertUser(t *testing.T) {
	t.Run("success insert user", func(t *testing.T) {
		req := InsertUserReq{
			Name:     "John Doe",
			Gender:   "Male",
			Email:    "john.doe@example.com",
			Password: "password",
		}
		user := domain.User{
			Id:       1,
			Name:     req.Name,
			Gender:   req.Gender,
			Email:    req.Email,
			Password: req.Password,
		}

		c, rec := requestJsonTestHelper("POST", req, "/user")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().InsertUser(c.Request().Context(), req.ToDomain(), "admin").Return(user, nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.InsertUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"user":{"id":1,"name":"John Doe","gender":"Male","email":"john.doe@example.com"}},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("service error returns 500", func(t *testing.T) {
		req := InsertUserReq{
			Name:     "John Doe",
			Gender:   "Male",
			Email:    "john.doe@example.com",
			Password: "password",
		}

		c, rec := requestJsonTestHelper("POST", req, "/user")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().InsertUser(c.Request().Context(), req.ToDomain(), "admin").Return(domain.User{}, errors.New("db error"))

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.InsertUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t,
			`{"message":"db error","data":{"user":{}},"version":"-"}`,
			rec.Body.String())
	})
}

// ---------------------------------------------------------------------------
// GetUsers
// ---------------------------------------------------------------------------

func TestUserRoute_GetUsers(t *testing.T) {
	t.Run("success get users", func(t *testing.T) {
		users := []domain.User{
			{Id: 1, Name: "Alice", Gender: "Female", Email: "alice@example.com"},
			{Id: 2, Name: "Bob", Gender: "Male", Email: "bob@example.com"},
		}

		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user?page=1&limit=10")
		// echo parses query params from the URL, so we need to set them on the context
		c.QueryParams().Set("page", "1")
		c.QueryParams().Set("limit", "10")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().GetUsers(c.Request().Context(), 0, 10).Return(users, nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUsers(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"users":[{"id":1,"name":"Alice","gender":"Female","email":"alice@example.com"},{"id":2,"name":"Bob","gender":"Male","email":"bob@example.com"}]},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("no query params defaults to page=1 limit=10", func(t *testing.T) {
		users := []domain.User{
			{Id: 1, Name: "Alice", Gender: "Female", Email: "alice@example.com"},
		}

		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().GetUsers(c.Request().Context(), 0, 10).Return(users, nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUsers(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"users":[{"id":1,"name":"Alice","gender":"Female","email":"alice@example.com"}]},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("service error returns 500", func(t *testing.T) {
		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().GetUsers(c.Request().Context(), 0, 10).Return(nil, errors.New("db error"))

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUsers(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t,
			`{"message":"db error","data":{"users":[]},"version":"-"}`,
			rec.Body.String())
	})
}

// ---------------------------------------------------------------------------
// GetUser
// ---------------------------------------------------------------------------

func TestUserRoute_GetUser(t *testing.T) {
	t.Run("success get user by id", func(t *testing.T) {
		user := domain.User{Id: 1, Name: "Alice", Gender: "Female", Email: "alice@example.com"}

		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().GetUser(c.Request().Context(), 1).Return(user, nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"user":{"id":1,"name":"Alice","gender":"Female","email":"alice@example.com"}},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("invalid user id returns 400", func(t *testing.T) {
		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user/abc")
		c.SetParamNames(USER_ID)
		c.SetParamValues("abc")

		mockUserService := mocks.NewMockIUserService(t)
		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		c, rec := requestJsonTestHelper("GET", struct{}{}, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().GetUser(c.Request().Context(), 1).Return(domain.User{}, errors.New("not found"))

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.GetUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t,
			`{"message":"not found","data":{"user":{}},"version":"-"}`,
			rec.Body.String())
	})
}

// ---------------------------------------------------------------------------
// UpdateUser
// ---------------------------------------------------------------------------

func TestUserRoute_UpdateUser(t *testing.T) {
	t.Run("success update user", func(t *testing.T) {
		req := UpdateUserReq{
			Name:   "Alice Updated",
			Gender: "Female",
			Email:  "alice.updated@example.com",
		}
		updatedUser := domain.User{Id: 1, Name: req.Name, Gender: req.Gender, Email: req.Email}

		c, rec := requestJsonTestHelper("PUT", req, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().UpdateUser(c.Request().Context(), domain.User{
			Id:     1,
			Name:   req.Name,
			Gender: req.Gender,
			Email:  req.Email,
		}, "admin").Return(updatedUser, nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.UpdateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"user":{"id":1,"name":"Alice Updated","gender":"Female","email":"alice.updated@example.com"}},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("invalid user id returns 400", func(t *testing.T) {
		req := UpdateUserReq{Name: "Alice", Gender: "Female", Email: "alice@example.com"}

		c, rec := requestJsonTestHelper("PUT", req, "/user/abc")
		c.SetParamNames(USER_ID)
		c.SetParamValues("abc")

		mockUserService := mocks.NewMockIUserService(t)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.UpdateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		req := UpdateUserReq{Name: "Alice", Gender: "Female", Email: "alice@example.com"}

		c, rec := requestJsonTestHelper("PUT", req, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().UpdateUser(c.Request().Context(), domain.User{
			Id:     1,
			Name:   req.Name,
			Gender: req.Gender,
			Email:  req.Email,
		}, "admin").Return(domain.User{}, errors.New("db error"))

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.UpdateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t,
			`{"message":"db error","data":{"user":{}},"version":"-"}`,
			rec.Body.String())
	})
}

// ---------------------------------------------------------------------------
// DeleteUser
// ---------------------------------------------------------------------------

func TestUserRoute_DeleteUser(t *testing.T) {
	t.Run("success delete user", func(t *testing.T) {
		c, rec := requestJsonTestHelper("DELETE", struct{}{}, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().DeleteUser(c.Request().Context(), 1, "admin").Return(nil)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.DeleteUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t,
			`{"message":"success","data":{"user":{}},"version":"-"}`,
			rec.Body.String())
	})

	t.Run("invalid user id returns 400", func(t *testing.T) {
		c, rec := requestJsonTestHelper("DELETE", struct{}{}, "/user/abc")
		c.SetParamNames(USER_ID)
		c.SetParamValues("abc")

		mockUserService := mocks.NewMockIUserService(t)

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.DeleteUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error returns 500", func(t *testing.T) {
		c, rec := requestJsonTestHelper("DELETE", struct{}{}, "/user/1")
		c.SetParamNames(USER_ID)
		c.SetParamValues("1")

		mockUserService := mocks.NewMockIUserService(t)
		mockUserService.EXPECT().DeleteUser(c.Request().Context(), 1, "admin").Return(errors.New("db error"))

		userRoute := NewUserRoute(mockUserService)
		err := userRoute.DeleteUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t,
			`{"message":"db error","data":{"user":{}},"version":"-"}`,
			rec.Body.String())
	})
}
