package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/hashing"
	"Dentify-X/app/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PatientSignup(db *gorm.DB, c *gin.Context) error {
	var user models.Patient
	var existingUser models.Patient
	var err error

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	fmt.Printf("User after ShouldBindJSON: %+v\n", user)

	if err := db.Where("p_email = ?", user.P_Email).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return err
	}

	user.P_Password, err = hashing.HashPassword(user.P_Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	email.PatientWelcomeEmail(user.P_Email, user.P_Name, user.Passcode, c)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	return nil
}

func PatientLogin(db *gorm.DB, c *gin.Context) error {
	var user models.Patient
	var existingUser models.Patient

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("p_email = ?", user.P_Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "you are not signed up"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.P_Password), []byte(user.P_Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}
	session := sessions.Default(c)
	session.Set("pid", existingUser.PatientID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.P_Name})
	GetMedicalHistory(db, c, session)
	return nil
}

func GetMedicalHistory(db *gorm.DB, c *gin.Context, s sessions.Session) {
	patientID := s.Get("pid")
	if patientID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	var medicalHistory []models.DoctorXray
	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", patientID).Find(&medicalHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
}

func PatientConfirmPasswordReset(email string, db *gorm.DB, c *gin.Context) {
	var patient models.Patient
	password := c.Param("password")
	confirmpassword := c.Param("confirmpassword")

	if err := db.Where("p_email = ?", email).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found, wrong email"})
		return
	}

	if password != confirmpassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	patient.P_Password = password
	if err := db.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
