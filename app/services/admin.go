package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/models"
	"errors"
	"net/http"
	"text/template"

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

	// if err := bcrypt.CompareHashAndPassword([]byte(existingAdmin.A_password), []byte(admin.A_password)); err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
	// 	return err
	// }
	session := sessions.Default(c)
	session.Set("aid", existingAdmin.AdminID)

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"welcome": existingAdmin.A_Name})
	GetDoctorRequests(db, c)
	return nil
}

// Define an HTML template
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Doctor Requests</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body style="background-color:#dfebeb">
    <h1>Doctor Requests</h1>
    <table>
        <thead>
            <tr>
                <th>Name</th>
				<th>phone</th>
				<th>password</th>
				<th>mln</th>
				<th>gender</th>
				<th>email</th>
				<th>clinic</th>
            </tr>
        </thead>
        <tbody>
            {{range .DoctorRequests}}
            <tr>
                <td>{{.D_Name}}</td>
				<td>{{.D_PhoneNumber}}</td>
				<td>{{.D_Password}}</td>
				<td>{{.MLN}}</td>
				<td>{{.D_Gender}}</td>
				<td>{{.D_Email}}</td>
				<td>{{.ClinicAddress}}</td>
				<td>
				<form method="POST" action="http://localhost:8080/accept-request">
					<input type="hidden" name="doctorRequestIDaccept" value="{{.ID}}">
					<button type="submit" class="accept-btn">Accept</button>
        		</form>

				<form method="POST" action="http://localhost:8080/decline-request">
					<input type="hidden" name="doctorRequestIDreject" value="{{.ID}}">
					<button type="submit" class="accept-btn">Reject</button>
        		</form>
				</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`

// Define GetDoctorRequests function
func GetDoctorRequests(db *gorm.DB, c *gin.Context) {
	var doctorRequests []models.DoctorRequests
	if err := db.Find(&doctorRequests).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve doctor requests: %s", err.Error())
		return
	}

	// Execute HTML template
	t, err := template.New("DoctorRequests").Parse(htmlTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse template: %s", err.Error())
		return
	}

	// Render HTML template with data and send it as response
	if err := t.Execute(c.Writer, gin.H{"DoctorRequests": doctorRequests}); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render template: %s", err.Error())
		return
	}
}

// Define an HTML template
const htmlTemplatee = `
<!DOCTYPE html>
<html>
<head>
    <title>Doctors</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body style="background-color:#dfebeb">
    <h1>Doctors</h1>
    <table>
        <thead>
            <tr>
                <th>Name</th>
				<th>phone</th>
				<th>password</th>
				<th>mln</th>
				<th>gender</th>
				<th>email</th>
				<th>clinic</th>
            </tr>
        </thead>
        <tbody>
            {{range .Doctor}}
            <tr>
                <td>{{.D_Name}}</td>
				<td>{{.D_PhoneNumber}}</td>
				<td>{{.D_Password}}</td>
				<td>{{.MLN}}</td>
				<td>{{.D_Gender}}</td>
				<td>{{.D_Email}}</td>
				<td>{{.ClinicAddress}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`

func GetDoctors(db *gorm.DB, c *gin.Context) {
	var doctors []models.Doctor
	if err := db.Find(&doctors).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve doctors: %s", err.Error())
		return
	}

	// Execute HTML template
	t, err := template.New("Doctor").Parse(htmlTemplatee)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse template: %s", err.Error())
		return
	}

	// Render HTML template with data and send it as response
	if err := t.Execute(c.Writer, gin.H{"Doctor": doctors}); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render template: %s", err.Error())
		return
	}
}

// Define an HTML template
const htmlTemplateee = `
<!DOCTYPE html>
<html>
<head>
    <title>Patients</title>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body style="background-color:#dfebeb">
    <h1>Patients</h1>
    <table>
        <thead>
            <tr>
                <th>ID</th>
				<th>Passcode</th>
				<th>Name</th>
				<th>gender</th>
				<th>phone</th>
				<th>email</th>
				<th>password</th>
            </tr>
        </thead>
        <tbody>
            {{range .Patient}}
            <tr>
                <td>{{.PatientID}}</td>
				<td>{{.Passcode}}</td>
				<td>{{.P_Name}}</td>
				<td>{{.P_Gender}}</td>
				<td>{{.P_PhoneNumber}}</td>
				<td>{{.P_Email}}</td>
				<td>{{.P_Password}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`

func GetPatients(db *gorm.DB, c *gin.Context) {
	var patients []models.Patient
	if err := db.Find(&patients).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to retrieve doctors: %s", err.Error())
		return
	}

	// Execute HTML template
	t, err := template.New("Patient").Parse(htmlTemplateee)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse template: %s", err.Error())
		return
	}

	// Render HTML template with data and send it as response
	if err := t.Execute(c.Writer, gin.H{"patients": patients}); err != nil {
		c.String(http.StatusInternalServerError, "Failed to render template: %s", err.Error())
		return
	}
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
