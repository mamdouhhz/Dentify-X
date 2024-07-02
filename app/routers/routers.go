package routers

import (
	"Dentify-X/app/email"
	"Dentify-X/app/handlers"
	"Dentify-X/app/middlewares"
	"Dentify-X/app/services"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
	"gorm.io/gorm"
)

func Rout(db *gorm.DB) *gin.Engine {
	middlewares.SaveLogs()
	r := gin.New()
	r.Use(gin.Recovery(), middlewares.Logger(), gindump.Dump())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://localhost/files"}
	config.AllowMethods = []string{"GET", "POST"}
	//config.AllowAllOrigins = true
	r.Use(cors.New(config))

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 5000})
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/resetPassConfEmail/:email", func(c *gin.Context) {
		email.EmailForgetPassword(db, c)
	})

	// Doctor
	r.POST("/dsignupreq", func(c *gin.Context) {
		services.DoctorSignupRequest(db, c)
	})
	r.POST("/dlogin", func(c *gin.Context) {
		services.Doctorlogin(db, c)
	})
	r.POST("/addpatient", func(c *gin.Context) {
		services.AddPatient(db, c)
	})
	r.POST("/existingpatient", func(c *gin.Context) {
		services.ExistingPatient(db, c)
	})
	r.POST("/upload", func(c *gin.Context) {
		services.UploadXray(db, c)
	})
	r.GET("/latest-predicted-image", services.ServeLatestPredictedImage)
	r.POST("/save-prescription", func(c *gin.Context) {
		services.CreatePrescriptionPDF(c, db)
	})

	// Patient
	r.POST("/psignup", func(c *gin.Context) {
		handlers.PsignupHandler(db, c)
	})
	r.POST("/plogin", func(c *gin.Context) {
		services.PatientLogin(db, c)
	})
	// r.POST("/google-login", func(c *gin.Context) {
	// 	services.GoogleLogin(c)
	// })
	r.GET("/medicalhistory", func(c *gin.Context) {
		services.GetMedicalHistory(db, c)
	})
	r.POST("/plogout", func(c *gin.Context) {
		services.PatientLogout(c)
	})

	// Admin
	r.POST("/alogin", func(c *gin.Context) {
		services.AdminLogin(db, c)
	})
	r.GET("/doctors", func(c *gin.Context) {
		services.GetDoctors(db, c)
	})
	r.GET("/patients", func(c *gin.Context) {
		services.GetPatients(db, c)
	})
	r.GET("/Requests", func(c *gin.Context) {
		services.GetDoctorRequests(db, c)
	})
	r.POST("/accept-request", func(c *gin.Context) {
		var requestData struct {
			DoctorRequestIDAccept uint64 `json:"doctorRequestIDaccept"`
		}

		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		services.AcceptDoctorRequest(db, c, uint(requestData.DoctorRequestIDAccept))
	})

	r.POST("/decline-request", func(c *gin.Context) {
		var requestData struct {
			DoctorRequestIDReject uint64 `json:"doctorRequestIDreject"`
		}

		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		services.DeclineDoctorRequest(db, c, uint(requestData.DoctorRequestIDReject))
	})
	return r
}
