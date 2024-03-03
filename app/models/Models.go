package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	DoctorID      uint   `gorm:"primary_key;autoIncrement" json:"doctor_id"`
	D_Name        string `json:"name"`
	D_PhoneNumber string `json:"phone_number"`
	D_Password    string `json:"password"`
	D_Gender      string `json:"gender"`
	D_Email       string `json:"email"`
	ClinicAddress string `json:"clinic_address"`
}

type Patient struct {
	gorm.Model
	PatientID      uint   `gorm:"primary_key;autoIncrement" json:"patient_id"`
	Passcode       string `json:"passcode"`
	MedicalHistory string `json:"medical_history"` // question.
	P_Name         string `json:"name"`
	P_Gender       string `json:"gender"`
	P_PhoneNumber  string `json:"phone_number"`
	P_Email        string `json:"email"`
	P_Password     string `json:"password"`
}

type Xray struct {
	gorm.Model
	XrayID    uint `gorm:"primary_key;autoIncrement" json:"xray_id"`
	PatientID uint `gorm:"references:PatientID" json:"patient_id"` // not sure if it is a foreign key or a primary key.
	// where are the images that are going to be stores ??
}

type DoctorXray struct {
	gorm.Model
	DoctorID     uint      `gorm:"references:DoctorID" json:"doctor_id"`
	XrayID       uint      `gorm:"references:XrayID" json:"xray_id"`
	PatientID    uint      `gorm:"references:PatientID" json:"patient_id"`
	Prescription string    `json:"prescription"` // not sure.
	Date         time.Time `json:"date"`         // not sure.
}

type DoctorRequests struct {
	gorm.Model
	DoctorID      uint   `gorm:"primary_key;autoIncrement" json:"doctor_id"`
	D_Name        string `json:"name"`
	D_PhoneNumber string `json:"phone_number"`
	D_Password    string `json:"password"`
	D_Gender      string `json:"gender"`
	D_Email       string `json:"email"`
	ClinicAddress string `json:"clinic_address"`
}