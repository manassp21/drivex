package main

import(
	"net/http"
	"fmt"
	"log"

	"github.com/manassp21/drivex/routes"
)

func main(){

	fmt.Println("Loading config...")

	fmt.Println("Connecting to Database...")
	
	r := routes.SetupRoutes()

	port := ":8080"
	fmt.Println("Server starting on port", port)

	err := http.ListenAndServe(port, r)
	if(err!=nil){
		log.Fatal("Server failed to start...")
	}
}