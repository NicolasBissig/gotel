package main

import (
	"fmt"
	"github.com/NicolasBissig/gotel"
	"github.com/NicolasBissig/gotel/gotelhttp"
	"io"
	"log"
	"math/rand"
	"net/http"
)

func init() {
	_, err := gotel.Setup()
	if err != nil {
		log.Fatalf("failed to setup gotel: %v", err)
	}
}

func main() {
	gotelhttp.HandleFunc("GET /rolldice", func(w http.ResponseWriter, r *http.Request) {
		value := rand.Int()%6 + 1

		_, err := io.WriteString(w, fmt.Sprintf("%d\n", value))
		if err != nil {
			log.Printf("Write failed: %v\n", err)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
