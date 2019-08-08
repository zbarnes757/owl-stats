package app

import (
	"log"
	leagueapi "owl-stats/league-api"
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
		log.Println("Scheduling player_fetch polling to run at: " + polling.NextRun.UTC().Format("2006-01-02T15:04:05"))
		time.Sleep(polling.NextRun.Sub(time.Now()))
		players := leagueapi.FetchAllPlayers()

		models.BulkCreateOWLPlayers(players)

		next := time.Now().Add(time.Hour * 24 * 7)
		log.Println("Updating player_fetch polling to run again at: " + next.UTC().Format("2006-01-02T15:04:05"))
		polling.Update(models.Polling{LastRan: time.Now(), NextRun: next})
	}
}
