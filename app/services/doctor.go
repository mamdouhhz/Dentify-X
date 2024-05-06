package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/hashing"
	"Dentify-X/app/models"
	"errors"
	"net/http"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	email.DoctorWelcomeEmail(user.D_Email, user.D_Name, c)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	return nil
}

func Doctorlogin(db *gorm.DB, c *gin.Context) error {
	var user models.Doctor
	var pengingUser models.DoctorRequests
	var existingUser models.Doctor

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("d_email = ?", user.D_Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Where("d_email = ?", user.D_Email).First(&pengingUser).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "you are not signed up"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return err
			}
			c.JSON(http.StatusOK, gin.H{"welcome": "you are still pending approval from our admins"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.D_Password), []byte(user.D_Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}
	session := sessions.Default(c)
	session.Set("did", existingUser.DoctorID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.D_Name})
	return nil
}

func AddPatient(db *gorm.DB, c *gin.Context) error {
	var requestData struct {
		DoctorID  uint   `json:"doctor_id"`
		Passcode  string `json:"passcode"`
		PatientID uint   `json:"patient_id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	var addedPatient models.DoctorPatient
	if err := db.Where("doctor_id = ? AND patient_id = ?", requestData.DoctorID, requestData.PatientID).First(&addedPatient).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "patient already added"})
		return err
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	var existingDoctor models.Doctor
	if err := db.Where("doctor_id = ?", requestData.DoctorID).First(&existingDoctor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "doctor does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	var existingPatient models.Patient
	if err := db.Where("passcode = ? AND patient_id = ?", requestData.Passcode, requestData.PatientID).First(&existingPatient).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient does not exist"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	doctorXray := models.DoctorPatient{
		DoctorID:  requestData.DoctorID,
		PatientID: requestData.PatientID,
	}

	if err := db.Create(&doctorXray).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "patient record added successfully to DoctorPatient"})
	return nil
}

func ViewPatientHistory(db *gorm.DB, c *gin.Context) {
	patientID := c.Param("pid")
	var medicalHistory models.DoctorXray

	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", patientID).Find(&medicalHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
	// na2esha validation: patient is assigned to doctor or not.
}

func UploadXray(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request"})
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".png": true, ".bmp": true, ".tiff": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, BMP, and TIFF files are allowed."})
		return
	}

	savePath := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/Dentify-X/cmd/" + file.Filename
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	cmd := exec.Command("python", "/Users/mamdouhhazem/yolov5/detect.py", "--weights", "/Users/mamdouhhazem/Desktop/Graduaiton_Project/DATASET/runs_results/newDatasetV5L/weights/best.pt", "--img", "512", "256", "--source", savePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform inference", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"detections": string(out)})
}

func DoctorConfirmPasswordReset(email string, db *gorm.DB, c *gin.Context) {
	var doctor models.Doctor
	password := c.Param("password")
	confirmpassword := c.Param("confirmpassword")

	if err := db.Where("d_email = ?", email).First(&doctor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found, wrong email"})
		return
	}

	if password != confirmpassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	doctor.D_Password = password
	if err := db.Save(&doctor).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
