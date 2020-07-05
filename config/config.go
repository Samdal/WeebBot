package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

var (
	// Public variables

	//BotPrefix : list of bot prefixes
	BotPrefix []string

	//maps (have to be public for Read and Write to work)

	//MapWaifu : map of current saved waifu's
	MapWaifu map[string]string
	//MapPrefix : contains all users custum prefixes
	MapPrefix map[string]string
	//MapUsr : list of users and asociated accounts
	MapUsr map[string]string
	//MapBlack : list of blacklisted users
	MapBlack map[string]bool

	// Private variables

	config *configStruct
)

type configStruct struct {
	Token     string   `json:"Token"`
	BotPrefix []string `json:"BotPrefix"`
}

//ReadConfig ...
//reads the config file that contains the ID and prefixes
func ReadConfig() error {
	filename := "./config/json/config.json"
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	BotPrefix = config.BotPrefix
	return nil
}

//ReadCustomPrefix ...
//reads the custm prefixes
func ReadCustomPrefix() {
	JSONText, err := ioutil.ReadFile("./config/json/prefix.json")
	if err != nil {
		log.Println(err)
	}
	// Unmarshal json data structure to map
	err = json.Unmarshal(JSONText, &MapPrefix)
	if err != nil {
		panic(err)
	}
}

//SetCustomPrefix ...
//sets a custum prefix for that user
func SetCustomPrefix(idImport string, prefixImport string) {
	delete(MapPrefix, idImport) //deletes value if it already exists in map
	//makes new key and value based on input in discord
	NewMapPrefix := map[string]string{
		idImport: prefixImport,
	}
	//in order to add NewMapPrefix to MapPrefix we use range in a for loop (append does not work with maps)
	for K, V := range MapPrefix {
		NewMapPrefix[K] = V
	}
	//marshal (makes a byte array out of map)
	jsonString, err := json.Marshal(NewMapPrefix)
	if err != nil {
		log.Println(err)
	}
	//pushes byte array to json file
	err = ioutil.WriteFile("./config/json/prefix.json", jsonString, 0644)
	return
}

//ReadWaifu ...
//reads all the waifus
func ReadWaifu() {
	JSONText, err := ioutil.ReadFile("./config/json/waifu.json")
	if err != nil {
		log.Println(err)
	}
	// Unmarshal json data structure to map
	err = json.Unmarshal(JSONText, &MapWaifu)
	if err != nil {
		panic(err)
	}
}

//SetWaifu ...
//sets a users waifu
func SetWaifu(idImport string, waifuImport string) {
	if MapWaifu[idImport] != "" { //deletes value if it already exists in map
		delete(MapWaifu, idImport)
	}
	//makes new key and value based on input in discord
	MewMapWaifu := map[string]string{
		idImport: waifuImport,
	}
	//in order to add MewMapWaifu to MapWaifu we use range in a for loop (append does not work with maps)
	for K, V := range MapWaifu {
		MewMapWaifu[K] = V
	}
	//marshal (makes a byte array out of map)
	jsonString, err := json.Marshal(MewMapWaifu)
	if err != nil {
		log.Println(err)
	}
	//pushes byte array to json file
	err = ioutil.WriteFile("./config/json/waifu.json", jsonString, 0644)
	ReadWaifu() //reads json new file
	return
}

// ReadUsr ...
//reads the file of anime users
func ReadUsr() {
	JSONText, err := ioutil.ReadFile("./config/json/usr.json")
	if err != nil {
		log.Println(err)
	}
	// Unmarshal json data structure to map
	err = json.Unmarshal(JSONText, &MapUsr)
	if err != nil {
		panic(err)
	}
}

//SetUsr ...
//sets the user's anime account
func SetUsr(idImport string, usrImport string) {
	if MapUsr[idImport] != "" { //deletes value if it already exists in map
		delete(MapUsr, idImport)
	}
	//makes new key and value based on input in discord
	MewMapUsr := map[string]string{
		idImport: usrImport,
	}
	//in order to add MewMapUsr to MapUsr we use range in a for loop (append does not work with maps)
	for K, V := range MapUsr {
		MewMapUsr[K] = V
	}
	//marshal (makes a byte array out of map)
	jsonString, err := json.Marshal(MewMapUsr)
	if err != nil {
		log.Println(err)
	}
	//pushes byte array to json file
	err = ioutil.WriteFile("./config/json/usr.json", jsonString, 0644)
	ReadUsr() //reads json new file
	return
}

//ReadBlacklist ...
//reads the black list
func ReadBlacklist() {
	JSONText, err := ioutil.ReadFile("./config/json/blacklist.json")
	if err != nil {
		log.Println(err)
	}
	// Unmarshal json data structure to map
	err = json.Unmarshal(JSONText, &MapBlack)
	if err != nil {
		panic(err)
	}
}

//SetBlacklist ...
//sets the blacklist
func SetBlacklist(idImport string, blacklistImport bool) {
	//makes new key and value based on input in discord
	delete(MapBlack, idImport)
	MewMapBlack := map[string]bool{
		idImport: blacklistImport,
	}
	//in order to add MewMapBlack to MapBlack we use range in a for loop (append does not work with maps)
	for K, V := range MapBlack {
		MewMapBlack[K] = V
	}
	//marshal (makes a byte array out of map)
	jsonString, err := json.Marshal(MewMapBlack)
	if err != nil {
		log.Println(err)
	}
	//pushes byte array to json file
	err = ioutil.WriteFile("./config/json/blacklist.json", jsonString, 0644)
	ReadBlacklist() //reads json new file
	return
}

//AnimeSearchStruct ...
//structure of the anime reply
type AnimeSearchStruct struct {
	Data struct {
		Page struct {
			PageInfo struct {
				LastPage int `json:"lastPage"`
			} `json:"pageInfo"`
			Media []struct {
				Status       string `json:"status"`
				Format       string `json:"format"`
				Episodes     int    `json:"episodes"`
				Duration     int    `json:"duration"`
				Chapters     int    `json:"chapters"`
				Volumes      int    `json:"volumes"`
				AverageScore int    `json:"averageScore"`
				Source       string `json:"source"`
				StartDate    struct {
					Day   int `json:"day"`
					Month int `json:"mont"`
					Year  int `json:"year"`
				}
				EndDate struct {
					Day   int `json:"day"`
					Month int `json:"mont"`
					Year  int `json:"year"`
				}
				Genres  []string `json:"genres"`
				Studios struct {
					Nodes []struct {
						Name string `json:"name"`
					} `json:"nodes"`
				}
				Title struct {
					Romaji  string `json:"romaji"`
					English string `json:"english"`
				}

				Description string `json:"description"`
				CoverImage  struct {
					ExtraLarge string `json:"extraLarge"`
				} `json:"coverImage"`
				SiteURL       string `json:"siteUrl"`
				ExternalLinks []struct {
					Site string `json:"site"`
					URL  string `json:"url"`
				} `json:"externalLinks"`
				Relations struct {
					Edges []struct {
						RelationType string `json:"relationType"`
					} `json:"edges"`
					Nodes []struct {
						Title struct {
							English string `json:"english"`
							Romaji  string `json:"romaji"`
						} `json:"title"`
						ID int `json:"id"`
					} `json:"nodes"`
				} `json:"relations"`
				NextAiringEpisode struct {
					TimeUntilAiring int `json:"timeUntilAiring"`
					Episode         int `json:"episode"`
				} `json:"nextAiringEpisode"`
			} `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

//AnimeSearch ...
//searches for anime
func AnimeSearch(query string, format string, id string) AnimeSearchStruct {

	cmd := exec.Command("python", "config/request.py", query, format, id)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	response, err := ioutil.ReadFile("config/response.json")
	if err != nil {
		log.Println(err)
	}

	var AnimeStruct AnimeSearchStruct

	// Unmarshal json data structure to map
	err = json.Unmarshal(response, &AnimeStruct)
	if err != nil {
		panic(err)
	}

	return AnimeStruct
}
