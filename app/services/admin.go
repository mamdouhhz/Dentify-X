package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/models"
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminLogin(db *gorm.DB, c *gin.Context) error {
	var admin models.Admin
	var existingAdmin models.Admin

	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("a_email = ?", admin.A_Email).First(&existingAdmin).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "you are not a regiistered admin"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if err := db.Where("a_password = ?", admin.A_password).First(&existingAdmin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "wrong password"})
		return nil
	}

	session := sessions.Default(c)
	session.Set("aid", existingAdmin.AdminID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}
	c.JSON(http.StatusOK, gin.H{
		"welcome":   existingAdmin.A_Name,
		"sessionid": session.Get("aid"),
		"email":     existingAdmin.A_Email,
		"password":  existingAdmin.A_password,
		"phone":     existingAdmin.A_PhoneNumber,
	})
	return nil
}

func GetDoctorRequests(db *gorm.DB, c *gin.Context) {
	var DoctorRequests []models.DoctorRequests
	if err := db.Find(&DoctorRequests).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor requests", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"doctor requests:": DoctorRequests})
}

func GetDoctors(db *gorm.DB, c *gin.Context) {
	var doctors []models.Doctor
	if err := db.Find(&doctors).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve doctors: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"doctors: ": doctors})
}

func GetPatients(db *gorm.DB, c *gin.Context) {
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve doctors: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"patients: ": patients})
}

func AcceptDoctorRequest(db *gorm.DB, c *gin.Context, doctorRequestID uint) {
	var doctorRequest models.DoctorRequests
	if err := db.First(&doctorRequest, doctorRequestID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor request", "details": err.Error()})
		return
	}

	newDoctor := models.Doctor{
		DoctorID:      doctorRequest.DoctorID,
		D_Name:        doctorRequest.D_Name,
		D_PhoneNumber: doctorRequest.D_PhoneNumber,
		D_Password:    doctorRequest.D_Password,
		MLN:           doctorRequest.MLN,
		D_Gender:      doctorRequest.D_Gender,
		D_Email:       doctorRequest.D_Email,
		ClinicAddress: doctorRequest.ClinicAddress,
	}

	if err := db.Create(&newDoctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert doctor record", "details": err.Error()})
		return
	}

	if err := db.Delete(&doctorRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor request", "details": err.Error()})
		return
	}
	email.DoctorAcceptanceEmail(doctorRequest.D_Email, doctorRequest.D_Name, c)
	c.JSON(http.StatusOK, gin.H{"message": "Doctor request accepted and recored removed from requests and added to permenant doctors"})
}

func DeclineDoctorRequest(db *gorm.DB, c *gin.Context, doctorID uint) {
	var doctorRequest models.DoctorRequests
	if err := db.First(&doctorRequest, doctorID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor request", "details": err.Error()})
		return
	}

	if err := db.Delete(&models.DoctorRequests{}, doctorID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor request", "details": err.Error()})
		return
	}
	email.DoctorRejectionEmail(doctorRequest.D_Email, doctorRequest.D_Name, c)
	c.JSON(http.StatusOK, gin.H{"message": "Doctor request declined and record deleted"})
}
