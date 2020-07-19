package bot //the package of this Go file

import (
	//standard libraries
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings" //(https://golang.org/pkg/strings/)
	"time"

	//custom libraries
	"../config"      //imports config file
	scraper "../web" //imports scraper file

	//third party libraries
	"github.com/bwmarrin/discordgo" //imports discord library
	"github.com/fatih/color"
)

var (
	newupdate         = "New and kawaii"
	pink              = 0x00FF00FF
	nameAuthor        = "DRAVA#8380 | Weeb BOT | Coded in Go"
	searchPage    int = 0
	srch          config.AnimeSearchStruct
	searchh       = srch.Data.Page.Media
	lastSearch    time.Time
	animesearch   *discordgo.Message
	animeSearchID string
	botPrefix     = [8]string{"W!", "w!", "!w", "!W", "1w", "1W", "w1", "W!"}
)

//Start is a function called in main that starts the bot
func Start() {
	//starts timer that prints how long it took to start the bot
	t := time.Now()

	//reads information from json files
	color.Blue("reading waifues")
	config.ReadWaifu()
	color.Blue("reading animelists")
	config.ReadUsr()

	//start a session
	var goBot *discordgo.Session

	//read token
	file, err := os.Open("bot/token")
	if checkerror(err) {
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	//convert to string and remvove potential line break at the end
	token := strings.Replace(string(b), "\n", "", -1)

	//connecting with the bot
	goBot, err = discordgo.New("Bot " + token)
	if checkerror(err) {
		return
	}

	//massegeHandler runs every time a message is sent
	goBot.AddHandler(messageHandler)
	//Open creates a websocket connection to Discord.
	err = goBot.Open()
	if checkerror(err) {
		return
	}

	//finished starting up
	color.Green("Bot is running!")
	fmt.Println("Starting bot took:", time.Since(t))

	//runs in parralell and changes the status
	//alternate between "w!help" and the custum status
	go func() {
		for true {
			goBot.UpdateStatus(1, "w!help")
			time.Sleep(time.Second * 5)
			goBot.UpdateStatus(1, newupdate)
			time.Sleep(time.Second * 5)
		}
	}()
}

//this functions work is to direct information to the right function
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	//i have decided to have a bunch of prefixes in an array
	//these are all variations of the same prefix, 10/10 would not reccomend
	for _, prefix := range botPrefix {
		if strings.HasPrefix(m.Content, string(prefix)+" ") {
			m.Content = strings.Replace(m.Content, string(prefix)+" ", "", 1)
			withPrefix(s, m)
			return

			//checks if someone used a space at the end
		} else if strings.HasPrefix(m.Content, string(prefix)) {
			m.Content = strings.Replace(m.Content, string(prefix), "", 1)
			withPrefix(s, m)
			return
		}
	}
	//if no prefix is found
	withoutprefix(s, m)
	return
}

//all commands that start with the prefix
func withPrefix(s *discordgo.Session, m *discordgo.MessageCreate) {

	//printing in console
	color.Magenta("Time: " + time.Now().String())
	color.Cyan("Username: " + m.Author.Username)
	color.Green("Message: w!" + m.Content)

	//here is a private command that can only be executed by me
	if m.Author.ID == "313703847656816642" {
		//updating bot status
		if stringsHasMultiple("prefix", []string{"update ", "status ", "presence "}, m.Content) {
			//different prefixes need to be removed in order to acces the wanted word
			m.Content = strings.Replace(m.Content, "update ", "", 1)
			m.Content = strings.Replace(m.Content, "status ", "", 1)
			m.Content = strings.Replace(m.Content, "presence ", "", 1)
			newupdate = m.Content
			go s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			return
			//print the emoji, instead of the discord way ":smirk:"
			//this is useful to get the unicode character of an emoji, wich is
			//needed to send an emoji
		}
	}

	//setting the anime library account linked to your discord ID.
	//this will simply check if the suggested new linked anime list has
	//the proper formating and is one of the sites that the bot supports.
	//then it will save the link to the account as a value, and the
	//user ID as the key in a json file
	if stringsHasMultiple("prefix", []string{"set usr ", "set user "}, m.Content) {
		//removing the prefix
		usr := strings.TrimPrefix(m.Content, "set user ")
		usr = strings.TrimPrefix(usr, "set usr ")

		//check if it is the same as the existing linked anime list
		if usr == config.MapUsr[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "Your discord account is already connected to this user")
			return
		}

		//format the url
		url, err := url.Parse(usr)
		if checkerror(err) {
			return
		}

		//check if there are any other "sub-sites" after the user
		path := url.Path
		var i string
		for _, i = range []string{"/users/", "/user/", "/profile/"} {
			if !strings.HasPrefix(path, i) {
				break
			}

			//remove the "/profile/" part
			str := strings.TrimPrefix(path, i)
			//remove the "/" at the end if there is one
			str = strings.TrimSuffix(str, "/")

			if strings.Contains(str, "/") || str == "" {
				break
			}

			//print message
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			go addemoji(url.Hostname(), m.ChannelID, m.ID, s)
			config.SetUsr(m.Author.ID, usr)

			return
		}

		//print error
		s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		str := "Please send the user adress of your anime account.\nthe current supported platforms are Anilist, MAL and kitsu\nit should look like this`"
		str += url.Scheme + url.Hostname() + i + "[your username]`"
		s.ChannelMessageSend(m.ChannelID, str)
		return

	}

	//setting waifu / husbando
	if stringsHasMultiple("prefix", []string{"set waifu ", "set husbando "}, m.Content) {
		//removing the prefix
		waifu := strings.TrimPrefix(m.Content, "set waifu ")
		waifu = strings.TrimPrefix(waifu, "set husbando ")

		//changing waifu
		s.ChannelMessageSend(m.ChannelID, "Your waifu/husbando is now \""+waifu+"\", your previus waifu/husbando was "+config.MapWaifu[m.Author.ID])

		config.SetWaifu(m.Author.ID, waifu)
		return

		//reading someones waifu/husbado
	} else if strings.HasPrefix(m.Content, "waifu <@") && strings.HasSuffix(m.Content, ">") || strings.HasPrefix(m.Content, "husbando <@") && strings.HasSuffix(m.Content, ">") {
		//removing prefix and making a mention
		mention := strings.TrimPrefix(m.Content, "waifu ")
		mention = strings.TrimPrefix(m.Content, "husbando  ")

		//makes a pure id string by removing the junk from a mention command
		id := strings.TrimPrefix(mention, "<@")
		id = strings.TrimSuffix(id, ">")
		id = strings.Replace(id, "!", "", 1)

		//if there is no waifu conected
		if config.MapWaifu[id] == "" {
			s.ChannelMessageSend(m.ChannelID, "currently no waifu/husbando connected with that user")
			return
		}

		//sending message
		s.ChannelMessageSend(m.ChannelID, mention+"'s waifu/husbando is "+config.MapWaifu[id])
		return
	}

	m.Content = strings.ToLower(m.Content)
	//
	// all lower case from this point
	//

	//some simple commands
	if simplecommands(s, m) {
		return
	}

	//getting an animes filler episodes from https://www.animefillerlist.com/
	if stringsHasMultiple("prefix", []string{"filler ", "fill ", "fil "}, m.Content) {
		//removing prefix
		scraperfiller := strings.TrimPrefix(m.Content, "fil ")
		scraperfiller = strings.TrimPrefix(m.Content, "fill ")
		scraperfiller = strings.TrimPrefix(m.Content, "filler ")

		//get scraper result
		res, _ := scraper.Scrape(scraperfiller, "fill")

		//transforming struct (output from web scraper)
		marshalbyte, err := json.Marshal(res) //marshal (makes byte array out of struct)
		if checkerror(err) {
			return
		}
		mapstring := string(marshalbyte) //byte to string

		//cleaning up
		if strings.Contains(mapstring, "Desc") {
			mapstring = strings.ReplaceAll(mapstring, `{`, "")
			mapstring = strings.ReplaceAll(mapstring, `}`, "")
			mapstring = strings.ReplaceAll(mapstring, `[`, "")
			mapstring = strings.ReplaceAll(mapstring, `]`, "")
			mapstring = strings.ReplaceAll(mapstring, `UNKNOWN`, "")
			mapstring = strings.ReplaceAll(mapstring, `"`, "")
			mapstring = strings.ReplaceAll(mapstring, `\n`, "\n")
			mapstring = strings.ReplaceAll(mapstring, `\`, "")

		}

		//check if it has proper formating
		if !strings.Contains(mapstring, "Desc") {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}

		//devide up the string
		mapstring = strings.Replace(mapstring, "Tittle:%//DEVIDER//%,Desc:", "", 1)
		fillerinfo := strings.Split(mapstring, ",Tittle:%//DEVIDER//%,Desc:")

		link := strings.ReplaceAll(scraperfiller, " ", "-")

		//make and send an embed
		filler := &discordgo.MessageEmbed{
			Color:  pink,
			Footer: &discordgo.MessageEmbedFooter{Text: "https://www.animefillerlist.com/shows/" + link},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{Name: "Manga Canon Episodes:", Value: fillerinfo[0]},
				&discordgo.MessageEmbedField{Name: "Mixed Canon/Filler Episodes:", Value: fillerinfo[1]},
				&discordgo.MessageEmbedField{Name: "Filler Episodes:", Value: fillerinfo[2]},
				&discordgo.MessageEmbedField{Name: "Anime Canon Episodes:", Value: fillerinfo[3]},
			},
		}

		s.ChannelMessageSendEmbed(m.ChannelID, filler)
		return
	}

	//if wednesday, send wednesday meme
	dag := &discordgo.MessageEmbed{Color: pink, Image: &discordgo.MessageEmbedImage{URL: "https://goo.gl/DDCsXo"}}
	if m.Content == "dag" {
		//only on wednesday
		if int(time.Now().Weekday()) == 3 {
			s.ChannelMessageSendEmbed(m.ChannelID, dag)
			return
		}
		go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
		return
	}

	//anime commands that can be used with your library
	animecommands := []string{"usr",
		"user", "lib", "library", "cur",
		"wat", "current", "watching", "pla",
		"wan", "planned", "want to watch",
		"w2w", "want 2 watch", "com", "completed",
		"on", "hold", "on hold", "pau", "paused", "dro", "drop", "dropped"}

	//checking if it has an anime command as a prefix, and formates the Users ID
	for _, i := range animecommands {
		if strings.HasPrefix(m.Content, string(i)) {
			//checking if message doesn't contains a user ping
			if config.MapUsr[m.Author.ID] != "" && !strings.Contains(m.Content, " <@") && !strings.HasSuffix(m.Content, ">") {
				m.Content = string(i)
				break
				//if there is no account connected
			} else if config.MapUsr[m.Author.ID] == "" {
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				s.ChannelMessageSend(m.ChannelID, "No account connected")
				return
			}
			//checking if message does contains a user ping
			if strings.Contains(m.Content, " <@") && strings.HasSuffix(m.Content, ">") {
				m.Content = strings.Replace(m.Content, "!", "", -1)
				m.Author.ID = strings.TrimPrefix(m.Content, " <@")
				m.Author.ID = strings.Replace(m.Author.ID, string(i), "", 1)
				m.Author.ID = strings.TrimSuffix(m.Author.ID, ">")
				m.Content = strings.Replace(m.Content, m.Content, string(i), 1)
				break
			} else {
				s.ChannelMessageSend(m.ChannelID, "Please send a valid command. Example: \n`w!completed <@486926277702320148>`")
				return
			}
		}
	}

	//sends the anime account, if it has prefix. There are multiple services
	usraccount := config.MapUsr[m.Author.ID]
	//kitsu
	if strings.Contains(usraccount, "https://kitsu.io/users") {
		if m.Content == "usr" || m.Content == "user" {
			s.ChannelMessageSend(m.ChannelID, usraccount)
			return
		} else if m.Content == "lib" || m.Content == "library" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/library")
			return
		} else if m.Content == "cur" || m.Content == "wat" || m.Content == "current" || m.Content == "watching" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/library?status=current")
			return
		} else if stringsHasMultiple("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m.Content) {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/library?status=planned")
			return
		} else if m.Content == "com" || m.Content == "completed" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/library?status=completed")
			return
		} else if m.Content == "on" || m.Content == "hold" || m.Content == "on hold" || m.Content == "pau" || m.Content == "paused" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/library?status=on_hold")
			return
		} else if m.Content == "dro" || m.Content == "dropped" || m.Content == "drop" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"library?status=dropped")
			return
		}
		//MAL
	} else if strings.Contains(usraccount, "https://myanimelist.net/animelist/") {
		if m.Content == "usr" || m.Content == "user" {
			usraccount = strings.Replace(usraccount, "animelist", "profile", 1)
			s.ChannelMessageSend(m.ChannelID, usraccount)
			return
		} else if m.Content == "lib" || m.Content == "library" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"?status=7")
			return
		} else if m.Content == "cur" || m.Content == "wat" || m.Content == "current" || m.Content == "watching" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"?status=1")
			return
		} else if stringsHasMultiple("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m.Content) {
			s.ChannelMessageSend(m.ChannelID, usraccount+"?status=6")
			return
		} else if m.Content == "com" || m.Content == "completed" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"?status=2")
			return
		} else if m.Content == "on" || m.Content == "hold" || m.Content == "on hold" || m.Content == "pau" || m.Content == "paused" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"?status=3")
			return
		} else if m.Content == "dro" || m.Content == "dropped" || m.Content == "drop" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"library?status=4")
			return
		}
		//anilist
	} else if strings.Contains(usraccount, "https://anilist.co/user/") {
		if m.Content == "usr" || m.Content == "user" {
			usraccount = strings.Replace(usraccount, "animelist", "profile", 1)
			s.ChannelMessageSend(m.ChannelID, usraccount)
			return
		} else if m.Content == "lib" || m.Content == "library" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/animelist")
			return
		} else if m.Content == "cur" || m.Content == "wat" || m.Content == "current" || m.Content == "watching" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/animelist/Watching")
			return
		} else if stringsHasMultiple("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m.Content) {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/animelist/Planning")
			return
		} else if m.Content == "com" || m.Content == "completed" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/animelist/Completed")
			return
		} else if m.Content == "on" || m.Content == "hold" || m.Content == "on hold" || m.Content == "pau" || m.Content == "paused" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"/animelist/Paused")
			return
		} else if m.Content == "dro" || m.Content == "dropped" || m.Content == "drop" {
			s.ChannelMessageSend(m.ChannelID, usraccount+"library/animelist/Dropped")
			return
		}
	}

	//Searching for anime information with the anilist API

	//searching for an anime with a format
	//default format is nothing
	format := "NONE"
	isAnimesearch := false

	//if it is just "a", don't defien a search term
	if strings.HasPrefix(m.Content, "a ") {
		m.Content = strings.TrimPrefix(m.Content, "a ")
		isAnimesearch = true
		//if it is tv short instead of just SHORT
	} else if strings.HasPrefix(m.Content, "tv short ") {
		m.Content = strings.TrimPrefix(m.Content, "tv short ")
		format = "SHORT"
		isAnimesearch = true
		//go through possible other formats
	} else {
		formats := []string{"TV", "SHORT", "MOVIE", "SPECIAL", "OVA", "ONA", "MANGA", "NOVEL", "ONE SHOT", "a"}
		for _, i := range formats {
			if strings.HasPrefix(m.Content, strings.ToLower(i)+" ") {
				format = strings.Replace(i, " ", "_", 1)
				m.Content = strings.Replace(m.Content, strings.ToLower(i)+" ", "", 1)
				isAnimesearch = true
				break
			}
		}
	}

	//search for anime with decided format or without formant
	if isAnimesearch {
		//reset the search page
		searchPage = 0

		//make search
		srch = config.AnimeSearch(strings.ReplaceAll(m.Content, " ", "-"), format, "NONE")

		//if returned data has no information, return an error
		if len(srch.Data.Page.Media) == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}

		//create embed
		animeEmbed := createAnimeEmbed(srch)
		var err error
		//send message
		animesearch, err = s.ChannelMessageSendEmbed(m.ChannelID, &animeEmbed)
		if checkerror(err) {
			return
		}
		//update variables for changing media and relations
		//see withoutPrefix()
		animeSearchID = animesearch.ChannelID
		lastSearch = time.Now()

		return
	}

	//time it takes to watch x amount of anime episodes
	if stringsHasMultiple("suffix", []string{"episode", "ep"}, m.Content) {
		x := strings.Replace(m.Content, " episodes", "", 1)
		x = strings.Replace(x, " ep", "", 1)
		//convert to int
		y, _ := strconv.Atoi(x)

		if y == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}

		//format time
		episodes := y * 23
		var hours int
		var day int

		//converting minutes into hours and days
		//there should be a better way of doint this, but it works
		for ; episodes-60 > 0; episodes = episodes - 60 {
			hours++
		}
		for ; hours-24 > 0; hours = hours - 24 {
			day++
		}

		//convert to string
		h := strconv.Itoa(hours)
		e := strconv.Itoa(episodes)
		d := strconv.Itoa(day)
		//format
		if d != "0" {
			h = d + " days+\n" + h
		}

		//make embed and send it
		embeddedepisode := &discordgo.MessageEmbed{
			Color:  pink,
			Author: &discordgo.MessageEmbedAuthor{Name: h + ":" + e},
			Footer: &discordgo.MessageEmbedFooter{Text: "h : m"},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embeddedepisode)
		return
	}

	color.Red("Exited at end of function. Command not recognized")
}

//all simple commands
func simplecommands(s *discordgo.Session, m *discordgo.MessageCreate) bool {

	//all embedded messages for the help commands
	help := &discordgo.MessageEmbed{
		Color:  pink,
		Author: &discordgo.MessageEmbedAuthor{Name: nameAuthor},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "anime user account", Value: "`w!help user`", Inline: true},
			&discordgo.MessageEmbedField{Name: "weakly pleasure", Value: "`w!dag`", Inline: true},
			&discordgo.MessageEmbedField{Name: "anime filler episodes", Value: "`w!help filler`", Inline: true},
			&discordgo.MessageEmbedField{Name: "anime searching", Value: "`w!help anime`", Inline: true},
			&discordgo.MessageEmbedField{Name: "animes this season", Value: "`w!help s`", Inline: true},
			&discordgo.MessageEmbedField{Name: "waifu/husbando", Value: "`w!help waifu/h...`", Inline: true},
			&discordgo.MessageEmbedField{Name: "lists simple commands", Value: "`w!help simple`", Inline: true},
		},
		//image at bottom and image in upper right corner
		Image:     &discordgo.MessageEmbedImage{URL: "https://goo.gl/KXZCw3"},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://goo.gl/1WzVwB"},
	}

	//help anime
	helpAnime := &discordgo.MessageEmbed{
		Color:  pink,
		Author: &discordgo.MessageEmbedAuthor{Name: "Makes a search on some streaming sites or for anime information on Anilist"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "w!cru/w!crunchyroll \"anime\"", Value: "Example: w!cunchyroll naruto", Inline: false},
			&discordgo.MessageEmbedField{Name: "w!kiss \"anime\"", Value: "Example: w!kiss dragon ball", Inline: false},
			&discordgo.MessageEmbedField{Name: "w!9 \"anime\"", Value: "Example: w!9 One piece", Inline: false},
			&discordgo.MessageEmbedField{Name: "w!t/w!twist \"anime\"", Value: "Example: w!t bakemonogatari", Inline: false},
			&discordgo.MessageEmbedField{Name: "w!a \"anime/manga/LN/movie\"", Value: "Example: w!a fma\nreturns information about your chosen show", Inline: false},
			&discordgo.MessageEmbedField{Name: "w!movie/manga/LN/OVA \"anime search\"", Value: "Example: w!movie Wolf children\n specifyes the search to a format", Inline: false},
			&discordgo.MessageEmbedField{Name: "r[n]", Value: "After searching with the two alst comands, you can select to view a relation\nExample: `r2`", Inline: false},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fanilist.co%2Fimg%2Ficons%2Fandroid-chrome-512x512.png&f=1&nofb=1"},
	}

	//help user
	helpuser := &discordgo.MessageEmbed{
		Color:  pink,
		Author: &discordgo.MessageEmbedAuthor{Name: "There are multiple commands that can be used with an anime account"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "w!set usr [user link]", Value: "Example:\nw!set user https://anilist.co/user/Drava/animelist\nw!set user https://kitsu.io/users/drava\nw!set user https://myanimelist.net/profile/Drava", Inline: false},
			&discordgo.MessageEmbedField{Name: "Looking up your own account", Value: "w!lib(library)\nw!cur/wat(current/watching)\nw!pla/wan (planned/want2watch)\nw!com(completed)\nw!hold(on hold)\nw!drop(dropped)", Inline: false},
			&discordgo.MessageEmbedField{Name: "Looking another users account", Value: "w!lib(library) @DiscordUser\nw!cur/wat(current/watching) @DiscordUser\nw!pla/wan (planned/want2watch) @DiscordUser\nw!com(completed) @DiscordUser\nw!hold(on hold) @DiscordUser\nw!drop(dropped) @DiscordUser", Inline: false},
		},
	}

	//help filler
	helpfiller := &discordgo.MessageEmbed{
		Color:  pink,
		Author: &discordgo.MessageEmbedAuthor{Name: "Makes a search in animefillerlist for the anime"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "w!filler \"anime\"", Value: "Example:\nw!filler naruto\nw!fill boruto naruto next generations\nw!fil Bleach", Inline: false},
		},
	}

	//help waifu
	helpwaifu := &discordgo.MessageEmbed{
		Color:  pink,
		Author: &discordgo.MessageEmbedAuthor{Name: "you can set your own personal waifu/husbando, "},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Setting your own waifu / husbando", Value: "w!set waifu \"waifu\" \nw!husbando \"husbando\"", Inline: false},
			&discordgo.MessageEmbedField{Name: "Viewing your own waifu / husbando", Value: "w!waifu\nw!husbando", Inline: false},
			&discordgo.MessageEmbedField{Name: "Viewing someone else's waifu / husbando", Value: "w!waifu @discord user\nhusbando @DiscordUser", Inline: false},
		},
	}

	//embeds
	embeds := map[string]*discordgo.MessageEmbed{
		"help":          help,
		"hjelp":         help,
		"help anime":    helpAnime,
		"help user":     helpuser,
		"help filler":   helpfiller,
		"help waifu":    helpwaifu,
		"help husbando": helpwaifu,
	}

	//itereate throug every prefix and send embed
	for i := range embeds {
		if m.Content == string(i) {
			if string(i) == "yare" || string(i) == "excuse" {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embeds[i])
			return true
		}
	}

	//all simple direct text commands commands
	simple := map[string]string{
		"waifu":       "Your waifu is " + config.MapWaifu[m.Author.ID],
		"husbando":    "Your husbando is " + config.MapWaifu[m.Author.ID],
		"help define": "define has sadly been removed from this bot ;(",
		"help filler": "makes a search in animefillerlist for the anime\n" + "`w!filler \"filler_anime\"`",
		"help s":      "opens anichart (popular anime this season)\nmore comands planned\n`" + `w!s` + "`",
		"help waifu":  "you can set your own personal waifu\n`" + "w!set waifu \"waifu\"\nw!waifu @user (returns user's waifu)\nw!waifu (returns your waifu)" + "`",
		"ping":        "i hear you loud and clear",
		"invite":      "https://discord.com/oauth2/authorize?client_id=486926277702320148&scope=bot&permissions=0",
	}

	var helpsimple string
	//range all commands (keys) and send the value
	for key := range simple {
		//sending the "help simple command" (see help command in discord)
		if m.Content == "help simple" {
			helpsimple = helpsimple + "*" + key + "*                `" + simple[key] + "`, "
			continue
		}
		//sending simple commands
		if m.Content == key {
			s.ChannelMessageSend(m.ChannelID, simple[key])
			return true
		}

	}
	if helpsimple != "" {
		s.ChannelMessageSend(m.ChannelID, "`"+helpsimple+"`")
	}

	return false
}

//withoutprefix does stuff where there is no prefix
func withoutprefix(s *discordgo.Session, m *discordgo.MessageCreate) {

	//if there has been an anime search in the current channel and it is 60 seconds since tha last search
	//continue on from this point
	//this only works because there are no multiple instances of the anime search, since this is still quite a small bot
	if m.ChannelID != animeSearchID || time.Since(lastSearch).Seconds() > 60 || len(srch.Data.Page.Media) == 0 {
		return
	}

	//next
	if m.Content == "n" {

		searchPage++
		//go to 0 if above max
		if searchPage >= len(srch.Data.Page.Media)-1 {
			searchPage = 0
		}
		go s.ChannelMessageDelete(m.ChannelID, m.ID)
		animeEmbed := createAnimeEmbed(srch)
		_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
		checkerror(err)
		return
	}
	//previus
	if m.Content == "b" {

		searchPage--
		//go to max if below 0
		if searchPage < 0 {
			searchPage = len(srch.Data.Page.Media) - 1
		}

		go s.ChannelMessageDelete(m.ChannelID, m.ID)
		animeEmbed := createAnimeEmbed(srch)
		_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
		checkerror(err)
		return
	}

	//choosing relations
	var chooseMax int
	if len(srch.Data.Page.Media[searchPage].Relations.Edges) > 3 {
		chooseMax = 3
	} else {
		chooseMax = len(srch.Data.Page.Media[searchPage].Relations.Edges)
	}

	for i := 0; i <= chooseMax; i++ {
		if strings.HasPrefix(m.Content, "r"+strconv.Itoa(i)) {
			//get and search the anime relation
			srch = config.AnimeSearch(strings.ReplaceAll(m.Content, " ", "-"), "NONE", strconv.Itoa(srch.Data.Page.Media[searchPage].Relations.Nodes[i].ID))

			//error check
			if len(srch.Data.Page.Media) == 0 {
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				return
			}
			//send and save new anime message
			go s.ChannelMessageDelete(m.ChannelID, m.ID)
			animeEmbed := createAnimeEmbed(srch)
			_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
			checkerror(err)
			return
		}
	}

	return
}

//adds emojis to the message
func addemoji(hostname string, channelid string, mid string, s *discordgo.Session) {
	switch hostname {
	case "kitsu.io":
		s.MessageReactionAdd(channelid, mid, "🦊")
	case "anilist.co":
		s.MessageReactionAdd(channelid, mid, "🇲")
		s.MessageReactionAdd(channelid, mid, "🇦")
		s.MessageReactionAdd(channelid, mid, "🇱")
	case "myanimelist.net":
		s.MessageReactionAdd(channelid, mid, "🇦")
		s.MessageReactionAdd(channelid, mid, "🇱")
	}
}

//Iterates through an array of strings and checks if it has any of
//the strings in the array as a prefix, suffix or contains it
//depends on the desired mdoe
func stringsHasMultiple(stringtype string, commands []string, content string) bool {
	for _, i := range commands {
		if stringtype == "prefix" {
			if strings.HasPrefix(content, i) {
				return true
			}
		}
		if stringtype == "content" {
			if strings.Contains(content, i) {
				return true
			}
		}
		if stringtype == "suffix" {
			if strings.HasSuffix(content, i) {
				return true
			}
		}
		if stringtype == "message" {
			if content == i {
				return true
			}
		}
	}
	return false
}

//Simply checks and prints the error, returns true if err != nil
func checkerror(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

//checks if a string is empty and fills it with none
func checkNil(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

//capitalizes the first letter of a word
func toLowerCapital(s string) string {
	if checkNil(s) == "none" {
		return checkNil(s)
	}
	var nstring string
	var nextupper bool
	for i := 0; i < len(s); i++ {
		if i == 0 {
			nstring += strings.ToUpper(string(s[i]))
		} else if s[i] == ' ' {
			nextupper = true
			nstring += string(s[i])
		} else if nextupper {
			nstring += strings.ToUpper(string(s[i]))
			nextupper = false
		} else {
			nstring += strings.ToLower(string(s[i]))
		}
	}
	return nstring
}

//creates an embed and does the anime search
func createAnimeEmbed(srch config.AnimeSearchStruct) discordgo.MessageEmbed {

	//getting the search and printing to console
	search := srch.Data.Page.Media[searchPage]
	fmt.Println(srch.Data.Page.Media[searchPage].Title.English)
	status := checkNil(search.Status) + "\n*"

	//formating date
	date := search.StartDate
	edate := search.EndDate
	if search.Status == "FINISHED" {
		if date.Day == edate.Day && date.Month == edate.Month && date.Year == edate.Year {
			status += strconv.Itoa(date.Day) + "-" + strconv.Itoa(date.Month) + "-" + strconv.Itoa(date.Year) + "*"
		} else {
			status += strconv.Itoa(date.Day) + "-" + strconv.Itoa(date.Month) + "-" + strconv.Itoa(date.Year)
			status += " -> " + strconv.Itoa(edate.Day) + "-" + strconv.Itoa(edate.Month) + "-" + strconv.Itoa(edate.Year) + "*"
		}

	} else {
		status += strconv.Itoa(date.Day) + "-" + strconv.Itoa(date.Month) + "-" + strconv.Itoa(date.Year) + " ->*"
		if search.Status == "RELEASING" {
			seconds := search.NextAiringEpisode.TimeUntilAiring
			days := search.NextAiringEpisode.TimeUntilAiring / 86400
			hours := (seconds - days*86400) / 3600
			minutes := (seconds - (days*86400 + hours*3600)) / 60

			timeUntill := strconv.Itoa(days) + "d " + strconv.Itoa(hours) + "h " + strconv.Itoa(minutes) + "m**"
			status += "\n**Ep " + strconv.Itoa(search.NextAiringEpisode.Episode) + ": " + timeUntill
		}
	}

	//formating the genres
	var genres string
	for _, i := range search.Genres {
		genres = genres + "\n" + i
	}

	//formating relations
	var relations string
	for i := range search.Relations.Edges {
		if i > 3 {
			break
		}
		relations += "[" + strconv.Itoa(i) + "]: **" + toLowerCapital(search.Relations.Edges[i].RelationType) + "**\n"
		if search.Relations.Nodes[i].Title.English != "" {
			relations += search.Relations.Nodes[i].Title.English + "\n"
		} else {
			relations += search.Relations.Nodes[i].Title.Romaji + "\n"
		}
	}
	if relations == "" {
		relations = "none"
	}

	//formating studios
	var studios string
	if len(search.Studios.Nodes) == 0 {
		studios = "none"
	} else {
		for _, i := range search.Studios.Nodes {
			studios = studios + "\n" + i.Name
		}
	}

	//episode duration and count
	aformat := toLowerCapital(strings.Replace(search.Format, "_", "", 1))
	episodes := "Episodes"
	epInfo := checkNil(strconv.Itoa(search.Episodes))
	duration := "Ep Duration"
	duInfo := checkNil(strconv.Itoa(search.Duration)) + "min"
	var externalLinks string

	//change to volumes and chapters if it is a novel or a manga
	if search.Format == "NOVEL" || search.Format == "MANGA" {
		episodes = "Volumes"
		epInfo = checkNil(strconv.Itoa(search.Volumes))
		duration = "Chapters"
		duInfo = checkNil(strconv.Itoa(search.Chapters))
		externalLinks = "none"
	} else {

		//format external links if it is not a novel or manga
		if len(search.ExternalLinks) != 0 {
			for _, i := range search.ExternalLinks {
				//only show certan websites
				if i.Site == "Crunchyroll" || i.Site == "Netflix" {
					externalLinks += "\n" + i.Site + ": " + i.URL
				}
				continue
			}
		} else {
			externalLinks = "none"
		}

		//generate a twist.moe link
		illegall := []string{":", "'", "/", ".", ",", ";", "\\", "*", "☆"}

		twist := strings.ToLower(strings.ReplaceAll(search.Title.Romaji, " ", "-"))
		for _, i := range illegall {
			if i == ";" || i == "☆" {
				twist = strings.ReplaceAll(twist, i, "-")
				continue
			}
			twist = strings.ReplaceAll(twist, i, "")
		}
		externalLinks += "\nTwist.moe: https://twist.moe/a/" + twist
	}

	//format description
	description := checkNil(search.Description)

	htmlstf := []string{"<br>", "<i>", "</i>"}
	for _, i := range htmlstf {
		if i == "<i>" || i == "</i>" {
			description = strings.ReplaceAll(description, i, "*")
			continue
		}
		description = strings.ReplaceAll(description, i, "")
	}

	characters := len(description)

	//make the max length of the description 1000
	//this is a horrible way to do it
	//someone please help me
	var desc string
	maxSize := 1000
	if characters >= maxSize {
		for i := 0; i < len(description); i++ {
			if i <= maxSize {
				desc += string(description[i])
			} else {
				break
			}
		}
		description = desc + "[...]"
	}

	//format score
	score := checkNil(strconv.Itoa(search.AverageScore)) + "/100"
	if score == "none/100" {
		score = "none"
	}

	//Foramt ittle
	title := search.Title.English + " | " + search.Title.Romaji
	if search.Title.English == "" {
		title = search.Title.Romaji
	}

	//make the embed
	animeEmbed := &discordgo.MessageEmbed{
		Color: pink,
		Author: &discordgo.MessageEmbedAuthor{
			Name: title,
			URL:  search.SiteURL,
		},
		Footer: &discordgo.MessageEmbedFooter{Text: "<<[b]       |" + strconv.Itoa(searchPage+1) + "/" + strconv.Itoa(len(srch.Data.Page.Media)) + "|       [n]>>"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Description", Value: description, Inline: false},
			&discordgo.MessageEmbedField{Name: "Status", Value: status, Inline: false},
			&discordgo.MessageEmbedField{Name: "Format", Value: aformat, Inline: true},
			&discordgo.MessageEmbedField{Name: episodes, Value: epInfo, Inline: true},
			&discordgo.MessageEmbedField{Name: duration, Value: duInfo, Inline: true},
			&discordgo.MessageEmbedField{Name: "Studios", Value: studios, Inline: true},
			&discordgo.MessageEmbedField{Name: "Average Score", Value: score, Inline: true},
			&discordgo.MessageEmbedField{Name: "Source", Value: checkNil(toLowerCapital(strings.ReplaceAll(search.Source, "_", " "))), Inline: true},
			&discordgo.MessageEmbedField{Name: "Relations", Value: relations, Inline: true},
			&discordgo.MessageEmbedField{Name: "External Links", Value: externalLinks, Inline: true},
			&discordgo.MessageEmbedField{Name: "Genres", Value: genres, Inline: true},
		},
		//image at bottom and image in upper right corner
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: search.CoverImage.ExtraLarge},
	}

	return *animeEmbed
}
