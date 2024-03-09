package services

import (
	"Dentify-X/app/models"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Doctorlogin(db *gorm.DB, c *gin.Context) error {
	var user models.Doctor
	var pengingUser models.DoctorRequests
	var existingUser models.Doctor

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("d_email = ?", user.D_Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			if err := db.Where("d_email = ?", user.D_Email).First(&pengingUser).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "you are not signed up"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return err
			}
			c.JSON(http.StatusOK, gin.H{"welcome": "you are still pending approval from our admins"})

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.D_Password), []byte(user.D_Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}
	session := sessions.Default(c)
	session.Set("did", existingUser.DoctorID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.D_Name})
	return nil
}

func AddPatient(db *gorm.DB, c *gin.Context) error {
	var requestData struct {
		DoctorID  uint   `json:"doctor_id"`
		Passcode  string `json:"passcode"`
		PatientID uint   `json:"patient_id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	var existingP models.DoctorXray
	if err := db.Where("doctor_id = ? AND patient_id = ?", requestData.DoctorID, requestData.PatientID).First(&existingP).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "patient already added"})
		return err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	var existingPatient models.Patient
	if err := db.Where("passcode = ? AND patient_id = ?", requestData.Passcode, requestData.PatientID).First(&existingPatient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	doctorXray := models.DoctorXray{
		DoctorID:  requestData.DoctorID,
		PatientID: requestData.PatientID,
	}

	if err := db.Create(&doctorXray).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "patient record added successfully to doctorXrays"})
	return nil
}