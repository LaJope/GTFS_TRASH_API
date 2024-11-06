package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig =
		&tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/api/stops", getStopsHandler)
	http.HandleFunc("/api/stops/{id}", getStopInfoHandler)
	http.HandleFunc("/api/vehicle/{id}", getVehicletripsInfoHandler)
	http.HandleFunc("/api/position", getVehiclePositionInfoHandler)

	http.HandleFunc("/api/downzip", getGTFSFeed)

	fmt.Println("Starting server on port 3334...")
	if err := http.ListenAndServe(":3334", nil); err != nil {
		log.Fatal(err)
	}
}
