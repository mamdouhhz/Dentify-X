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

type DoctorRequests struct {
	gorm.Model
	DoctorID      uint   `gorm:"primary_key;autoIncrement" json:"doctor_id"`
	D_Name        string `json:"name"`
	D_PhoneNumber string `json:"phone_number"`
	D_Password    string `json:"password"`
	MLN           string `json:"mln"`
	D_Gender      string `json:"gender"`
	D_Email       string `json:"email"`
	ClinicAddress string `json:"clinic_address"`
}

type Patient struct {
	gorm.Model
	PatientID     uint   `gorm:"primary_key;autoIncrement" json:"patient_id" unique:"true"`
	Passcode      string `json:"passcode"`
	P_Name        string `json:"name"`
	P_Gender      string `json:"gender"`
	P_PhoneNumber string `json:"phone_number"`
	P_Email       string `json:"email"`
	P_Password    string `json:"password"`
}

type DoctorPatient struct {
	DoctorID  uint `gorm:"references:DoctorID" json:"doctor_id"`
	PatientID uint `gorm:"references:PatientID" json:"patient_id"`
}

type Xray struct {
	gorm.Model
	XrayID    uint   `gorm:"primary_key;autoIncrement" json:"xray_id"`
	PatientID uint   `gorm:"references:PatientID" json:"patient_id"` // not sure if it is a foreign key or a primary key.
	XrayImage []byte `json:"xray_image"`
}

type DoctorXray struct {
	gorm.Model
	MedicalHistory   uint      `gorm:"primary_key;autoIncrement;column:medicalhistory" json:"medicalhistory"`
	DoctorID         uint      `gorm:"column:doctor_id" json:"doctor_id"`
	PatientID        uint      `gorm:"column:patient_id" json:"patient_id"`
	XrayID           []byte    `gorm:"column:xray_id;type:bytea" json:"XrayID"`
	PredictedXray    []byte    `gorm:"column:predicted_xray;type:bytea" json:"PredictedXray"`
	XrayPDFPath      string    `json:"xray_pdf_path"`
	PredictedPDFPath string    `json:"predicted_pdf_path"`
	Prescription     string    `gorm:"column:prescription" json:"Prescription"`
	Date             time.Time `gorm:"column:date" json:"date"`
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
