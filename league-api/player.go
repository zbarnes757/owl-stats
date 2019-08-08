package leagueapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"owl-stats/models"
	"time"
)

var playerURI = "https://api.overwatchleague.com/players/"

var playerClient = &http.Client{
	Timeout: time.Second * 3,
}

// Player detailed information
type Player struct {
	ID           int
	HomeLocation string
	FamilyName   string
	GivenName    string
}

type playerResponse struct {
	Data playerData `json:"data"`
}

type playerData struct {
	Player Player `json:"player"`
}

type playersResp struct {
	Content []models.OWLPlayer `json:"content"`
}

// FetchAllPlayers will return a full list of all players in the OW League
func FetchAllPlayers() []models.OWLPlayer {
	log.Println("Sending request for all player data to " + playerURI)

	resp, err := playerClient.Get(playerURI)
	if err != nil {
		return []models.OWLPlayer{}
	}
	defer resp.Body.Close()

	var res playersResp
	json.NewDecoder(resp.Body).Decode(&res)
	return res.Content
}

func fetchPlayer(id int) (Player, error) {
	fmt.Printf("Fetching data for player #%d\n", id)

	uri := fmt.Sprintf("%s%d", playerURI, id)
	resp, err := playerClient.Get(uri)
	if err != nil {
		return Player{}, err
	}
	defer resp.Body.Close()

	var r playerResponse
	json.NewDecoder(resp.Body).Decode(&r)
	return r.Data.Player, nil
}
