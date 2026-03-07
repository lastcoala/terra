package repo

import "time"

type BaseModel struct {
	Id        int `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string
	IsDeleted bool
}
