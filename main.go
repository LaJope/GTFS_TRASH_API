package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"google.golang.org/protobuf/proto"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig =
		&tls.Config{InsecureSkipVerify: true}

	http.HandleFunc("/stop", getStopVehiclesHaldler)
	http.HandleFunc("/aaaa", getStopVehiclesHaldler)
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func getStopFromGTFS(Stop_ID int) []string {
	var (
		username = "test"
		password = "test"
	)

	get_url := fmt.Sprintf("https://transport.orgp.spb.ru/"+
		"Portal/transport/internalapi/gtfs/realtime/"+
		"stopforecast?stopID=%d", Stop_ID)

	client := &http.Client{}
	req, err := http.NewRequest("GET", get_url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

  // for _, entity := range feed.Entity {
  //   fmt.Println(entity.GetTripUpdate())
  // }

  var vehicles []string


	for _, entity := range feed.Entity {
		tripUpdate := entity.GetTripUpdate()
		vehicle := tripUpdate.GetVehicle()
		vehicleId := vehicle.GetId()
    vehicles = append(vehicles, vehicleId)
	}
	//

	return vehicles
}

type Stop struct {
	ID int
}

// func getAAAAHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Only POST method is supported",
// 			http.StatusMethodNotAllowed)
// 		return
// 	}

	// var stop_info Stop
	// err := json.NewDecoder(r.Body).Decode(&stop_info)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// Вывод полученных данных в консоль
	// fmt.Printf("Received user: %+v\n", stop_info)

	// message := []string()

	// Отправка ответа клиенту
// 	w.Header().Set("Content-Type", "application/json")
// 	response := map[string][]string{"status": {"success"}, "message": message}
// 	_ = json.NewEncoder(w).Encode(response)
// }


func getStopVehiclesHaldler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is supported",
			http.StatusMethodNotAllowed)
		return
	}

	var stop_info Stop
	err := json.NewDecoder(r.Body).Decode(&stop_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received user: %+v\n", stop_info)

	message := getStopFromGTFS(stop_info.ID)


  w.Header().Set("Content-Type", "application/json")
	response := map[string][]string{"status": {"success"}, "message": message}
	_ = json.NewEncoder(w).Encode(response)

	// w.Header().Set("Content-Type", "application/json")
	// response := map[string]string{"status": "success", "message": message}
	// _ = json.NewEncoder(w).Encode(response)
}
