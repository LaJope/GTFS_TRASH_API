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

func getStopForecastRealtimeInfo(stopID int64) []map[string]interface{} {
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

	stopInfo := make([]map[string]interface{}, 0)

	for _, entity := range feed.Entity {
		entityInfo := make(map[string]interface{})

		tripUpdate := entity.GetTripUpdate()
		trip := tripUpdate.GetTrip()
		vehicle := tripUpdate.GetVehicle()
		stopTimeUpdate := tripUpdate.GetStopTimeUpdate()

		entityInfo["route_id"], _ =
			strconv.ParseInt(trip.GetRouteId(), 10, 0)
		entityInfo["vehicle_id"], _ =
			strconv.ParseInt(vehicle.GetId(), 10, 0)
		entityInfo["arrival"] =
			time.Unix(stopTimeUpdate[0].GetArrival().GetTime(), 0)

		stopInfo = append(stopInfo, entityInfo)
	}

	return stopInfo
}

func getVehicleForecastRealtimeInfo(vehicleID int64) []map[string]interface{} {
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

	vehicleInfo := make([]map[string]interface{}, 0)

	for _, entity := range feed.Entity {
		entityInfo := make(map[string]interface{})

		tripUpdate := entity.GetTripUpdate()
		stopTimeUpdate := tripUpdate.GetStopTimeUpdate()

		entityInfo["id"], _ = strconv.ParseInt(entity.GetId(), 10, 0)
		entityInfo["forecast"] = make([]map[string]interface{},
			len(stopTimeUpdate))

		stopsInfo := entityInfo["forecast"].([]map[string]interface{})

		for ind, stop := range stopTimeUpdate {
			stopsInfo[ind] = make(map[string]interface{})
			stopsInfo[ind]["stop_id"], _ =
				strconv.ParseInt(stop.GetStopId(), 10, 0)
			stopsInfo[ind]["arrival"] =
				time.Unix(stopTimeUpdate[0].GetArrival().GetTime(), 0)
		}

		vehicleInfo = append(vehicleInfo, entityInfo)
	}

	return vehicleInfo
}

func getVehiclePositionRealtimeInfo(params url.Values) []map[string]interface{} {

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

	vehiclePositionInfo := make([]map[string]interface{}, 0)

	for _, entity := range feed.Entity {
		entityInfo := make(map[string]interface{})

		vehicle := entity.GetVehicle()
		position := vehicle.GetPosition()

		entityInfo["vehicle_id"], _ = strconv.ParseInt(entity.GetId(), 10, 0)
		entityInfo["route_id"], _ =
			strconv.ParseInt(vehicle.GetTrip().GetRouteId(), 10, 0)
		entityInfo["lat"] = position.GetLatitude()
		entityInfo["lon"] = position.GetLongitude()
		entityInfo["bearing"] = position.GetBearing()

		vehiclePositionInfo = append(vehiclePositionInfo, entityInfo)
	}

	return vehiclePositionInfo
}
