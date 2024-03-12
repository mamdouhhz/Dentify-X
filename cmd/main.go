package main

import (
	"Dentify-X/app/models"
	"Dentify-X/app/routers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	db, err := models.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	fmt.Println("Database connection successful!")

	r := routers.Rout(db)

	// openssl ecparam -name prime256v1 -genkey -noout -out server.key
	// openssl req -x509 -new -key server.key -out server.crt -days 365
	err = http.ListenAndServeTLS(":443", "server.crt", "server.key", r)
	if err != nil {
		log.Fatal("Error running server:", err)
	}
}
