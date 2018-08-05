package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Sensor struct {
	Type  string `json:"value_type"`
	Value string `json:"value"`
}

type Payload struct {
	ID              string   `json:"esp8266id"`
	SoftwareVersion string   `json:"software_version"`
	Sensors         []Sensor `json:"sensordatavalues"`
}

func handler(w http.ResponseWriter, r *http.Request, c MQTT.Client, prefix string) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)

	var p Payload
	err := json.Unmarshal(buf.Bytes(), &p)
	if err != nil {
		log.Print("Can't unmarshal payload into struct")
		log.Fatal(err)
	}
	log.Printf("Got data from %q\n", r.RemoteAddr)
	for _, sensor := range p.Sensors {
		t := fmt.Sprintf("%s/%s/%s", prefix, p.ID, sensor.Type)
		log.Printf("Publishing to %q\n", t)
		if token := c.Publish(t, 0, false, sensor.Value); token.Wait() && token.Error() != nil {
			log.Fatal(token.Error())
		}
	}
}

func main() {
	listenAddr := flag.String("listen", ":8080", "Address to listen on ")
	mqttBroker := flag.String("broker", "", "Address to the mqtt broker (Required)")
	mqttClientID := flag.String("clientID", "feinstaub", "ClientID for this tool")
	mqttUser := flag.String("mqttUser", "", "Username for MQTT")
	mqttPass := flag.String("mqttPass", "", "Passwort for MQTT")
	topicPrefix := flag.String("prefix", "feinstaub", "Prefix for the sensordata")
	flag.Parse()

	if *mqttBroker == "" {
		log.Fatal("No MQTT broker was specified! Exiting")
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(*mqttBroker)
	opts.SetClientID(*mqttClientID)
	if *mqttUser != "" && *mqttPass != "" {
		opts.SetUsername(*mqttUser)
		opts.SetPassword(*mqttPass)
	}
	client := MQTT.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Connected to MQTT broker")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, client, *topicPrefix)
	})
	http.ListenAndServe(*listenAddr, nil)
}
