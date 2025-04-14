package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type FirebaseConfig struct {
	APIKey            string
	AuthDomain        string
	ProjectID         string
	StorageBucket     string
	MessagingSenderID string
	AppID             string
	MeasurementID     string
	DatabaseURL       string
}

type Paths struct {
	Count     string
	Bugs      string
	Debt      string
	Upgrades  string
	AscPoints string
}

type UpgradePaths struct {
	ClickPower    string
	DecreaseDelay string
	// AscMulti string
	// AutoRate string
}

type UpgradeSubpaths struct {
	Cost  string
	Name  string
	Level string
	Power string
	Type  string
}

type Database struct {
	conf            FirebaseConfig
	Paths           Paths
	UpgradePaths    UpgradePaths
	UpgradeSubpaths UpgradeSubpaths
}

type Upgrade struct {
	Cost  int    `json:"cost"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Power int    `json:"power"`
	Type  string `json:"type"`
}

func NewDatabase() Database {

	err := godotenv.Load()
	if err != nil {
		panic("Failed to load .env")
	}

	return Database{
		FirebaseConfig{
			APIKey:      os.Getenv("APIKey"),
			DatabaseURL: os.Getenv("DatabaseURL"),
		},
		Paths{
			Count:     "/rooms/global/game/count",
			AscPoints: "/rooms/global/game/ascPoints",
			Bugs:      "/rooms/global/game/bugs",
			Debt:      "/rooms/global/game/techDebt",
			Upgrades:  "/rooms/global/upgrades",
		},
		UpgradePaths{
			ClickPower:    "/clickPower",
			DecreaseDelay: "/decreaseDelay",
		},
		UpgradeSubpaths{
			Cost:  "/cost",
			Name:  "/name",
			Level: "/level",
			Power: "/power",
			Type:  "/type",
		},
	}
}

func (db *Database) FetchUpgrade(path string) (*Upgrade, error) {

	url := fmt.Sprintf("%s%s.json?auth=%s", db.conf.DatabaseURL, path, db.conf.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var upgrade Upgrade
	if err := json.Unmarshal(body, &upgrade); err != nil {
		return nil, err
	}

	return &upgrade, nil
}

func (db *Database) IncrementValueBy(path string, incrementBy int) error {

	currentValue := db.GetValue(path)

	newValue := currentValue + incrementBy

	return db.SetValue(path, newValue)
}

func (db *Database) SetValue(path string, value int) error {
	// Increment count
	newCountJSON, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling %s: %w", path, err)
	}

	// Send new count to db
	postURL := fmt.Sprintf("%s%s.json?auth=%s", db.conf.DatabaseURL, path, db.conf.APIKey)

	// Create request
	req, err := http.NewRequest(http.MethodPut, postURL, bytes.NewBufferString(string(newCountJSON)))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	// Do the thing
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error putting %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	log.Printf("Successfully updated %s to %d", path, value)
	return nil
}

func (db *Database) GetValue(path string) int {

	url := fmt.Sprintf("%s%s.json?auth=%s", db.conf.DatabaseURL, path, db.conf.APIKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Parse the JSON response
	var count int
	if err := json.Unmarshal(body, &count); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return count
}
