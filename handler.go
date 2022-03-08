package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func httpHandler(client mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var ot_data OT_struct
		err := json.NewDecoder(r.Body).Decode(&ot_data)
		if err != nil {
			log.Printf("ERROR: JSON decode error: %s\n", err)
			return
		}
		log.Printf("Topic: %s, Lat: %f, Lon: %f\n", ot_data.Topic, ot_data.Lat, ot_data.Lon)

		var jsonData []byte
		jsonData, err = json.Marshal(ot_data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(jsonData))

		token := client.Publish(ot_data.Topic, 0, false, jsonData)
		token.Wait()

		fmt.Fprintf(w, "[]\n")

	}
}
