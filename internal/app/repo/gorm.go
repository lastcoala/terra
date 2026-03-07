package repo

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGormRepo(path string, nConn int) (*gorm.DB, error) {
	var err error
	client, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{Conn: client}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	if err != nil {
		return nil, err
	}

	dbX, err := db.DB()
	if err != nil {
		return nil, err
	}
	dbX.SetMaxIdleConns(nConn)
	dbX.SetMaxOpenConns(nConn)
	dbX.SetConnMaxLifetime(1 * time.Hour)

	return db, nil
}
