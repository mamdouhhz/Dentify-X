package routers

import (
	"Dentify-X/app/handlers"
	"Dentify-X/app/middleware"
	"Dentify-X/app/services"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"gorm.io/gorm"
)

func Rout(db *gorm.DB) *gin.Engine {
	middleware.SaveLogs()

	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(), gindump.Dump())

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/psignup", func(c *gin.Context) {
		handlers.PsignupHandler(db, c)
	})

	r.POST("/dsignupreq", func(c *gin.Context) {
		handlers.DsignupHandler(db, c)
	})

	r.POST("/plogin", func(c *gin.Context) {
		handlers.Ploginhandler(db, c)
	})

	r.POST("/alogin", func(c *gin.Context) {
		handlers.Aloginhandler(db, c)
	})

	r.POST("/accept-request/:id", func(c *gin.Context) {
		doctorRequestID := c.Param("id")
		idUint, err := strconv.ParseUint(doctorRequestID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctorRequestID"})
			return
		}
		services.AcceptDoctorRequest(db, c, uint(idUint))
	})

	r.POST("/decline-request/:id", func(c *gin.Context) {
		doctorRequestID := c.Param("id")
		idUint, err := strconv.ParseUint(doctorRequestID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctorRequestID"})
			return
		}
		services.DeclineDoctorRequest(db, c, uint(idUint))
	})

	return r
}
