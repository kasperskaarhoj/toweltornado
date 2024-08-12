package utils

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"os"
	"strings"

	skconfig "github.com/SKAARHOJ/ibeam-lib-config"
	log "github.com/s00500/env_logger"
)

// updates the Json
func UpdateJson(data interface{}, fileName string) {
	file, err := json.MarshalIndent(data, "", "\t")
	if !log.Should(err) {
		err = os.WriteFile(fileName, file, 0755)
		log.Should(err)
	}
}

// Reads the JSON file into data and creates a new file if not found
func ReadJson(data interface{}, fileName string) {
	file, err := os.ReadFile(fileName) // Read File

	// Write a new file
	if log.Should(err) {
		jsonData, err := json.Marshal(data)
		if log.Should(err) {
			return
		}

		err = os.WriteFile(fileName, jsonData, 0644)
		log.Should(err)
		return
	}

	err = json.Unmarshal(file, data)
	log.Should(err)
}

// Create a folder on disk if it does not already exsist
func Makedir(folder string) string {

	folderPath := fmt.Sprintf("%s/%s", skconfig.GetConfigPath(), folder)

	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(folderPath, 0755)
		if errDir != nil {
			log.Should(err)
		}
	}
	return folderPath
}

// Get preferred outbound ip of this machine
func FetchOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, log.Wrap(err, "on fetching outbound IP address")
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func TrimExplode(str string, token string) []string {
	output := []string{}
	splitString := strings.Split(str, token)
	for _, s := range splitString {
		trimmed := strings.TrimSpace(s)
		if trimmed != "" {
			output = append(output, trimmed)
		}
	}

	return output
}

func MapValueCeil(x int, in_min int, in_max int, out_min int, out_max int) int {
	floatDiv := float64((x-in_min)*(out_max-out_min)) / float64((in_max - in_min))

	return int(math.Ceil(floatDiv)) + out_min
}

func ConstrainValue(input int, minimumValue int, maximumValue int) int {
	if input < minimumValue {
		return minimumValue
	} else if input > maximumValue {
		return maximumValue
	} else {
		return input
	}
}

func MapAndConstrainValueCeil(x int, in_min int, in_max int, out_min int, out_max int) int {
	return ConstrainValue(MapValueCeil(x, in_min, in_max, out_min, out_max), out_min, out_max)
}

//go:embed resources
var wsFS embed.FS

// Read contents from embedded file
func ReadResourceFile(fileName string) []byte {
	fileName = strings.ReplaceAll(fileName, "\\", "/")
	byteValue, err := wsFS.ReadFile(fileName)
	log.Should(err)
	return byteValue
}
