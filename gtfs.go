package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

const (
	GTFS_URL = "https://transport.orgp.spb.ru/Portal/transport/internalapi/gtfs"
)

func getFeedZipArchive() []byte {
	url := fmt.Sprintf("%s/feed.zip", GTFS_URL)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getStopForecastRealtimeInfo(stopID int64) []stopRealtimeInfo {
	url := fmt.Sprintf("%s/realtime/stopforecast?stopID=%d", GTFS_URL, stopID)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	stopInfo := make([]stopRealtimeInfo, 0)

	for _, entity := range feed.Entity {
		var entityInfo stopRealtimeInfo

		tripUpdate := entity.GetTripUpdate()
		trip := tripUpdate.GetTrip()
		vehicle := tripUpdate.GetVehicle()
		stopTimeUpdate := tripUpdate.GetStopTimeUpdate()

		entityInfo.Route_id, _ =
			strconv.ParseInt(trip.GetRouteId(), 10, 0)
		entityInfo.Vehicle_id, _ =
			strconv.ParseInt(vehicle.GetId(), 10, 0)
		stopTime := time.Unix(stopTimeUpdate[0].GetArrival().GetTime(), 0)
		entityInfo.Arrival = stopTime.Format(time.UnixDate)

		stopInfo = append(stopInfo, entityInfo)
	}

	return stopInfo
}

func getVehicleForecastRealtimeInfo(vehicleID int64) []vehicleRealtimeInfo {
	url := fmt.Sprintf("%s/realtime/vehicletrips?vehicleIDs=%d",
		GTFS_URL, vehicleID)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	vehicleInfo := make([]vehicleRealtimeInfo, 0)

	for _, entity := range feed.Entity {
		var entityInfo vehicleRealtimeInfo

		tripUpdate := entity.GetTripUpdate()
		stopTimeUpdate := tripUpdate.GetStopTimeUpdate()

		entityInfo.Id, _ = strconv.ParseInt(entity.GetId(), 10, 0)
		entityInfo.Forecast = make([]vehicleStopRealtimeInfo,
			len(stopTimeUpdate))

		for ind, stop := range stopTimeUpdate {
			entityInfo.Forecast[ind].StopId, _ =
				strconv.ParseInt(stop.GetStopId(), 10, 0)
			stopTime := time.Unix(stopTimeUpdate[0].GetArrival().GetTime(), 0)
			entityInfo.Forecast[ind].Arrival = stopTime.Format(time.UnixDate)
		}

		vehicleInfo = append(vehicleInfo, entityInfo)
	}

	return vehicleInfo
}

func getVehiclePositionRealtimeInfo(params url.Values) []vehiclePositionRealtimeInfo {

	reqParams := make([]string, 0)
	if bbox, ok := params["bbox"]; ok {
		reqParams = append(reqParams, fmt.Sprintf("bbox=%s", bbox[0]))
	}
	if transports, ok := params["transports"]; ok {
		reqParams = append(reqParams, fmt.Sprintf("transports=%s", transports[0]))
	}
	if routeIDs, ok := params["routeIDs"]; ok {
		reqParams = append(reqParams, fmt.Sprintf("routeIDs=%s", routeIDs[0]))
	}

	url := fmt.Sprintf("%s/realtime/vehicle?%s", GTFS_URL,
		strings.Join(reqParams, "&"))
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	vehiclePositionInfo := make([]vehiclePositionRealtimeInfo, 0)

	for _, entity := range feed.Entity {
		var entityInfo vehiclePositionRealtimeInfo

		vehicle := entity.GetVehicle()
		position := vehicle.GetPosition()

		entityInfo.Id, _ = strconv.ParseInt(entity.GetId(), 10, 0)
		entityInfo.Route_id, _ =
			strconv.ParseInt(vehicle.GetTrip().GetRouteId(), 10, 0)
		entityInfo.Lat = position.GetLatitude()
		entityInfo.Lon = position.GetLongitude()
		entityInfo.Bearing = position.GetBearing()

		vehiclePositionInfo = append(vehiclePositionInfo, entityInfo)
	}

	return vehiclePositionInfo
}
