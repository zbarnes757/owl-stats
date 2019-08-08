package app

import (
	"log"
	"owl-stats/models"
	"time"
)

// StartScheduledProcesses will kick off all background jobs
func StartScheduledProcesses() {
	log.Println("Kicking off all scheduled processes...")
	go schedulePlayerData()
}

func schedulePlayerData() {
	log.Println("Player Data fetch process started...")
	var nextRun *time.Time
	var polling *models.Polling
	for {
		if nextRun == nil {
			polling = models.GetPolling("player_fetch")

			if polling == nil || polling.NextRun == nil {
				t := time.Now().Add(time.Second * 10)
				nextRun = &t
			} else {
				nextRun = polling.NextRun
			}
		}

		// sleep until it's the next time
		time.Sleep(nextRun.Sub(time.Now()))
		log.Println("Fetching player data...")

		lastRan := time.Now()
		t := time.Now().Add(time.Hour * 24 * 7)
		nextRun = &t
		if polling != nil {
			log.Println("Updating player_fetch polling to run again at: " + nextRun.String())
			models.GetDB().Model(polling).Updates(models.Polling{LastRan: &lastRan, NextRun: nextRun})
		}
	}
}
