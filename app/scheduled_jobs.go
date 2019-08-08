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
	var polling *models.Polling
	for {
		if polling == nil {
			polling = models.GetPolling("player_fetch")

			if polling.NextRun == (time.Time{}) {
				polling.NextRun = time.Now().Add(time.Second * 10)
			}
		}

		// sleep until it's the next time
		time.Sleep(polling.NextRun.Sub(time.Now()))
		log.Println("Fetching player data...")

		next := time.Now().Add(time.Hour * 24 * 7)
		models.GetDB().Model(polling).Updates(models.Polling{LastRan: time.Now(), NextRun: next})
		log.Println("Updating player_fetch polling to run again at: " + polling.NextRun.String())
	}
}
