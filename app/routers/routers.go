package routers

import (
	"Dentify-X/app/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Rout(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/psignup", func(c *gin.Context) {
		handlers.PsignupHandler(db, c)
	})

	r.POST("/dsignupreq", func(c *gin.Context) {
		handlers.DsignupHandler(db, c)
	})

	return r
}
