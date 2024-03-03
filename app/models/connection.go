package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres dbname=Dentify-X password=123 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Doctor{}, &Patient{}, &Xray{}, &DoctorXray{}, &DoctorRequests{})

	return db, nil
}
