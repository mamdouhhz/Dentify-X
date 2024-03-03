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

	err = r.Run(":8080") // listen and serve on 0.0.0.0:8080
	if err != nil {
		return
	}
	log.Println("SERVER RUNNING ON 8080")
}
