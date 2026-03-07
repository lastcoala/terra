package v1

import "github.com/lastcoala/terra/internal/app/domain"

// UserResp is the user response DTO
type UserResp struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
}

// InsertUserReq is the insert user request DTO
type InsertUserReq struct {
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u InsertUserReq) ToDomain() domain.User {
	return domain.User{
		Name:     u.Name,
		Gender:   u.Gender,
		Email:    u.Email,
		Password: u.Password,
	}
}

// UpdateUserReq is the update user request DTO
type UpdateUserReq struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Email  string `json:"email"`
}

// ChangePasswordReq is the change password request DTO
type ChangePasswordReq struct {
	NewPassword string `json:"new_password"`
}

func userToResp(u domain.User) UserResp {
	return UserResp{
		Id:     u.Id,
		Name:   u.Name,
		Gender: u.Gender,
		Email:  u.Email,
	}
}

func usersToResp(users []domain.User) []UserResp {
	userResp := make([]UserResp, len(users))
	for i, u := range users {
		userResp[i] = userToResp(u)
	}
	return userResp
}
