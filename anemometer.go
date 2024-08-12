package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/tarm/serial"
)

type AnemometerArray struct {
	sync.RWMutex

	Meters []string // Associated
}

type Anemometer struct {
	DeviceNumber int

	WindSpeed   float64
	Temperature float64
	Humidity    float64
}

func (aa *AnemometerArray) Scan(anemometerChannel chan<- Anemometer) {
	for {
		devices, _ := findUSBSerialDevices()
		for _, device := range devices {

			// If not found, add it
			if aa.newDevice(device) {
				log.Printf("Found device %s and adding it", device)

				aa.Lock()

				aa.Meters = append(aa.Meters, device)
				go func(device string, index int) {
					for {
						readSerialDevice(device, index, anemometerChannel)
						time.Sleep(time.Second)
					}
				}(device, len(aa.Meters))

				aa.Unlock()
			}
		}

		time.Sleep(time.Second)
	}
}

func (aa *AnemometerArray) newDevice(device string) bool {
	aa.RLock()
	defer aa.RUnlock()

	for _, aMeter := range aa.Meters {
		if aMeter == device {
			return false
		}
	}
	return true
}

var USBserialRE = regexp.MustCompile(`^(tty\.usbserial-.*)|(ttyUSB.*)$`)

func findUSBSerialDevices() ([]string, error) {
	var devices []string

	files, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if USBserialRE.MatchString(file.Name()) {
			devices = append(devices, "/dev/"+file.Name())
		}
	}

	return devices, nil
}

func readSerialDevice(device string, index int, anemometerChannel chan<- Anemometer) {

	// Open the serial port
	port, err := serial.OpenPort(&serial.Config{
		Name: device,
		Baud: 9600,
	})
	if err != nil {
		log.Println("Error opening serial port:", err)
		return
	}
	defer port.Close()

	// Read data from the serial port
	buf := make([]byte, 128)
	//prevTime := time.Now()
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Println("Error reading from serial port:", err)
			return
		}

		if n == 18 {
			//Beaufort := int(buf[8])
			windSpeed := binary.LittleEndian.Uint16(buf[5:7])
			humidity := binary.LittleEndian.Uint16(buf[9:11])
			temperatur := binary.LittleEndian.Uint16(buf[13:15])

			data := Anemometer{
				DeviceNumber: index,

				WindSpeed:   float64(windSpeed) / 100,
				Temperature: float64(temperatur) / 100,
				Humidity:    float64(humidity) / 10,
			}
			anemometerChannel <- data

			//log.Printf("#%d, %s: Received %d bytes: %s. Wind: %.2f (%d). Temp: %.2f. Humidity: %.2f. deltaT: %v\n", index, device, n, printHex(buf[:n]), float64(windSpeed)/100, Beaufort, float64(temperatur)/100, float64(humidity)/10, time.Since(prevTime))
			//prevTime = time.Now()
		}
	}
}

func printHex(theBytes []byte) string {
	output := ""
	// Iterate over the bytes
	for i, b := range theBytes {
		// Print two hexadecimal characters for the byte
		output += fmt.Sprintf("%02X ", b)

		// Add an extra space after every 8 bytes (except the first 8 bytes)
		if (i+1)%8 == 0 && i != 0 {
			output += " "
		}
	}
	return output
}
