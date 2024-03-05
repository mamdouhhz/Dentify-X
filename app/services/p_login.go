package services

import (
	"Dentify-X/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var user models.Patient
var existingUser models.Patient

func PatientLogin(db *gorm.DB, c *gin.Context) error {

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
	// session := sessions.Default(c)
	// session.Set("pid", existingUser.PatientID)

	// if err := session.Save(); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
	// 	return err
	// }
	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.P_Name})
	GetMedicalHistory(db, c)
	return nil
}

func GetMedicalHistory(db *gorm.DB, c *gin.Context) {
	// patientID, ok := c.Get("pid")
	// c.JSON(http.StatusOK, gin.H{"value of PatientID after getting": patientID})

	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	var medicalHistory []models.DoctorXray

	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", existingUser.PatientID).Find(&medicalHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
}
