package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/exec"
)

var (
	//MapWaifu : map of current saved waifu's
	MapWaifu map[string]string
	//MapUsr : users and asociated accounts
	MapUsr map[string]string
)

//ReadWaifu ...
//reads all the waifus
func ReadWaifu() {
	//read json file
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
	//deletes value if it already exists in map
	if MapWaifu[idImport] != "" {
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

	//reads json new file
	ReadWaifu()
	return
}

// ReadUsr ...
//reads the file of anime users
func ReadUsr() {
	//reads json file
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
	//deletes value if it already exists in map
	if MapUsr[idImport] != "" {
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

	//reads json new file
	ReadUsr()
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

	//do the HTTP POST request in python
	//i do this since i couldn't find out how to do it in GO
	//someone please help
	cmd := exec.Command("python", "config/request.py", query, format, id)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	//read json
	response, err := ioutil.ReadFile("config/response.json")
	if err != nil {
		log.Println(err)
	}

	//Format into struct
	var AnimeStruct AnimeSearchStruct
	err = json.Unmarshal(response, &AnimeStruct)
	if err != nil {
		panic(err)
	}

	//returns the search struct
	return AnimeStruct
}
