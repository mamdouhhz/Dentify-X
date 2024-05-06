package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Doctor{}, &DoctorRequests{}, &Patient{}, &DoctorPatient{}, &Xray{}, &DoctorXray{}, &Admin{})
	return db, nil
}
