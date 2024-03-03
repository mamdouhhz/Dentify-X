package handlers

import (
	"Dentify-X/app/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DsignupHandler(db *gorm.DB, c *gin.Context) {
	if err := services.DoctorSignupRequest(db, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("can't create this user: %v", err)})
		return
	}
}
