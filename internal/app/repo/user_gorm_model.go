package repo

import "github.com/lastcoala/terra/internal/app/domain"

type UserGormModel struct {
	BaseModel
	Name     string `gorm:"not null"`
	Gender   string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

func NewUserGormModel(user domain.User) UserGormModel {
	return UserGormModel{
		BaseModel: BaseModel{Id: user.Id},
		Name:      user.Name,
		Gender:    user.Gender,
		Email:     user.Email,
		Password:  user.Password,
	}
}

func (m UserGormModel) ToDomain() domain.User {
	return domain.User{
		Id:       m.Id,
		Name:     m.Name,
		Gender:   m.Gender,
		Email:    m.Email,
		Password: m.Password,
	}
}

func (m UserGormModel) TableName() string {
	return "users"
}

type UserGormModels []UserGormModel

func (m UserGormModels) ToDomains() []domain.User {
	var users []domain.User
	for _, user := range m {
		users = append(users, user.ToDomain())
	}
	return users
}

func (m UserGormModels) TableName() string {
	return "users"
}
