package main

import (
	leagueapi "owl-stats/league-api"
)

func main() {
	// players, _ := leagueapi.GetAllPlayerData()
	match, _ := leagueapi.GetMatchData()
	leagueapi.PrintMatchData(match)
}
