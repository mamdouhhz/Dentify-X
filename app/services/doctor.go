package services

import (
	"Dentify-X/app/email"
	"Dentify-X/app/hashing"
	"Dentify-X/app/models"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jung-kurt/gofpdf"
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
		email.PendingDoctorEmail(user.D_Email, user.D_Name, c)
		c.JSON(http.StatusConflict, gin.H{"error": "check your email, you are still pending approval from our admins"})
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
	c.JSON(http.StatusOK, gin.H{"message": "check your email"})
	return nil
}

// Doctorlogin handles doctor login
func Doctorlogin(db *gorm.DB, c *gin.Context) {
	var user models.Doctor
	var pendingUser models.DoctorRequests
	var existingUser models.Doctor

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Where("d_email = ?", user.D_Email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Where("d_email = ?", user.D_Email).First(&pendingUser).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "You are not signed up"})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				}
				return
			}
			c.JSON(http.StatusOK, gin.H{"error": "You are still pending approval from our admins"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.D_Password), []byte(user.D_Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("did", existingUser.DoctorID)

	// Save session
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"welcome":   existingUser.D_Name,
		"sessionid": session.Get("did"),
		"email":     existingUser.D_Email,
		"password":  existingUser.D_PhoneNumber,
		"phone":     existingUser.D_PhoneNumber,
	})
}

func AddPatient(db *gorm.DB, c *gin.Context) error {
	var requestData struct {
		PatientID uint   `json:"PatientID"`
		Passcode  string `json:"Passcode"`
		DoctorID  uint   `json:"DoctorID"`
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

	// var existingDoctor models.Doctor
	// if err := db.Where("doctor_id = ?", requestData.DoctorID).First(&existingDoctor).Error; err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "doctor does not exist"})
	// 	} else {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	}
	// 	return err
	// }

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

func ExistingPatient(db *gorm.DB, c *gin.Context) error {
	var doctorPatient models.DoctorPatient

	// Define a struct to bind the incoming JSON data
	type RequestBody struct {
		DoctorID  uint `json:"doctor_id"`
		PatientID uint `json:"patient_id"`
	}

	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}

	// Check if the patient exists for the given doctor
	if err := db.Where("patient_id = ? AND doctor_id = ?", requestBody.PatientID, requestBody.DoctorID).First(&doctorPatient).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient is not added to your list"})
		return err
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}
	c.JSON(http.StatusOK, gin.H{"message": "patient found"})
	return nil
}

func createPDF(imageBytes []byte, filename string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	imageData := imageBytes // Assuming imageBytes contains valid JPEG or PNG image data
	imageType := "JPG"      // Adjust based on your image type

	imageOptions := gofpdf.ImageOptions{
		ImageType: imageType,
	}

	// Add image to PDF
	pdf.RegisterImageOptionsReader(filename, imageOptions, bytes.NewReader(imageData))
	pdf.ImageOptions(filename, 10, 10, 190, 0, false, imageOptions, 0, "")

	// Set font for table
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0) // Set text color to black (RGB values: 0, 0, 0)

	// Add table
	tableX, tableY := 10.0, 150.0 // Adjust X, Y position of the table
	pdf.SetXY(tableX, tableY)
	cellWidth, cellHeight := 95.0, 10.0 // Adjust cell width and height

	staticData := [][]string{
		{"0: ", "Implants"},
		{"1: ", "Fillings"},
		{"3: ", "Impacted tooth"},
		{"4: ", "caries"},
	}

	for _, row := range staticData {
		for _, col := range row {
			pdf.CellFormat(cellWidth, cellHeight, col, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}

	// Add prescription text
	pdf.SetY(-30) // Adjust Y position based on your layout
	pdf.MultiCell(0, 10, "Prescription: Your prescription text here.", "", "C", false)

	return pdf.OutputFileAndClose(filename)
}

func UploadXray(db *gorm.DB, c *gin.Context) {
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
		return
	}

	patientID := c.PostForm("patient_id")
	doctorID := c.PostForm("doctor_id")
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request", "details": err.Error()})
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded file", "details": err.Error()})
		return
	}

	// Channels to handle results and errors
	pdfChan := make(chan string, 1)
	yoloChan := make(chan string, 1)
	errChan := make(chan error, 2)

	// Save original x-ray as PDF
	go func() {
		originalPDFPath := filepath.Join("/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/Project_Grad/htmlandcssandimages", header.Filename+".pdf")
		if err := createPDF(fileBytes, originalPDFPath); err != nil {
			errChan <- fmt.Errorf("failed to create PDF from original x-ray: %w", err)
			return
		}
		pdfChan <- originalPDFPath
	}()

	// Run YOLO prediction
	go func() {
		tempFilePath := filepath.Join("/tmp", header.Filename)
		if err := os.WriteFile(tempFilePath, fileBytes, 0644); err != nil {
			errChan <- fmt.Errorf("failed to save uploaded file: %w", err)
			return
		}
		defer os.Remove(tempFilePath)

		yoloCmd := "/Users/mamdouhhazem/opt/anaconda3/envs/dentex/bin/yolo"
		cmd := exec.Command(yoloCmd, "predict", "model=/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/runs_results/newDatasetV8M/weights/best.pt", "source="+tempFilePath)
		if err := cmd.Run(); err != nil {
			errChan <- fmt.Errorf("failed to perform inference: %w", err)
			return
		}

		runsDir := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/cmd/runs/detect"
		latestPredictFolder, err := findLatestPredictFolder(runsDir)
		if err != nil {
			errChan <- fmt.Errorf("failed to find latest predict folder: %w", err)
			return
		}

		predictedFilePath := filepath.Join(latestPredictFolder, header.Filename)
		predictedFileBytes, err := os.ReadFile(predictedFilePath)
		if err != nil {
			errChan <- fmt.Errorf("failed to read predicted file: %w", err)
			return
		}

		predictedPDFPath := filepath.Join("/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/Project_Grad/htmlandcssandimages", header.Filename+".predicted.pdf")
		if err := createPDF(predictedFileBytes, predictedPDFPath); err != nil {
			errChan <- fmt.Errorf("failed to create PDF from predicted x-ray: %w", err)
			return
		}
		yoloChan <- predictedPDFPath
	}()

	// Wait for results
	var originalPDFPath, predictedPDFPath string
	for i := 0; i < 2; i++ {
		select {
		case pdf := <-pdfChan:
			originalPDFPath = pdf
		case yolo := <-yoloChan:
			predictedPDFPath = yolo
		case err := <-errChan:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Convert IDs to uint
	patientIDUint, err := strconv.ParseUint(patientID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID", "details": err.Error()})
		return
	}
	doctorIDUint, err := strconv.ParseUint(doctorID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID", "details": err.Error()})
		return
	}

	// Save data to database
	medicalHistory := models.DoctorXray{
		DoctorID:         uint(doctorIDUint),
		PatientID:        uint(patientIDUint),
		XrayPDFPath:      originalPDFPath,
		PredictedPDFPath: predictedPDFPath,
		Date:             time.Now(),
	}
	if err := db.Create(&medicalHistory).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
}

// CreatePrescriptionPDF creates a prescription PDF and updates the database record
func CreatePrescriptionPDF(c *gin.Context, db *gorm.DB) {
	var request struct {
		Prescription string `json:"Prescription"`
		PatientID    uint   `json:"patient_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the existing record for the given patient ID
	var doctorXray models.DoctorXray
	if err := db.Where("patient_id = ?", request.PatientID).Last(&doctorXray).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record not found"})
		return
	}

	// Create a unique filename for the prescription PDF
	timestamp := time.Now().Format("2006-01-02-15-04-05") // Format: YYYY-MM-DD-HH-MM-SS
	filename := fmt.Sprintf("%d_prescription_%s.pdf", request.PatientID, timestamp)

	// Generate PDF content
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Doctor's Prescription")
	pdf.Ln(10)

	// Prescription text
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, "Patient ID: "+fmt.Sprint(request.PatientID), "", "L", false)
	pdf.MultiCell(0, 10, "Prescription: "+request.Prescription, "", "L", false)
	pdf.Ln(10)

	// Doctor's signature and timestamp
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(0, 10, fmt.Sprintf("Doctor's Signature: ______________________   Date: %s", time.Now().Format("2006-01-02")), "", 0, "L", false, 0, "")
	pdf.Ln(10)

	// Save the prescription PDF file directly
	filepath := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/Project_Grad/htmlandcssandimages/" + filename
	if err := pdf.OutputFileAndClose(filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save prescription PDF", "details": err.Error()})
		return
	}

	// Update the existing record with the prescription PDF path
	doctorXray.Prescription = filepath
	if err := db.Save(&doctorXray).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record", "details": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "Prescription updated successfully", "prescription_pdf_path": filepath})
}

// Helper function to find the latest predict folder
func findLatestPredictFolder(runsDir string) (string, error) {
	files, err := os.ReadDir(runsDir)
	if err != nil {
		return "", err
	}

	var latestFolder string
	var latestTime time.Time

	for _, file := range files {
		if file.IsDir() && strings.HasPrefix(file.Name(), "predict") {
			fileInfo, err := file.Info()
			if err != nil {
				return "", err
			}

			if fileInfo.ModTime().After(latestTime) {
				latestTime = fileInfo.ModTime()
				latestFolder = filepath.Join(runsDir, file.Name())
			}
		}
	}

	if latestFolder == "" {
		return "", fmt.Errorf("no predict folder found")
	}

	return latestFolder, nil
}

// Handler to serve the latest predicted image
func ServeLatestPredictedImage(c *gin.Context) {
	runsDir := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/cmd/runs/detect"
	latestPredictFolder, err := findLatestPredictFolder(runsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find latest predict folder", "details": err.Error()})
		return
	}

	files, err := os.ReadDir(latestPredictFolder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read files in predict folder", "details": err.Error()})
		return
	}

	if len(files) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No files found in the latest predict folder"})
		return
	}

	imagePath := filepath.Join(latestPredictFolder, files[0].Name())
	c.File(imagePath)
}

// // Handler to save prescription for the predicted x-ray
// func SavePrescription(db *gorm.DB, c *gin.Context) {
// 	// Parse prescription data from form or JSON body
// 	prescription := c.PostForm("prescription")
// 	if prescription == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Prescription is required"})
// 		return
// 	}

// 	// Fetch the latest predicted x-ray record from the database
// 	var latestXray models.DoctorXray
// 	if err := db.Order("id DESC").First(&latestXray).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find latest x-ray record", "details": err.Error()})
// 		return
// 	}

// 	// Update the prescription field
// 	latestXray.Prescription = prescription

// 	// Save updated record back to the database
// 	if err := db.Save(&latestXray).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save prescription", "details": err.Error()})
// 		return
// 	}

// 	// Read the existing PDF file
// 	pdfBytes, err := os.ReadFile(latestXray.PredictedPDFPath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF file", "details": err.Error()})
// 		return
// 	}

// 	// Append prescription to PDF content
// 	updatedPDFBytes := appendPrescriptionToPDF(pdfBytes, prescription)

// 	// Write updated PDF back to file (optional step)
// 	if err := os.WriteFile(latestXray.PredictedPDFPath, updatedPDFBytes, 0644); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated PDF file", "details": err.Error()})
// 		return
// 	}

// 	// Respond with URL of the last predicted PDF file
// 	pdfURL := fmt.Sprintf("http://localhost:8000/%s", filepath.Base(latestXray.PredictedPDFPath))
// 	c.JSON(http.StatusOK, gin.H{"pdf_url": pdfURL})
// }

// Function to append prescription text to PDF content
// func appendPrescriptionToPDF(pdfBytes []byte, prescription string) []byte {
// 	// Example implementation: appending prescription text to the end of PDF
// 	// This is a simplistic example and may need adjustment based on your PDF structure

// 	pdfStr := string(pdfBytes)

// 	// Append prescription text at the end of the PDF content
// 	updatedPDFStr := pdfStr + fmt.Sprintf("\nPrescription: %s\n", prescription)

// 	return []byte(updatedPDFStr)
// }

// func UploadXray(db *gorm.DB, c *gin.Context) {
// 	// Parse form data
// 	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
// 		return
// 	}

// 	// Retrieve form values
// 	patientID := c.PostForm("patient_id")
// 	doctorID := c.PostForm("doctor_id")
// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request", "details": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	// Read file into memory
// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded file", "details": err.Error()})
// 		return
// 	}

// 	// Save the file to a temporary location for YOLO inference
// 	tempFilePath := "/tmp/" + header.Filename
// 	if err := os.WriteFile(tempFilePath, fileBytes, 0644); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file", "details": err.Error()})
// 		return
// 	}
// 	defer os.Remove(tempFilePath)

// 	// Path to the YOLO command
// 	yoloCmd := "/Users/mamdouhhazem/opt/anaconda3/envs/dentex/bin/yolo"

// 	// Execute the YOLO command
// 	cmd := exec.Command(yoloCmd, "predict", "model=/Users/mamdouhhazem/Desktop/Graduaiton_Project/runs_results/newDatasetV8M/weights/best.pt", "source="+tempFilePath)
// 	predictedXray, err := cmd.CombinedOutput()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform inference", "details": err.Error(), "output": string(predictedXray)})
// 		return
// 	}

// 	// Find the latest predict folder
// 	runsDir := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/Dentify-X/cmd/runs/detect"
// 	latestPredictFolder, err := findLatestPredictFolder(runsDir)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find latest predict folder", "details": err.Error()})
// 		return
// 	}

// 	// Read the predicted image file
// 	predictedFilePath := latestPredictFolder + "/" + header.Filename // Adjust based on your YOLO output
// 	predictedFileBytes, err := os.ReadFile(predictedFilePath)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read predicted file", "details": err.Error()})
// 		return
// 	}

// 	patientIDUint, err := strconv.ParseUint(patientID, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID", "details": err.Error()})
// 		return
// 	}
// 	doctorIDUint, err := strconv.ParseUint(doctorID, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID", "details": err.Error()})
// 		return
// 	}

// 	medicalHistory := models.DoctorXray{
// 		DoctorID:      uint(doctorIDUint),
// 		PatientID:     uint(patientIDUint),
// 		PredictedXray: predictedFileBytes,
// 		XrayID:        fileBytes,
// 	}
// 	if err := db.Create(&medicalHistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
// }

// // Helper function to find the latest predict folder
// func findLatestPredictFolder(runsDir string) (string, error) {
// 	files, err := os.ReadDir(runsDir)
// 	if err != nil {
// 		return "", err
// 	}

// 	var latestFolder string
// 	var latestTime time.Time

// 	for _, file := range files {
// 		if file.IsDir() && strings.HasPrefix(file.Name(), "predict") {
// 			fileInfo, err := file.Info()
// 			if err != nil {
// 				return "", err
// 			}

// 			if fileInfo.ModTime().After(latestTime) {
// 				latestTime = fileInfo.ModTime()
// 				latestFolder = runsDir + "/" + file.Name()
// 			}
// 		}
// 	}

// 	if latestFolder == "" {
// 		return "", fmt.Errorf("no predict folder found")
// 	}

// 	return latestFolder, nil
// }

// func ViewPatientHistory(db *gorm.DB, c *gin.Context) {
// 	patientID := c.Param("pid")
// 	var medicalHistory models.DoctorXray

// 	if err := db.Select("xray_id, prescription, date").Where("patient_id = ?", patientID).Find(&medicalHistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medical history", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"medicalHistory": medicalHistory})
// 	// na2esha validation: patient is assigned to doctor or not.
// }

// func UploadXray(db *gorm.DB, c *gin.Context) {
// 	// Parse form data
// 	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form", "details": err.Error()})
// 		return
// 	}

// 	// Retrieve form values
// 	patientID := c.PostForm("patient_id")
// 	doctorID := c.PostForm("doctor_id")
// 	file, header, err := c.Request.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request", "details": err.Error()})
// 		return
// 	}
// 	defer file.Close()

// 	// Read file into memory
// 	fileBytes, err := io.ReadAll(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read uploaded file", "details": err.Error()})
// 		return
// 	}

// 	// Save the file to a temporary location for YOLO inference
// 	tempFilePath := "/tmp/" + header.Filename
// 	if err := os.WriteFile(tempFilePath, fileBytes, 0644); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file", "details": err.Error()})
// 		return
// 	}
// 	defer os.Remove(tempFilePath)

// 	// Path to the YOLO command
// 	yoloCmd := "/Users/mamdouhhazem/opt/anaconda3/envs/dentex/bin/yolo"

// 	// Execute the YOLO command
// 	cmd := exec.Command(yoloCmd, "predict", "model=/Users/mamdouhhazem/Desktop/Graduaiton_Project/runs_results/newDatasetV8M/weights/best.pt", "source="+tempFilePath)
// 	predictedXray, err := cmd.CombinedOutput()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform inference", "details": err.Error(), "output": string(predictedXray)})
// 		return
// 	}

// 	patientIDUint, err := strconv.ParseUint(patientID, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID", "details": err.Error()})
// 		return
// 	}
// 	doctorIDUint, err := strconv.ParseUint(doctorID, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID", "details": err.Error()})
// 		return
// 	}

// 	// predictedXrayBase64 := base64.StdEncoding.EncodeToString([]byte(predictedXray))

// 	medicalHistory := models.DoctorXray{
// 		DoctorID:      uint(doctorIDUint),
// 		PatientID:     uint(patientIDUint),
// 		PredictedXray: predictedXray,
// 		XrayID:        fileBytes,
// 	}
// 	if err := db.Create(&medicalHistory).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Data saved successfully"})
// }

// func UploadXray(c *gin.Context) {
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file from request"})
// 		return
// 	}

// 	ext := filepath.Ext(file.Filename)
// 	allowedExts := map[string]bool{".jpg": true, ".png": true, ".bmp": true, ".tiff": true}
// 	if !allowedExts[ext] {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, BMP, and TIFF files are allowed."})
// 		return
// 	}

// 	savePath := "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/Dentify-X/cmd/" + file.Filename
// 	if err := c.SaveUploadedFile(file, savePath); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
// 		return
// 	}

// 	// Path to the yolo command
// 	yoloCmd := "/Users/mamdouhhazem/opt/anaconda3/envs/dentex/bin/yolo"

// 	// Execute the yolo command with the correct model and source
// 	cmd := exec.Command(yoloCmd, "predict", "model=/Users/mamdouhhazem/Desktop/Graduaiton_Project/runs_results/newDatasetV8M/weights/best.pt", "source="+savePath)
// 	out, err := cmd.CombinedOutput()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform inference", "details": err.Error(), "output": string(out)})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"detections": string(out)})
// }

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
