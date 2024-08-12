package main

import (
	"net/http"

	utils "github.com/kasperskaarhoj/toweltornado/utils"
	log "github.com/s00500/env_logger"
)

//go:generate sh injectGitVars.sh

func main() {

	tt := &TowelTornado{}
	frontEndPort := 8031 // Web Frontend Port

	// Create frontend server:
	frontendServer := &http.Server{}
	wsserver := &wsServer{}
	go func() {
		frontendServer = wsserver.createServer(frontEndPort, tt)
		if frontendServer != nil {
			IP, _ := utils.FetchOutboundIP()
			log.Infof("Frontend is available on http://%s:%d\n", IP, frontEndPort)
			frontendServer.ListenAndServe()
		}
	}()

	// Array of anemometers and scanning for them, adding them when found and setting up readings posted into the channel:
	Anemometers := AnemometerArray{}
	anemometerChannel := make(chan Anemometer, 100)
	go Anemometers.Scan(anemometerChannel)

	// Reading the channel:
	for data := range anemometerChannel {
		wsserver.BroadcastMessage(&wsToClient{DeviceNumber: data.DeviceNumber, WindSpeed: &data.WindSpeed})
		//log.Printf("#%d: Wind: %.2f. Temp: %.2f. Humidity: %.2f\n", data.DeviceNumber, data.WindSpeed, data.Temperature, data.Humidity)
	}
}
