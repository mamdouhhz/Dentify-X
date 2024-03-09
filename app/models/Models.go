package models

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model
	DoctorID      uint   `gorm:"primary_key" json:"doctor_id"`
	D_Name        string `json:"name"`
	D_PhoneNumber string `json:"phone_number"`
	D_Password    string `json:"password"`
	MLN           string `json:"mln"`
	D_Gender      string `json:"gender"`
	D_Email       string `json:"email"`
	ClinicAddress string `json:"clinic_address"`
}

// for the admin to view this table to accpet or decline.
// if accepted, the record is inserted into Doctor table.
type DoctorRequests struct {
	gorm.Model
	DoctorID      uint   `gorm:"primary_key;autoIncrement" json:"doctor_id"`
	D_Name        string `json:"name"`
	D_PhoneNumber string `json:"phone_number"`
	D_Password    string `json:"d_password"`
	MLN           string `json:"mln"`
	D_Gender      string `json:"gender"`
	D_Email       string `json:"d_email"`
	ClinicAddress string `json:"clinic_address"`
}

type Patient struct {
	gorm.Model
	PatientID      uint   `gorm:"primary_key;autoIncrement" json:"patient_id" unique:"true"`
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
	// where are the images that are going to be stored ??
}

type DoctorXray struct {
	gorm.Model
	DoctorID     uint      `gorm:"references:DoctorID" json:"doctor_id"`
	XrayID       uint      `gorm:"references:XrayID" json:"xray_id"`
	PatientID    uint      `gorm:"references:PatientID" json:"patient_id"`
	Prescription string    `json:"prescription"` // not sure.
	Date         time.Time `json:"date"`         // not sure.
}

type Admin struct {
	gorm.Model
	AdminID       uint   `gorm:"primary_key;autoIncrement" json:"id" unique:"true"`
	A_Name        string `json:"name"`
	A_password    string `json:"password"`
	A_gender      string `json:"gender"`
	A_PhoneNumber string `json:"phone_number"`
	A_Email       string `json:"email"`
}
