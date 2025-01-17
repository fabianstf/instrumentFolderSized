package main

import (
	"instrumentFolderSized/handlers"
	"log"
	"net/http"
)

func main() {
	router := http.ServeMux{}
	router.HandleFunc("/", handlers.IndexHandler)
	router.HandleFunc("/createPDF", handlers.CreatePdfHandler)

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	port := ":8080"
	log.Printf("Server started at http://localhost%s", port)
	err := http.ListenAndServe(port, &router)
	if err != nil {
		log.Fatal(err)
	}
}
