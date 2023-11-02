package database

import (
	"context"
	"fmt"
	"jobportalapi/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connection() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=root dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	pg, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err = pg.PingContext(ctx)
	if err != nil {

		return nil, fmt.Errorf("database is not connected: %w", err)
	}

	db.Migrator().AutoMigrate(&models.User{}, &models.Company{}, &models.Job{})
	return db, err

}
