package handlers

import (
	"Dentify-X/app/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Dloginhandler(db *gorm.DB, c *gin.Context) {
	if err := services.Doctorlogin(db, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("can't Login: %v", err)})
		return
	}
}
