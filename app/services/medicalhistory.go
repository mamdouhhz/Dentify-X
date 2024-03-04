package services

// import (
// 	"Dentify-X/app/models"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// // func GetMedicalHistory(db *gorm.DB, c *gin.Context) {
// // 	var p models.Patient
// // 	var medicalHistory []models.DoctorXray

// // 	if err := c.ShouldBindJSON(&p); err != nil {
// // 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
// // 		return
// // 	}

// // 	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", p.PatientID).Find(&medicalHistory).Error; err != nil {
// // 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
// // 		return
// // 	}
// // 	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
// // }

// func GetMedicalHistory(db *gorm.DB, c *gin.Context) {
// 	patientID, ok := c.Get("patientID")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	var medicalHistory []models.DoctorXray

// 	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", patientID).Find(&medicalHistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
// }
