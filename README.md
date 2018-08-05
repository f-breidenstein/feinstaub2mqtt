# feinstaub2mqtt

Dieses Tool kann man als "Custom API" in einem Feinstaubsensor von [http://luftdaten.info](luftdaten.info) eintragen um
die Sensorwerte mittels MQTT im Netzwerk zu verteilen. Damit lässt sich ein Feinstaubsensor bspw. einfacher in Home Assistant integrieren als über eine InfluxDB.


# Installation
Entweder mittels `git clone` oder `go get` den Sourcecode runterladen und selber kompilieren, oder unter [Releases](https://github.com/f-breidenstein/feinstaub2mqtt/releases) eine fertige Binary runterladen.

# MQTT Topics
Für jeden Sensorwert im übertragenen JSON, wird ein eigener MQTT Topic angelegt. Das format ist
```
PREFIX/<6stellige ID des verbauten ESP8266>/<Name des jeweiligen Sensors>
```
Den `PREFIX` kann man beim Starten mit `-prefix` ändern. Beispielhaft könten die erstellten Topics dann so aussehen:
```
ibg10/feinstaub/911386/SDS_P1 2.63
ibg10/feinstaub/911386/SDS_P2 2.43
ibg10/feinstaub/911386/BME280_temperature 36.79
ibg10/feinstaub/911386/BME280_humidity 22.66
ibg10/feinstaub/911386/BME280_pressure 100271.30
ibg10/feinstaub/911386/samples 581728
ibg10/feinstaub/911386/min_micro 243
ibg10/feinstaub/911386/max_micro 28265
ibg10/feinstaub/911386/signal -56
```

# Anwendung
Am einfachsten mit Hilfe eines SystemD Services auf einem kleinen Server laufen lassen. Hier ein beispielhaftes Service-File:
```
[Unit]
Description=Provide an API endpoint for a fine dust sensor from luftdaten.info

[Service]
Type=simple
ExecStart=/usr/bin/feinstaub2mqtt -mqttBroker tcp://broker.lol.nope:1883 
Restart=always

[Install]
WantedBy=multi-user.target

```
