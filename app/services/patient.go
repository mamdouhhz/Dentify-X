package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/hashing"
	"Dentify-X/app/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PatientSignup(db *gorm.DB, c *gin.Context) error {
	var user models.Patient
	var existingUser models.Patient
	var err error

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	fmt.Printf("User after ShouldBindJSON: %+v\n", user)

	if err := db.Where("p_email = ?", user.P_Email).First(&existingUser).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return err
	}

	user.P_Password, err = hashing.HashPassword(user.P_Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	user.Passcode = GenerateRandomPasscode()
	if err = db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	email.PatientWelcomeEmail(user.P_Email, user.P_Name, user.Passcode, c)
	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	return nil
}

func GenerateRandomPasscode() string {
	bytes := make([]byte, 4) // 4 bytes will give us an 8-digit hexadecimal number
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func PatientLogin(db *gorm.DB, c *gin.Context) error {
	var user models.Patient
	var existingUser models.Patient
	session := sessions.Default(c)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	if err := db.Where("p_email = ?", user.P_Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "you are not signed up"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.P_Password), []byte(user.P_Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return err
	}
	session.Set("pid", existingUser.PatientID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"welcome": existingUser.P_Name,
		"sessionid": session.Get("pid"),
	})
	return nil
}

// // GoogleLogin handles Google OAuth login
// func GoogleLogin(c *gin.Context) {
// 	var req struct {
// 		IDToken string `json:"id_token"`
// 	}
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Configure OAuth client credentials
// 	conf := &oauth2.Config{
// 		ClientID:     "1017790372651-ljdrkqcqbm08truclh6p8hf3ve7n9hm3.apps.googleusercontent.com", // Replace with your actual client ID
// 		ClientSecret: "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",                                       // Replace with your actual client secret
// 		RedirectURL:  "http://localhost:8000",                                                     // Must match your Google OAuth settings
// 		Scopes:       []string{"profile", "email"},
// 		Endpoint:     google.Endpoint,
// 	}

// 	// Verify ID token with Google
// 	payload, err := idtoken.Validate(context.Background(), req.IDToken, conf.ClientID)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
// 		return
// 	}

// 	// Extract user information from payload
// 	var userInfo struct {
// 		P_Name  string `json:"name"`
// 		P_Email string `json:"email"`
// 	}

// 	userInfo.P_Name = payload.Claims["name"].(string)
// 	userInfo.P_Email = payload.Claims["email"].(string)

// 	// Implement your own logic to handle user authentication and session management
// 	// For example, check if user with userInfo.Email exists in your database

// 	// Dummy response for demonstration
// 	responseData := struct {
// 		Welcome   string `json:"welcome"`
// 		SessionID string `json:"sessionid"`
// 	}{
// 		Welcome:   "Welcome, " + userInfo.P_Name,
// 		SessionID: "123456", // Example session ID
// 	}

// 	c.JSON(http.StatusOK, responseData)
// }

// Handler to get medical history
func GetMedicalHistory(db *gorm.DB, c *gin.Context) {
	var requestBody struct {
		PatientID uint `form:"patient_id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestBody.PatientID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var medicalHistory []models.DoctorXray
	if err := db.Select("xray_pdf_path, predicted_pdf_path, Prescription, medicalhistory, doctor_id, created_at").Where("patient_id = ?", requestBody.PatientID).Find(&medicalHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
}

// func GetMedicalHistory(db *gorm.DB, c *gin.Context) {
// 	var requestBody struct {
// 		PatientID uint `form:"patient_id" binding:"required"`
// 	}

// 	// Bind query parameters to RequestBody struct
// 	if err := c.ShouldBindQuery(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Ensure patient ID is provided
// 	if requestBody.PatientID == 0 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	var medicalHistory []models.DoctorXray
// 	if err := db.Select("xray_id, predicted_xray, doctor_id, medicalhistory, date").Where("patient_id = ?", requestBody.PatientID).Find(&medicalHistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
// }

func PatientConfirmPasswordReset(email string, db *gorm.DB, c *gin.Context) {
	var patient models.Patient
	password := c.Param("password")
	confirmpassword := c.Param("confirmpassword")

	if err := db.Where("p_email = ?", email).First(&patient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found, wrong email"})
		return
	}

	if password != confirmpassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	patient.P_Password = password
	if err := db.Save(&patient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func PatientLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"session": session.Get("pid")})
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
