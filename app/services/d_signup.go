package services

import (
	"Dentify-X/app/hashing"
	"Dentify-X/app/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DoctorSignupRequest(db *gorm.DB, c *gin.Context) error {
	var user models.DoctorRequests
	var existingUser models.Doctor
	var pendingUser models.DoctorRequests
	var err error

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("mln = ?", user.MLN).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return err
	}
	if err := db.Where("mln = ?", user.MLN).First(&pendingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "you are still pending approval from our admins"})
		return err
	}

	user.D_Password, err = hashing.HashPassword(user.D_Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	return nil
}
