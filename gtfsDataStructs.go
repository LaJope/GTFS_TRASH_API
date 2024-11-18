package main

type stopRealtimeInfo struct {
	Route_id   int64  `json:"route_id"`
	Vehicle_id int64  `json:"vehicle_id"`
	Arrival    string `json:"arrival"`
}

type vehicleStopRealtimeInfo struct {
	StopId  int64  `json:"stop_id"`
	Arrival string `json:"arrival"`
}

type vehicleRealtimeInfo struct {
	Id       int64                     `json:"id"`
	Forecast []vehicleStopRealtimeInfo `json:"forecast"`
}

type vehiclePositionRealtimeInfo struct {
	Id       int64   `json:"id"`
	Route_id int64   `json:"route_id"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
	Bearing  float32 `json:"bearing"`
}
