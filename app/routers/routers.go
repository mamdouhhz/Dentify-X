package routers

import (
	"Dentify-X/app/handlers"
	"Dentify-X/app/middleware"
	"Dentify-X/app/services"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST"}
	r.Use(cors.New(config))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// Root path
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Welcome to Dentify-X!")
	})

	// Doctor
	r.POST("/dsignupreq", func(c *gin.Context) {
		handlers.DsignupHandler(db, c)
	})
	r.POST("/dlogin", func(c *gin.Context) {
		handlers.Dloginhandler(db, c)
	})
	r.POST("/addpatient", func(c *gin.Context) {
		services.AddPatient(db, c)
	})

	// Patient
	r.POST("/psignup", func(c *gin.Context) {
		handlers.PsignupHandler(db, c)
	})
	r.POST("/plogin", func(c *gin.Context) {
		handlers.Ploginhandler(db, c)
	})

	// Admin
	r.POST("/alogin", func(c *gin.Context) {
		handlers.Aloginhandler(db, c)
	})
	r.GET("/getrequests", func(c *gin.Context) {
		services.GetDoctorRequests(db, c)
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
