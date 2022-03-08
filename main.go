package main

import (
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type OT_struct struct {
	Type  string  `json:"_type"`
	BSSID string  `json:"BSSID"`
	SSID  string  `json:"SSID"`
	T     string  `json:"t"`
	Batt  int     `json:"batt"`
	Lat   float32 `json:"lat"`
	Lon   float32 `json:"lon"`
	Acc   int     `json:"acc"`
	Alt   int     `json:"alt"`
	Vac   int     `json:"vac"`
	Vel   int     `json:"vel"`
	Tid   string  `json:"tid"`
	Tst   int32   `json:"tst"`
	Topic string  `json:"topic"`
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Print("INFO: Connected to MQTT broker")

}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("ERROR: MQTT Connection lost: %v", err)
}

func main() {
	log.Print("INFO: ot2mqtt starting...")

	broker_url := os.Getenv("BROKER_URL")
	broker_username := os.Getenv("BROKER_USERNAME")
	broker_password := os.Getenv("BROKER_PASSWORD")
	listen_address := os.Getenv("LISTEN_ADDRESS")

	opts := mqtt.NewClientOptions().AddBroker(broker_url)
	// opts.SetClientID("ot-pub")
	opts.SetUsername(broker_username)
	opts.SetPassword(broker_password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(time.Second * 5)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pub", httpHandler(client))

	log.Printf("INFO: http listener on %s", listen_address)

	err := http.ListenAndServe(listen_address, mux)
	log.Fatal(err)
}
