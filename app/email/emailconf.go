package email

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func PatientConfirmationEmail(email string, name string, pass string, c *gin.Context) {
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

func DoctorConfirmationEmail(email string, name string, c *gin.Context) {
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

func DoctorAcceptanceEmail(email string, name string, c *gin.Context) {
	m := gomail.NewMessage()

	m.SetHeader("From", "dentifyx24@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Test Email")

	messageBody := "hi dr " + name + " you can now login to our system"
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

	messageBody := "hi " + name + " you are rejected"
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
