package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Signal struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Repeaters   []string `json:"repeaters"`
}

type MarshallingHill struct {
	Tracks []Signal `json:"tracks"`
}

type Station struct {
	Name  string            `json:"name"`
	Hash  string            `json:"hash"`
	Hills []MarshallingHill `json:"hills"`
}

// App struct
type App struct {
	ctx      context.Context
	Stations map[string]Station
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.Stations = loadStationDefinitions(".")
	log.Print(a.Stations)
	a.ctx = ctx
}

func loadStationDefinitions(root string) map[string]Station {
	stations := map[string]Station{}
	err := filepath.Walk(root+"/stations", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.ReadFile(path)
		if err != nil {
			return errors.WithMessage(err, "cannot read file")
		}

		h := sha256.New()
		h.Write(file)

		bs := h.Sum(nil)

		var station Station
		err = json.Unmarshal(file, &station)
		if err != nil {
			return errors.WithMessage(err, "cannot unmarshal station")
		}

		station.Hash = fmt.Sprintf("%x", bs)

		stations[station.Name] = station
		return nil
	})

	if err != nil {
		log.Fatalf("cannot load station definitions: %s", err)
	}
	return stations
}

func (a *App) GetStations() map[string]Station {
	return a.Stations
}

func (a *App) GetStation(name string) Station {
	return a.Stations[name]
}
