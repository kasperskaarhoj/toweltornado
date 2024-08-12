package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"
)

func LoadHiScores(filename string) ([]*HiScoreEntry, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return []*HiScoreEntry{}, err
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return []*HiScoreEntry{}, err
	}

	// Unmarshal the JSON data
	var hiScores []*HiScoreEntry
	err = json.Unmarshal(data, &hiScores)
	if err != nil {
		return []*HiScoreEntry{}, err
	}

	return hiScores, nil
}

func SaveHiScores(filename string, hiScores []*HiScoreEntry) error {
	// Marshal the slice into JSON
	data, err := json.MarshalIndent(hiScores, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddToHiscore(hiScore *HiScoreEntry) []*HiScoreEntry {

	// Get existing:
	hiScores, _ := LoadHiScores("hiscores.json")

	// Search for record and update
	found := false
	for i, hiSc := range hiScores {
		if hiSc.Name == hiScore.Name {
			hiScores[i] = hiScore
			found = true
			break
		}
	}

	// Add if not found:
	if !found {
		hiScores = append(hiScores, hiScore)
	}

	// Sort the slice based on Time field
	sort.Slice(hiScores, func(i, j int) bool {
		return hiScores[i].Time < hiScores[j].Time
	})

	SaveHiScores("hiscores.json", hiScores)

	return hiScores
}
