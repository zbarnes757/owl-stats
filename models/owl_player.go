package models

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// OWLPlayer represents data for a player from the overwatch leage
type OWLPlayer struct {
	Base
	EsportsID    int    `json:"id" gorm:"unique;not null"`
	Name         string `json:"name"`
	HomeLocation string `json:"homeLocation"`
	FamilyName   string `json:"familyName"`
	GivenName    string `json:"givenName"`
	Nationality  string `json:"nationality"`
	Headshot     string `json:"headshot"`
}

// BulkCreateOWLPlayers will save multiple records at once
func BulkCreateOWLPlayers(players []OWLPlayer) error {
	valueStrings := []string{}
	valueArgs := []interface{}{}
	now := time.Now()

	for _, player := range players {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

		valueArgs = append(valueArgs, uuid.NewV4())
		valueArgs = append(valueArgs, now)
		valueArgs = append(valueArgs, now)
		valueArgs = append(valueArgs, player.EsportsID)
		valueArgs = append(valueArgs, player.Name)
		valueArgs = append(valueArgs, player.HomeLocation)
		valueArgs = append(valueArgs, player.FamilyName)
		valueArgs = append(valueArgs, player.GivenName)
		valueArgs = append(valueArgs, player.Nationality)
		valueArgs = append(valueArgs, player.Headshot)
	}

	smt := `
		INSERT INTO owl_players(
			id,
			created_at,
			updated_at,
			esports_id,
			name,
			home_location,
			family_name,
			given_name,
			nationality,
			headshot
		) VALUES %s
		ON CONFLICT ON CONSTRAINT owl_players_esports_id_key
		DO UPDATE SET
			updated_at = EXCLUDED.updated_at,
			name = EXCLUDED.name,
			home_location = EXCLUDED.home_location,
			family_name = EXCLUDED.family_name,
			given_name = EXCLUDED.given_name,
			nationality = EXCLUDED.nationality,
			headshot = EXCLUDED.headshot;
	`

	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	tx := db.Begin()
	if err := tx.Exec(smt, valueArgs...).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// {
//             "id": 5865,
//             "availableLanguages": [
//                 "en",
//                 "zh-cn"
//             ],
//             "handle": "ksf.4600",
//             "name": "KSF",
//             "homeLocation": "Puyallup, WA",
//             "accounts": [
//                 {
//                     "id": 7351,
//                     "competitorId": 5865,
//                     "value": "https://www.twitch.tv/ksfksf",
//                     "accountType": "TWITCH",
//                     "isPublic": true
//                 },
//                 {
//                     "id": 7288,
//                     "competitorId": 5865,
//                     "value": "https://twitter.com/ksf",
//                     "accountType": "TWITTER",
//                     "isPublic": true
//                 }
//             ],
//             "game": "OVERWATCH",
//             "attributes": {
//                 "heroes": [
//                     "genji",
//                     "soldier-76",
//                     "sombra",
//                     "widowmaker"
//                 ],
//                 "player_number": 27,
//                 "preferred_slot": "1",
//                 "role": "offense"
//             },
//             "attributesVersion": "1.0.0",
//             "familyName": "Frandanisa",
//             "givenName": "Kyle",
//             "nationality": "US",
//             "headshot": "https://bnetcmsus-a.akamaihd.net/cms/gallery/FDHS3GHI2JY31549580254953.png",
//             "teams": [
//                 {
//                     "team": {
//                         "id": 4405,
//                         "availableLanguages": [
//                             "en",
//                             "zh-cn",
//                             "th",
//                             "ko",
//                             "ru",
//                             "ja",
//                             "zh-tw"
//                         ],
//                         "handle": "los-angeles-a.2779",
//                         "name": "Los Angeles Valiant",
//                         "homeLocation": "Los Angeles, CA",
//                         "primaryColor": "004438",
//                         "secondaryColor": "000000",
//                         "accounts": [
//                             {
//                                 "id": 2293,
//                                 "competitorId": 4405,
//                                 "value": "https://www.twitch.tv/lavaliant",
//                                 "accountType": "TWITCH",
//                                 "isPublic": true
//                             },
//                             {
//                                 "id": 2291,
//                                 "competitorId": 4405,
//                                 "value": "https://www.facebook.com/LAValiant",
//                                 "accountType": "FACEBOOK",
//                                 "isPublic": true
//                             },
//                             {
//                                 "id": 2292,
//                                 "competitorId": 4405,
//                                 "value": "https://twitter.com/lavaliant",
//                                 "accountType": "TWITTER",
//                                 "isPublic": true
//                             },
//                             {
//                                 "id": 2379,
//                                 "competitorId": 4405,
//                                 "value": "https://www.instagram.com/lavaliant",
//                                 "accountType": "INSTAGRAM",
//                                 "isPublic": true
//                             },
//                             {
//                                 "id": 2294,
//                                 "competitorId": 4405,
//                                 "value": "https://www.youtube.com/channel/UCYH1uEM9EBVP6ddxa5gsvGw",
//                                 "accountType": "YOUTUBE_CHANNEL",
//                                 "isPublic": true
//                             },
//                             {
//                                 "id": 2429,
//                                 "competitorId": 4405,
//                                 "value": "https://discord.gg/lavaliant",
//                                 "accountType": "DISCORD",
//                                 "isPublic": true
//                             }
//                         ],
//                         "game": "OVERWATCH",
//                         "attributes": {
//                             "city": null,
//                             "hero_image": null,
//                             "manager": null,
//                             "team_guid": "0x0D70000000000005"
//                         },
//                         "attributesVersion": "1.0.0",
//                         "divisions": [
//                             {
//                                 "competitor": {
//                                     "id": 4405
//                                 },
//                                 "division": {
//                                     "id": 61
//                                 }
//                             },
//                             {
//                                 "competitor": {
//                                     "id": 4405
//                                 },
//                                 "division": {
//                                     "id": 80
//                                 }
//                             }
//                         ],
//                         "abbreviatedName": "VAL",
//                         "addressCountry": "US",
//                         "logo": "https://bnetcmsus-a.akamaihd.net/cms/page_media/0D8BNUWVZP6B1508792362890.PNG",
//                         "icon": "https://bnetcmsus-a.akamaihd.net/cms/template_resource/L3U59GQVS1ZK1507822882879.svg",
//                         "secondaryPhoto": "https://bnetcmsus-a.akamaihd.net/cms/template_resource/L3U59GQVS1ZK1507822882879.svg",
//                         "type": "TEAM"
//                     },
//                     "player": {
//                         "id": 5865,
//                         "type": "PLAYER"
//                     },
//                     "flags": []
//                 }
//             ],
//             "user": {
//                 "id": 151
//             },
//             "type": "PLAYER"
//         }
