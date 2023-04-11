package main

import (
	"log"
	"pulse/config"
	"pulse/database"
	"pulse/efx"
	"pulse/internal"
	"pulse/structs"
	"time"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	//Connect to the database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Create a ticker that ticks every 5 seconds

	ticker := time.NewTicker(time.Second * time.Duration(cfg.Interval))
	defer ticker.Stop()

	// Create a channel for spread records
	spreadRecordsCh := make(chan *structs.SpreadRecord)

	// Create a new rest connector
	connector := efx.NewRestConnector(cfg.APIKey, cfg.APISecret, false)

	// Start a goroutine to periodically collect spread records
	go func() {
		for {
			// Wait for the ticker to tick
			<-ticker.C

			// Get the positions
			positions, err := connector.GetPositions()
			if err != nil {
				log.Printf("Failed to get positions: %v", err)
				continue
			}

			// Group the positions by symbol
			orderBookMap := structs.GroupOrdersBySymbol(positions)

			// Calculate the spreads for each level in the configuration
			scrapper := internal.SpreadScrapper{
				OrderBooks: orderBookMap,
				Config:     cfg.SpreadConfig,
			}
			scrapper.Run(spreadRecordsCh)
		}
	}()

	// Start a goroutine to save spread records to the database
	go func() {
		for spreadRecord := range spreadRecordsCh {

			err := db.SaveSpreadRecord(spreadRecord)
			if err != nil {
				log.Printf("Failed to save spread record to database: %v", err)
			} else {
				log.Printf("Saved spread record to database: %+v\n", spreadRecord)
			}
		}
	}()

	// Block forever
	select {}
}
