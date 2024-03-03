package services

import (
	"Dentify-X/app/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Psignup(db *gorm.DB, c *gin.Context) error {
	var user models.Patient
	var existingUser models.Patient
	var err error

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	fmt.Printf("User after ShouldBindJSON: %+v\n", user)

	if err := db.Where("p_phone_number = ?", user.P_PhoneNumber).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return err
	}

	user.P_Password, err = HashPassword(user.P_Password)
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}