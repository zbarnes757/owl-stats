package leagueapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

const statsURI = "https://api.overwatchleague.com/stats/players"

type data struct {
	Data []StatPlayer `json:"data"`
}

// StatPlayer struct
type StatPlayer struct {
	PlayerID                 int     `json:"playerId"`
	TeamID                   int     `json:"teamId"`
	Role                     string  `json:"role"`
	Name                     string  `json:"name"`
	Team                     string  `json:"team"`
	EliminationsAvgPer10M    float64 `json:"eliminations_avg_per_10m"`
	DeathsAvgPer10M          float64 `json:"deaths_avg_per_10m"`
	HeroDamageAvgPer10M      float64 `json:"hero_damage_avg_per_10m"`
	HealingAvgPer10M         float64 `json:"healing_avg_per_10m"`
	UltimatesEarnedAvgPer10M float64 `json:"ultimates_earned_avg_per_10m"`
	FinalBlowsAvgPer10M      float64 `json:"final_blows_avg_per_10m"`
	TimePlayedTotal          float64 `json:"time_played_total"`
}

func setMax(old float64, new float64) float64 {
	if new > old || old == 0 {
		return new
	}

	return old
}

func setMin(old float64, new float64) float64 {
	if new < old || old == 0 {
		return new
	}

	return old
}

func getMedian(list []float64) float64 {
	sort.Float64s(list)
	mid := len(list) / 2
	return list[mid]
}

// GetAllPlayerData get aggregated data for all players
func GetAllPlayerData() ([]StatPlayer, error) {
	resp, err := http.Get(statsURI)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer resp.Body.Close()

	var d data
	json.NewDecoder(resp.Body).Decode(&d)

	// total := float64(len(d.Data))

	var (
		minHeroDamage,
		minEliminations,
		minDeaths,
		minHealing,
		minUltimates,
		minFinalBlows,
		minTime,
		maxHeroDamage,
		maxEliminations,
		maxDeaths,
		maxHealing,
		maxUltimates,
		maxFinalBlows,
		maxTime float64
		heroDamage,
		eliminations,
		deaths,
		healing,
		ultimates,
		finalBlows,
		times []float64
	)

	for _, player := range d.Data {
		minHeroDamage = setMin(minHeroDamage, player.HeroDamageAvgPer10M)
		minEliminations = setMin(minEliminations, player.EliminationsAvgPer10M)
		minDeaths = setMin(minDeaths, player.DeathsAvgPer10M)
		minHealing = setMin(minHealing, player.HealingAvgPer10M)
		minUltimates = setMin(minUltimates, player.UltimatesEarnedAvgPer10M)
		minFinalBlows = setMin(minFinalBlows, player.FinalBlowsAvgPer10M)
		minTime = setMin(minTime, player.TimePlayedTotal)

		maxHeroDamage = setMax(maxHeroDamage, player.HeroDamageAvgPer10M)
		maxEliminations = setMax(maxEliminations, player.EliminationsAvgPer10M)
		maxDeaths = setMax(maxDeaths, player.DeathsAvgPer10M)
		maxHealing = setMax(maxHealing, player.HealingAvgPer10M)
		maxUltimates = setMax(maxUltimates, player.UltimatesEarnedAvgPer10M)
		maxFinalBlows = setMax(maxFinalBlows, player.FinalBlowsAvgPer10M)
		maxTime = setMax(maxTime, player.TimePlayedTotal)

		heroDamage = append(heroDamage, player.HeroDamageAvgPer10M)
		eliminations = append(eliminations, player.EliminationsAvgPer10M)
		deaths = append(deaths, player.DeathsAvgPer10M)
		healing = append(healing, player.HealingAvgPer10M)
		ultimates = append(ultimates, player.UltimatesEarnedAvgPer10M)
		finalBlows = append(finalBlows, player.FinalBlowsAvgPer10M)
		times = append(times, player.TimePlayedTotal)
	}

	// fmt.Printf(
	// 	"Total Players: %.f\n\n"+
	// 		"Max Damage: %.2f\n"+
	// 		"Max Eliminations: %.2f\n"+
	// 		"Max Deaths: %.2f\n"+
	// 		"Max Healing: %.2f\n"+
	// 		"Max Ultimates Earned: %.2f\n"+
	// 		"Max Final Blows: %.2f\n"+
	// 		"Max Time Played: %.2f\n\n"+
	// 		"Min Damage: %.2f\n"+
	// 		"Min Eliminations: %.2f\n"+
	// 		"Min Deaths: %.2f\n"+
	// 		"Min Healing: %.2f\n"+
	// 		"Min Ultimates Earned: %.2f\n"+
	// 		"Min Final Blows: %.2f\n"+
	// 		"Min Time Played: %.2f\n\n"+
	// 		"Median Damage: %.2f\n"+
	// 		"Median Eliminations: %.2f\n"+
	// 		"Median Deaths: %.2f\n"+
	// 		"Median Healing: %.2f\n"+
	// 		"Median Ultimates Earned: %.2f\n"+
	// 		"Median Final Blows: %.2f\n"+
	// 		"Median Time Played: %.2f\n",
	// 	total,
	// 	maxHeroDamage,
	// 	maxEliminations,
	// 	maxDeaths,
	// 	maxHealing,
	// 	maxUltimates,
	// 	maxFinalBlows,
	// 	maxTime,
	// 	minHeroDamage,
	// 	minEliminations,
	// 	minDeaths,
	// 	minHealing,
	// 	minUltimates,
	// 	minFinalBlows,
	// 	minTime,
	// 	getMedian(heroDamage),
	// 	getMedian(eliminations),
	// 	getMedian(deaths),
	// 	getMedian(healing),
	// 	getMedian(ultimates),
	// 	getMedian(finalBlows),
	// 	getMedian(times),
	// )

	return d.Data, nil
}
