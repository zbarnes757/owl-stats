package leagueapi

import (
	"encoding/json"
	"fmt"
	"net/http"
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
