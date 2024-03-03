package services

import (
	"Dentify-X/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PatientLogin(db *gorm.DB, c *gin.Context) error {

	var user models.Patient
	//var medical_history models.DoctorXray
	var existingUser models.Patient
	var err error

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
	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.P_Name})
	//db.Select("xray_id", "prescription", "date").Find(&medical_history).Rows()
	return err
}
