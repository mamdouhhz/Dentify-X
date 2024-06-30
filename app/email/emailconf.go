package email

import (
	"Dentify-X/app/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func PatientWelcomeEmail(email string, name string, pass string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "Welcome " + name + " to Dentify-X, this is your confirmation email, your passcode: " + pass
	m.SetBody("text/plain", messageBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func DoctorWelcomeEmail(email string, name string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "Welcome " + name + " to Dentify-X, your are pending approval from our admins"
	m.SetBody("text/plain", messageBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func PendingDoctorEmail(email string, name string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "Hi " + name + " You are still pending approval from our admins, you will recieve a confirmation email once our admins review your request"
	m.SetBody("text/plain", messageBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func DoctorAcceptanceEmail(email string, name string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "Welcome dr " + name + " you can now login to our system"
	m.SetBody("text/plain", messageBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func DoctorRejectionEmail(email string, name string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "Hi " + name + " you are rejected"
	m.SetBody("text/plain", messageBody)

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func EmailForgetPassword(db *gorm.DB, c *gin.Context) {
	var patient models.Patient
	var doctor models.Doctor

	m := gomail.NewMessage()
	email := c.Param("email")

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	if err := db.Where("p_email = ?", email).First(&patient).Error; err != nil {
		if err := db.Where("d_email = ?", email).First(&doctor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email not found"})
			return
		}
		DoctorMessageBody := "Click this link to reset your password, (link to Doctor)"
		m.SetBody("text/plain", DoctorMessageBody)

	} else {
		PatientMessageBody := "Click this link to reset your password: (link to patient)"
		m.SetBody("text/plain", PatientMessageBody)
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, "dentifyx24@gmail.com", "yyyx rysz tgef bxik")

	if err := d.DialAndSend(m); err != nil {
		errorMsg := fmt.Sprintf("Error sending email: %s", err)
		fmt.Println(errorMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errorMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
