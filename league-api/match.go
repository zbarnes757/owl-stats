package leagueapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const matchURI = "https://api.overwatchleague.com/stats/matches/10528/maps/1"

// GetMatchData get all information for a match
func GetMatchData() (Match, error) {
	newResp, err := http.Get(matchURI)
	if err != nil {
		fmt.Print(err)
		return Match{}, err
	}
	defer newResp.Body.Close()

	var m Match
	json.NewDecoder(newResp.Body).Decode(&m)

	return m, nil
}

// PrintMatchData to print all the data
func PrintMatchData(m Match) {
	// playerData := make(map[int]StatPlayer)

	// for _, player := range players {
	// 	playerData[player.PlayerID] = player
	// }

	for _, team := range m.Teams {
		for _, player := range team.Players {
			p, err := fetchPlayer(player.ESportsPlayerID)
			if err == nil {
				printPlayerData(p, player)
			}
		}
	}

}

func printPlayerData(player Player, mP MatchPlayer) {
	fmt.Printf("%s %s:\n", player.GivenName, player.FamilyName)

	for _, stat := range mP.Stats {
		fmt.Printf("%s: %d\n", stat.Name, stat.Value)
	}

	fmt.Println()
}

// Stat struct representation of various stats
type Stat struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
	ID    string `json:"ID"`
}

// Match is the total data for the a match
type Match struct {
	GameID              string `json:"game_id"`
	MatchID             string `json:"match_id"`
	TournamentID        string `json:"tournament_id"`
	TournamentType      string `json:"tournament_type"`
	SeasonID            string `json:"season_id"`
	GameNumber          int    `json:"game_number"`
	MapID               string `json:"map_id"`
	MapType             string `json:"map_type"`
	ESportsTournamentID int    `json:"esports_tournament_id"`
	ESportsMatchID      int    `json:"esports_match_id"`
	Teams               []Team `json:"teams"`
	Stats               []Stat `json:"stats"`
}

// Team information with breakdown of teamwide stats and players
type Team struct {
	ESportsTeamID int           `json:"esports_team_id"`
	Stats         []Stat        `json:"stats"`
	Players       []MatchPlayer `json:"players"`
}

// MatchPlayer information with their stats for the match
type MatchPlayer struct {
	ESportsPlayerID int    `json:"esports_player_id"`
	Stats           []Stat `json:"stats"`
	Heroes          []Hero `json:"heroes"`
}

// Hero played during the match along with stats for the match
type Hero struct {
	HeroID string `json:"hero_id"`
	Name   string `json:"name"`
	Stats  []Stat `json:"stats"`
}
