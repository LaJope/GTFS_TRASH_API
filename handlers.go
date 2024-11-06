package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func getStopsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported",
			http.StatusMethodNotAllowed)
		return
	}

	stopsInfo := parseCSVFiles("stops.txt", "stop_",
		"stop_id", "stop_name", "stop_lat", "stop_lon")

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stopsInfo)

}

func getStopInfoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported",
			http.StatusMethodNotAllowed)
		return
	}

	id_param := req.PathValue("id")
	if id_param == "" {
		http.Error(w, "To access this endpoint you need to provide ID in PATH",
			http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(id_param, 10, 0)
	if err != nil {
		http.Error(w, "Stop ID in PATH must be an integer",
			http.StatusBadRequest)
		return
	}

	stopInfo := getStopForecastRealtimeInfo(id)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(stopInfo)
}

func getVehicletripsInfoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported",
			http.StatusMethodNotAllowed)
		return
	}

	id_param := req.PathValue("id")
	if id_param == "" {
		http.Error(w, "To access this endpoint you need to provide ID in PATH",
			http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(id_param, 10, 0)
	if err != nil {
		http.Error(w, "Vehicle ID in PATH must be an integer",
			http.StatusBadRequest)
		return
	}

	vehicleInfo := getVehicleForecastRealtimeInfo(id)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(vehicleInfo)
}

func getVehiclePositionInfoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Only GET method is supported",
			http.StatusMethodNotAllowed)
		return
	}

	vehicleParams, _ := url.ParseQuery(req.URL.RawQuery)
	vehicleInfo := getVehiclePositionRealtimeInfo(vehicleParams)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(vehicleInfo)
}
