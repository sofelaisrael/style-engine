package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sofelaisrael/style-engine/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/transform", api.HandleTransform)
	http.HandleFunc("/styles", api.HandleListStyles)
	http.HandleFunc("/", api.HandleIndex)

	fmt.Printf("Style Engine running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}