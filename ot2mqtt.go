package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/beego/beego/v2/core/config"
	"github.com/coreos/go-systemd/journal"
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
	log.Println("MQTT Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("MQTT Connection lost: %v", err)
}
var logvars = map[string]string{
	"SYSLOG_IDENTIFIER": "ot2mqtt",
}

func main() {
	conf, cfgerr := config.NewConfig("ini", "app.conf")
	if cfgerr != nil {
		log.Fatal(cfgerr)
	}
	if !journal.Enabled() {
		log.Fatal("Needs a systemd journal")
	}
	journal.Send("ot2mqtt starting...", journal.PriInfo, logvars)

	broker_url, _ := conf.String("mqtt::url")
	broker_username, _ := conf.String("mqtt::username")
	broker_password, _ := conf.String("mqtt::password")

	opts := mqtt.NewClientOptions().AddBroker(broker_url)
	// opts.SetClientID("ot-pub")
	opts.SetUsername(broker_username)
	opts.SetPassword(broker_password)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/pub", pubHandler(client))

	listen_address, _ := conf.String("http::listen_address")
	journal.Send(fmt.Sprintf("http listener on %s", listen_address), journal.PriInfo, logvars)
	err := http.ListenAndServe(listen_address, mux)
	log.Fatal(err)
}

// func httpHandler(w http.ResponseWriter, r *http.Request) {
func pubHandler(client mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("Hi, from Service: " + name))

		fmt.Printf("%s %s %s\n", r.Method, r.URL, r.Proto)
		for k, v := range r.Header {
			fmt.Printf("Header[%q] = %q\n", k, v)
		}
		fmt.Printf("Host = %q\n", r.Host)
		fmt.Printf("RemoteAddr = %q\n", r.RemoteAddr)

		var ot_data OT_struct
		err := json.NewDecoder(r.Body).Decode(&ot_data)
		if err != nil {
			fmt.Printf("JSON decode error: %s\n", err)
			return
		}
		fmt.Printf("Lat: %f\n", ot_data.Lat)
		fmt.Printf("Lon: %f\n", ot_data.Lon)
		fmt.Printf("Topic: %s\n", ot_data.Topic)

		var jsonData []byte
		jsonData, err = json.Marshal(ot_data)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(jsonData))

		token := client.Publish(ot_data.Topic, 0, false, jsonData)
		token.Wait()

		fmt.Fprintf(w, "{}\n")

	}
}