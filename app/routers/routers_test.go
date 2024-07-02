package routers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Dentify-X/app/models"
	"Dentify-X/app/routers"

	"github.com/stretchr/testify/assert"
)

func TestDoctorSignupRequest(t *testing.T) {
	db, err := models.InitDB()
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}

	r := routers.Rout(db)

	// Define a JSON payload for the request body if required
	requestBody := map[string]string{
		"name":         "ahmed",
		"gender":       "male",
		"mln":          "isdfkll",
		"phone_number": "01097277",
		"email":        "ahmed@gmail.com",
		"password":     "123",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling JSON request: %v", err)
	}

	// Create a test request to /dsignupreq with JSON body
	req, err := http.NewRequest("POST", "/dsignupreq", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request to the handler
	r.ServeHTTP(rr, req)

	expectedResponseBody := `{"error":"check your email, fail"}`
	assert.Equal(t, expectedResponseBody, rr.Body.String(), "response body not as expected")
}

func TestPatientLogin(t *testing.T) {
	db, err := models.InitDB()
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}

	r := routers.Rout(db)

	// Define a JSON payload for the request body if required
	requestBody := map[string]string{
		"email":    "ahmed@gmail.com",
		"password": "123",
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling JSON request: %v", err)
	}

	// Create a test request to /dsignupreq with JSON body
	req, err := http.NewRequest("POST", "/plogin", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Serve the HTTP request to the handler
	r.ServeHTTP(rr, req)

	expectedResponseBody := `{"error":"you are not signed up"}`
	assert.Equal(t, expectedResponseBody, rr.Body.String(), "response body not as expected")
}
