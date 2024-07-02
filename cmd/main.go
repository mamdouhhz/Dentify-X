package main

import (
	"Dentify-X/app/models"
	"Dentify-X/app/routers"
	"fmt"
	"log"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	fmt.Println("Database connection successful!")

	r := routers.Rout(db)

	// Serve static files from the specified directory
	r.Static("/files", "/Users/mamdouhhazem/Desktop/Graduaiton_Project/Project_II/AI_Enabled_Dental_Diagnostic_Tool_project2/Dentify-X/Project_Grad/htmlandcssandimages")

	// Run the HTTPS server
	// openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt -config localhost.cnf
	err = r.RunTLS(":443", "server.crt", "server.key")
	if err != nil {
		log.Fatal("Error running server:", err)
	}
}
