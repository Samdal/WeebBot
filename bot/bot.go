package bot //the package of this Go file

import (
	//standard libraries
	"encoding/json"
	"fmt" //format
	"io/ioutil"
	"math/rand"
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
	newupdate     = "New and kawaii"
	lastreport    string
	timeupdate    time.Duration = 5
	farge                       = 0x00FF00FF
	nameAuthor                  = "DRAVA#8380 | Weeb BOT | Coded in Go"
	bug                         = false
	searchPage    int           = 0
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
	color.Blue("reading custom prefixes")
	config.ReadCustomPrefix()
	color.Blue("reading waifues")
	config.ReadWaifu()
	color.Blue("reading user list")
	config.ReadUsr()

	//start a session
	var goBot *discordgo.Session

	//read the token
	file, err := os.Open("bot/token")
	if checkerror(err) {
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	token := string(b)

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
	go func() { //runs in parralell and changes the status
		for true {
			goBot.UpdateStatus(1, "w!help")
			time.Sleep(time.Second * timeupdate)
			goBot.UpdateStatus(1, newupdate)
			time.Sleep(time.Second * timeupdate)
		}
	}()
}

//this function's work is to direct information to the right function
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	//i have decided to have a bunch of prefixes in an array
	for _, prefix := range botPrefix { //checks all prefixes in slice, with space behind it, and without
		if strings.HasPrefix(m.Content, string(prefix)+" ") {
			m.Content = strings.Replace(m.Content, string(prefix)+" ", "", 1)
			withPrefix(s, m)
			return

			//checks if someone used a space at the end
		} else if strings.HasPrefix(m.Content, string(prefix)) {
			m.Content = strings.Replace(m.Content, string(prefix), "", 1)
			withPrefix(s, m)
			return

			//checks if a user used a custum prefix
		} else if strings.HasPrefix(m.Content, config.MapPrefix[m.Author.ID]) {
			//Checks if the user has a custom prefix wich is being called
			if strings.HasPrefix(m.Content, config.MapPrefix[m.Author.ID]+" ") {
				m.Content = strings.Replace(m.Content, config.MapPrefix[m.Author.ID]+" ", "", 1)
				withPrefix(s, m)
				return
			}
		}
	}
	//if no prefix is found
	withoutprefix(s, m)
	return
}

//withoutprefix does stuff where there is no prefix
func withoutprefix(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != animeSearchID || time.Since(lastSearch).Seconds() > 60 {
		return
	}

	var chooseMax int
	if m.Content == "n" {
		if len(srch.Data.Page.Media) == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}
		searchPage++
		if searchPage >= len(srch.Data.Page.Media)-1 {
			searchPage = 0
		}
		go s.ChannelMessageDelete(m.ChannelID, m.ID)
		animeEmbed := createAnimeEmbed(srch)
		_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
		checkerror(err)
		return
	}
	if m.Content == "b" {
		if len(srch.Data.Page.Media) == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}
		searchPage--
		if searchPage < 0 {
			searchPage = len(srch.Data.Page.Media) - 1
		}
		go s.ChannelMessageDelete(m.ChannelID, m.ID)
		animeEmbed := createAnimeEmbed(srch)
		_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
		checkerror(err)
		return
	}

	if strings.HasPrefix(m.Content, "f") {
		if len(srch.Data.Page.Media) == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}
		if len(srch.Data.Page.Media[searchPage].Relations.Edges) > 3 {
			chooseMax = 3
		} else {
			chooseMax = len(srch.Data.Page.Media[searchPage].Relations.Edges)
		}

		for i := 0; i <= chooseMax; i++ {
			if strings.HasPrefix(m.Content, "r"+strconv.Itoa(i)) {
				srch = config.AnimeSearch(strings.ReplaceAll(m.Content, " ", "-"), "NONE", strconv.Itoa(srch.Data.Page.Media[searchPage].Relations.Nodes[i].ID))
				if len(srch.Data.Page.Media) == 0 {
					go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
					return
				}
				go s.ChannelMessageDelete(m.ChannelID, m.ID)
				animeEmbed := createAnimeEmbed(srch)
				_, err := s.ChannelMessageEditEmbed(animesearch.ChannelID, animesearch.ID, &animeEmbed)
				checkerror(err)
				return
			}
		}
	}

	return
}

//all commands that start with the prefix
func withPrefix(s *discordgo.Session, m *discordgo.MessageCreate) {

	//printing in console
	color.Magenta("Time: " + time.Now().String())
	color.Cyan("Username: " + m.Author.Username)
	color.Green("Message: w!" + m.Content)

	if bug {
		s.ChannelMessageSend(m.ChannelID, "Bot is currently in debugging mode, som features might not work as intended")
	}
	var ( //variables
		//other
		weekday = time.Now().Weekday()
		//int
		day = 3
		//string
	)

	//here are some private commands
	if m.Author.ID == "313703847656816642" {

		//updating bot status
		if stringsIterate("prefix", []string{"update ", "status ", "presence "}, m) {
			//different prefixes need to be removed in order to acces the wanted word
			m.Content = strings.Replace(m.Content, "update ", "", 1)
			m.Content = strings.Replace(m.Content, "status ", "", 1)
			m.Content = strings.Replace(m.Content, "presence ", "", 1)
			newupdate = m.Content
			go s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			return
			//sending message back to somone who uses the report function of the bot
		} else if strings.HasPrefix(m.Content, "emoji ") {
			m.Content = strings.Replace(m.Content, "emoji ", "", 1)
			return
			//says what is stated in the message
		} else if strings.HasPrefix(m.Content, "say ") {
			m.Content = strings.Replace(m.Content, "say ", "", 1)
			s.ChannelMessageSend(m.ChannelID, m.Content)
			return
		} else if strings.HasPrefix(m.Content, "bug") {
			if bug {
				bug = false
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			} else {
				bug = true
				s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			}
			return
		}

	}
	//
	//Global commands
	//
	if m.Content == "test" {
		t := time.Now()
		s.ChannelMessageSend(m.ChannelID, "test")
		var end = time.Since(t)
		s.ChannelMessageSend(m.ChannelID, "it took: "+end.String()+" to recive and send the message")
		return
	}
	//use this as a prefix to dm the message to yourself
	if stringsIterate("prefix", []string{"dm ", "pm "}, m) {
		//creats a direct message channel / private message channel
		dm, _ := s.UserChannelCreate(m.Author.ID)
		m.Content = strings.Replace(m.Content, "dm ", "", 1)
		m.Content = strings.Replace(m.Content, "pm ", "", 1)
		m.ChannelID = dm.ID
	}

	//setting the anime library account linked to your discord adress
	if stringsIterate("prefix", []string{"set usr ", "set user "}, m) {
		usr := strings.Replace(m.Content, "set user ", "", 1)
		usr = strings.Replace(usr, "set usr ", "", 1)

		if usr == config.MapUsr[m.Author.ID] {
			s.ChannelMessageSend(m.ChannelID, "Your discord account is already connected to this user")
			return

			//checks if the string sent is correct, and checks what service the link is sent from
		} else if strings.Contains(usr, "https://kitsu.io/users/") {
			check := strings.Replace(m.Content, "https://kitsu.io/users/", "", 1)
			if strings.Contains(check, "/") {
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				go addemoji("kitsu", m.ChannelID, m.ID, s)
				s.ChannelMessageSend(m.ChannelID, "Please send the user adress of your kitsu account.\nit should look lik this:\n`https://kitsu.io/users/[your_user_name/number]`")
				return
			}
			go s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			go addemoji("kitsu", m.ChannelID, m.ID, s)
			//config.go saves the url and user adress in a JSON
			config.SetUsr(m.Author.ID, usr)
			return
			//checks if the string sent is correct
		} else if stringsIterate("contains", []string{"https://myanimelist.net/profile/", "https://myanimelist.net/animelist/"}, m) {
			usr := strings.Replace(usr, "profile", "animelist", 1)
			check := strings.Replace(usr, "https://myanimelist.net/animelist/", "", 1)
			if strings.Contains(check, "?") || strings.Contains(check, "/") {
				s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				go addemoji("mal", m.ChannelID, m.ID, s)
				s.ChannelMessageSend(m.ChannelID, "Please send the user adress of your MAL account.\nit should look lik this:\n`https://myanimelist.net/profile/[your_user_name/number]`")
				return
			}
			//"????" is just a placeholder for new
			usr = strings.Replace(usr, "https://myanimelist.net/animelist/", "????", 1)
			usr = strings.Replace(usr, "/", "", 1)
			usr = strings.Replace(usr, "????", "https://myanimelist.net/animelist/", 1)
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			go addemoji("mal", m.ChannelID, m.ID, s)
			//config.go saves the url and user adress in a JSON
			config.SetUsr(m.Author.ID, usr)
			return
		} else if strings.Contains(usr, "https://anilist.co/user/") {
			check := strings.Replace(m.Content, "https://anilist.co/user/", "", 1)
			if strings.Contains(check, "/") {
				s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				go addemoji("anilist", m.ChannelID, m.ID, s)

				s.ChannelMessageSend(m.ChannelID, "Please send the user adress of your Anilist account.\nit should look lik this:\n`https://anilist.co/user/[your_user_name]`")
				return
			}
			s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
			go addemoji("anilist", m.ChannelID, m.ID, s)

			//config.go saves the url and user adress in a JSON
			config.SetUsr(m.Author.ID, usr)
			return
		}
		//TODO add more options for different websites
		s.ChannelMessageSend(m.ChannelID, "Please send the user adress of your anime account.\nit should look lik this:`https://[your_prefered_anime_site].[com/net/co]/[user_or_profile]/[your_user_name/number]`")
		s.ChannelMessageSend(m.ChannelID, "Currently we only support Anilist, MAL and kitsu. \nDoes the anime service your trying to connect to not exist? please contact with  w!report")
		return
	}

	//setting waifu / husbando
	if stringsIterate("prefix", []string{"set waifu ", "set husbando "}, m) {
		//sets a new waifu/husbando
		waifu := strings.Replace(m.Content, "set waifu ", "", 1)
		waifu = strings.Replace(waifu, "set husbando ", "", 1)

		if waifu != "" {
			s.ChannelMessageSend(m.ChannelID, "Your waifu/husbando is now \""+waifu+"\", your previus waifu/husbando was "+config.MapWaifu[m.Author.ID])
			config.SetWaifu(m.Author.ID, waifu)
			return
		}
		//reading someones waifu/husbado
	} else if strings.HasPrefix(m.Content, "waifu <@") && strings.HasSuffix(m.Content, ">") || strings.HasPrefix(m.Content, "husbando <@") && strings.HasSuffix(m.Content, ">") {
		mention := strings.Replace(m.Content, "waifu ", "", 1)
		mention = strings.Replace(m.Content, "husbando  ", "", 1)

		//makes a pure id string by removing the junk from a mention command
		id := strings.Replace(m.Content, "waifu <@", "", 1)
		id = strings.Replace(m.Content, "husbando <@", "", 1)
		id = strings.Replace(id, ">", "", 1)
		id = strings.Replace(id, "!", "", 1)
		if config.MapWaifu[id] == "" {
			s.ChannelMessageSend(m.ChannelID, "currently no waifu/husbando connected with that user")
			return
		}
		s.ChannelMessageSend(m.ChannelID, mention+"'s waifu/husbando is "+config.MapWaifu[id])
		return
	}
	m.Content = strings.ToLower(m.Content)
	//
	// all lower case from this point
	//

	if simplecommands(s, m) { //some simple commands (see function for reference)
		return
	}

	//sending report messages to me
	if m.Author.ID != "313703847656816642" && stringsIterate("prefix", []string{"report ", "bug ", "msg "}, m) {
		id := "user " + m.Author.ID + " sent a report: "
		msg := strings.Replace(m.Content, "msg ", id, 1)
		msg = strings.Replace(m.Content, "report ", id, 1)
		msg = strings.Replace(m.Content, "bug ", id, 1)
		s.ChannelMessageSend("493708042420879360", msg)
		go s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
		lastreport = m.Author.ID //lastreport is a global wariable in this file
		return
	}
	//setting a custom personal prefix
	if strings.HasPrefix(m.Content, "set prefix") {
		prefix := strings.Replace(m.Content, "set prefix ", "", 1)
		if prefix != "" {
			s.ChannelMessageSend(m.ChannelID, "Your new custom prefix is `"+prefix+"`, your previus custom prefix was `"+config.MapPrefix[m.Author.ID]+"`")
			//config.go function saves your personal prefix in a JSON
			config.SetCustomPrefix(m.Author.ID, prefix)
			return
		}
	}

	if stringsIterate("prefix", []string{"filler ", "fill ", "fil "}, m) {
		scraperfiller := strings.Replace(m.Content, "fil ", "", 1)
		scraperfiller = strings.Replace(m.Content, "fill ", "", 1)
		scraperfiller = strings.Replace(m.Content, "filler ", "", 1)

		//cleanup, cleans and activates the scrape
		mstring := cleanup(scraperfiller, "fill", 2)
		if !strings.Contains(mstring, "Desc") {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}
		mstring = strings.Replace(mstring, "Tittle:%//DEVIDER//%,Desc:", "", 1)

		fillerinfo := strings.Split(mstring, ",Tittle:%//DEVIDER//%,Desc:")
		footertxt := strings.Replace(scraperfiller, " ", "-", -1)

		//sends filler episodes
		filler := &discordgo.MessageEmbed{
			Color:  farge,
			Footer: &discordgo.MessageEmbedFooter{Text: "https://www.animefillerlist.com/shows/" + footertxt},
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

	//translating latin characters into japanese characters
	if stringsIterate("prefix", []string{"hiragana ", "hir "}, m) {

		//legal characters
		const legall = `っゃゅょかきくけこがぎぐげごさしつすせそざじずぜぞたちてとだぢづでどなにぬねのはひふへほばびぶべぼぱぴぺぽまみぷめもやゆよらりるれろわをあいうえおん。()[]{}「」〜:!`
		var hiragana string

		if strings.HasPrefix(m.Content, "hir ") {
			hiragana = strings.Replace(m.Content, "hir ", "", 1)
		} else if strings.HasPrefix(m.Content, "hiragana ") {
			hiragana = strings.Replace(m.Content, "hiragana ", "", 1)
		}

		//turns characters into japanese
		leg := language("hir", hiragana)

		//checks if it still contains latin characters
		for _, char := range leg {
			if !strings.Contains(legall, string(char)) {
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				return
			}
		}

		hiraganaemb := &discordgo.MessageEmbed{
			Color:  farge,
			Author: &discordgo.MessageEmbedAuthor{Name: leg},
			Footer: &discordgo.MessageEmbedFooter{Text: hiragana},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, hiraganaemb)
		return
	}
	{
		shouldreturn := false
		var str = m.Content
		if strings.HasPrefix(str, "weeb neko ") {
			str = strings.Replace(str, "weeb neko ", "", 1)
			str = language("weeb", str)
			str = language("neko", str)
			shouldreturn = true
		} else if strings.HasPrefix(str, "weeb ") {
			str = strings.Replace(str, "weeb ", "", 1)
			str = language("weeb", str)
			shouldreturn = true
		} else if strings.HasPrefix(str, "neko ") {
			str = strings.Replace(str, "neko ", "", 1)
			str = language("neko", str)
			shouldreturn = true
		}
		if shouldreturn {
			strembed := &discordgo.MessageEmbed{
				Color: farge,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{Name: m.Author.Username, Value: str, Inline: true},
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, strembed)
			return
		}
	}
	//if wednesday, send wednesday meme
	dag := &discordgo.MessageEmbed{Color: farge, Image: &discordgo.MessageEmbedImage{URL: "https://goo.gl/DDCsXo"}}
	if m.Content == "dag" {
		if int(weekday) == day {
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
			if config.MapUsr[m.Author.ID] != "" && !strings.Contains(m.Content, " <@") && !strings.HasSuffix(m.Content, ">") {
				m.Content = string(i)
				break
			} else if config.MapUsr[m.Author.ID] == "" && !strings.Contains(m.Content, " <@") && !strings.HasSuffix(m.Content, ">") {
				go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
				s.ChannelMessageSend(m.ChannelID, "No account connected")
				return
			}
			if strings.Contains(m.Content, " <@") && strings.HasSuffix(m.Content, ">") {
				m.Content = strings.Replace(m.Content, "!", "", -1)
				m.Author.ID = strings.Replace(m.Content, " <@", "", 1)
				m.Author.ID = strings.Replace(m.Author.ID, string(i), "", 1)
				m.Author.ID = strings.Replace(m.Author.ID, ">", "", 1)
				m.Content = strings.Replace(m.Content, m.Content, string(i), 1)
				break
			} else {
				s.ChannelMessageSend(m.ChannelID, "Please send a valid command. Example: \n`w!completed <@486926277702320148>`")
				return
			}
		}
	}

	//TODO more services

	//sends the anime account, if it has prefix. There are multiple services
	usraccount := config.MapUsr[m.Author.ID]
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
		} else if stringsIterate("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m) {
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
		} else if stringsIterate("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m) {
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
		} else if stringsIterate("prefix", []string{"pla", "wan", "planned", "want to watch", "w2w", "want 2 watch"}, m) {
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

	formats := []string{"TV", "SHORT", "MOVIE", "SPECIAL", "OVA", "ONA", "MANGA", "NOVEL", "ONE SHOT", "a"}
	format := "NONE"
	var isAnimesearch bool
	for _, i := range formats {
		if strings.HasPrefix(m.Content, "a ") {
			m.Content = strings.Replace(m.Content, "a ", "", 1)
			isAnimesearch = true
			break
		} else if strings.HasPrefix(m.Content, "tv short ") {
			m.Content = strings.Replace(m.Content, "tv short ", "", 1)
			isAnimesearch = true
			break
		}
		if strings.HasPrefix(m.Content, strings.ToLower(i)+" ") {
			format = strings.Replace(i, " ", "_", 1)
			m.Content = strings.Replace(m.Content, strings.ToLower(i)+" ", "", 1)
			isAnimesearch = true
			break
		}
	}

	if isAnimesearch {
		searchPage = 0

		srch = config.AnimeSearch(strings.ReplaceAll(m.Content, " ", "-"), format, "NONE")

		if len(srch.Data.Page.Media) == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return
		}

		animeEmbed := createAnimeEmbed(srch)
		var err error
		animesearch, err = s.ChannelMessageSendEmbed(m.ChannelID, &animeEmbed)
		checkerror(err)
		animeSearchID = animesearch.ChannelID
		lastSearch = time.Now()
		return
	}

	color.Red("Exited at end of function. Command not recognized")
}

//all simple commands
func simplecommands(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	var helpsimple string

	//help embed
	help := &discordgo.MessageEmbed{
		Color:  farge,
		Author: &discordgo.MessageEmbedAuthor{Name: nameAuthor},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "anime user account", Value: "`w!help user`", Inline: true},
			&discordgo.MessageEmbedField{Name: "weakly pleasure", Value: "`w!dag`", Inline: true},
			&discordgo.MessageEmbedField{Name: "anime filler episodes", Value: "`w!help filler`", Inline: true},
			&discordgo.MessageEmbedField{Name: "anime searching", Value: "`w!help anime`", Inline: true},
			&discordgo.MessageEmbedField{Name: "animes this season", Value: "`w!help s`", Inline: true},
			&discordgo.MessageEmbedField{Name: "language of the weeb", Value: "`w!help hir`", Inline: true},
			&discordgo.MessageEmbedField{Name: "custom prefix", Value: "`w!help prefix`", Inline: true},
			&discordgo.MessageEmbedField{Name: "waifu/husbando", Value: "`w!help waifu/h...`", Inline: true},
			&discordgo.MessageEmbedField{Name: "weeb life made easy", Value: "`w!help tips`", Inline: true},
			&discordgo.MessageEmbedField{Name: "is there a bug?", Value: "`w!help report`", Inline: true},
			&discordgo.MessageEmbedField{Name: "list's simple commands", Value: "`w!help simple`", Inline: true},
			&discordgo.MessageEmbedField{Name: "sentence changing", Value: "`w!help language`", Inline: true},
		},
		//image at bottom and image in upper right corner
		Image:     &discordgo.MessageEmbedImage{URL: "https://goo.gl/KXZCw3"},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: "https://goo.gl/1WzVwB"},
	}

	//embeded messages (posts a picture instead of url link)
	daze := &discordgo.MessageEmbed{Color: farge, Image: &discordgo.MessageEmbedImage{URL: "https://goo.gl/9kcYHA"}}
	Ex := &discordgo.MessageEmbed{Color: farge, Image: &discordgo.MessageEmbedImage{URL: "https://goo.gl/M6Wtz7"}}
	pat := &discordgo.MessageEmbed{Color: farge, Image: &discordgo.MessageEmbedImage{URL: "https://goo.gl/ekHWCV"}}

	helpAnime := &discordgo.MessageEmbed{
		Color:  farge,
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

	helpuser := &discordgo.MessageEmbed{
		Color:  farge,
		Author: &discordgo.MessageEmbedAuthor{Name: "There are multiple commands that can be used with an anime account"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "w!set usr [user link]", Value: "Example:\nw!set user https://anilist.co/user/Drava/animelist\nw!set user https://kitsu.io/users/drava\nw!set user https://myanimelist.net/profile/Drava", Inline: false},
			&discordgo.MessageEmbedField{Name: "Looking up your own account", Value: "w!lib(library)\nw!cur/wat(current/watching)\nw!pla/wan (planned/want2watch)\nw!com(completed)\nw!hold(on hold)\nw!drop(dropped)", Inline: false},
			&discordgo.MessageEmbedField{Name: "Looking another users account", Value: "w!lib(library) @DiscordUser\nw!cur/wat(current/watching) @DiscordUser\nw!pla/wan (planned/want2watch) @DiscordUser\nw!com(completed) @DiscordUser\nw!hold(on hold) @DiscordUser\nw!drop(dropped) @DiscordUser", Inline: false},
		},
	}

	helpfiller := &discordgo.MessageEmbed{
		Color:  farge,
		Author: &discordgo.MessageEmbedAuthor{Name: "Makes a search in animefillerlist for the anime"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "w!filler \"anime\"", Value: "Example:\nw!filler naruto\nw!fill boruto naruto next generations\nw!fil Bleach", Inline: false},
		},
	}

	helpwaifu := &discordgo.MessageEmbed{
		Color:  farge,
		Author: &discordgo.MessageEmbedAuthor{Name: "you can set your own personal waifu/husbando, "},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Setting your own waifu / husbando", Value: "w!set waifu \"waifu\" \nw!husbando \"husbando\"", Inline: false},
			&discordgo.MessageEmbedField{Name: "Viewing your own waifu / husbando", Value: "w!waifu\nw!husbando", Inline: false},
			&discordgo.MessageEmbedField{Name: "Viewing someone else's waifu / husbando", Value: "w!waifu @discord user\nhusbando @DiscordUser", Inline: false},
		},
	}

	helplanguage := &discordgo.MessageEmbed{
		Color:  farge,
		Author: &discordgo.MessageEmbedAuthor{Name: "You can chage your sentence in amazingly weeb ways!!"},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Translate romanji to hiragana", Value: "w!hir \"sentence\" \nw!hiragana \"sentence\"", Inline: false},
			&discordgo.MessageEmbedField{Name: "nyaow you can talk like a neko nya", Value: "w!neko \"sentence\"", Inline: false},
			&discordgo.MessageEmbedField{Name: "anta kan also talk like a true weeb, my tomodachi", Value: "w!weeb \"sentence\"", Inline: false},
			&discordgo.MessageEmbedField{Name: "anyata kan combine the powers myasterfully", Value: "w!weeb neko \"sentence\"", Inline: false},
		},
	}

	//embeds
	embeds := map[string]*discordgo.MessageEmbed{
		"yare":          daze,
		"excuse":        Ex,
		"pat":           pat,
		"help":          help,
		"hjelp":         help,
		"help anime":    helpAnime,
		"help user":     helpuser,
		"help filler":   helpfiller,
		"help waifu":    helpwaifu,
		"help husbando": helpwaifu,
		"help language": helplanguage,
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

	//all simple commands
	simple := map[string]string{
		"waifu":       "Your waifu is " + config.MapWaifu[m.Author.ID],
		"husbando":    "Your husbando is " + config.MapWaifu[m.Author.ID],
		"prefix":      "Your custom prefix is `" + config.MapPrefix[m.Author.ID] + "`",
		"help define": "define has sadly been removed from this bot ;(",
		"ani":         "https://kitsu.io/anime/",
		"cru":         "http://www.crunchyroll.com/",
		"man":         "https://kitsu.io/manga/",
		"kiss":        "https://kissanime.ac/",
		"9":           "https://www1.9anime.to/",
		"s":           "https://www.livechart.me", //TODO improve with scraper
		"t":           "https://twist.moe",
		"twist":       "https://twist.moe",
		"help hir":    "converts latin alphabet into hiragana\n`w!hir`",
		"help filler": "makes a search in animefillerlist for the anime\n" + "`w!filler \"filler_anime\"`",
		"help s":      "opens anichart (popular anime this season)\nmore comands planned\n`" + `w!s` + "`",
		"help prefix": "you can set your own personal custom prefix\n`" + `w!set prefix "your new prefix", w!prefix (returns your custom prefix)` + "`",
		"help waifu":  "you can set your own personal waifu\n`" + "w!set waifu \"waifu\"\nw!waifu @user (returns user's waifu)\nw!waifu (returns your waifu)" + "`",
		"help report": "you can send direct messeges to Author with these commands, if you spam you will be blacklisted\n`" + `w!msg/bug/report <your message>` + "`",
		"help tips":   "Use the prefix `dm` or `pm` after the main prefix to direct message your answer\n",
		"update":      "you need weebgod lisence to do that!",
		"status":      "you need weebgod lisence to do that!",
		"presence":    "you need weebgod lisence to do that!",
		"ping":        "i hear you loud and clear",
		"kitsu":       "https://kitsu.io",
		"kit":         "https://kitsu.io",
		"k":           "https://kitsu.io",
		"mal":         "https://myanimelist.net/",
		"m":           "https://myanimelist.net/",
		"anilist":     "https://anilist.co",
		"lenny":       "( ͡° ͜ʖ ͡°)",
		"invite":      "https://discord.com/oauth2/authorize?client_id=486926277702320148&scope=bot&permissions=0",
	}

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

	link := map[string]string{
		"cru ":         "http://www.crunchyroll.com/",
		"crunchyroll ": "http://www.crunchyroll.com/",
		"kiss ":        "http://kissanime.ru/Anime/",
		"twist ":       "https://twist.moe/a/",
		"t ":           "https://twist.moe/a/",
	}
	//changes links based on input
	for urlprefix := range link {
		if strings.HasPrefix(m.Content, urlprefix) {
			urlresult := strings.Replace(m.Content, urlprefix, link[urlprefix], -1)
			urlresult = strings.Replace(urlresult, " ", "-", -1)
			s.ChannelMessageSend(m.ChannelID, urlresult)
			return true
		}
	}

	//the only link generator that uses + instead of -
	if strings.HasPrefix(m.Content, "9 ") {
		anime9 := strings.Replace(m.Content, "9 ", "https://www1.9anime.to/search?keyword=", -1)
		anime9 = strings.Replace(anime9, " ", "+", -1)
		s.ChannelMessageSend(m.ChannelID, anime9)
		return true
	}

	//time it takes to watch x amount of anime episodes
	if stringsIterate("suffix", []string{"episode", "ep"}, m) {
		x := strings.Replace(m.Content, " episodes", "", 1)
		x = strings.Replace(x, " ep", "", 1)
		/** converting the str1 variable into an int using Atoi method */
		y, _ := strconv.Atoi(x)

		if y == 0 {
			go s.MessageReactionAdd(m.ChannelID, m.ID, "❌")
			return true
		}

		episodes := y * 23
		var hours int
		var day int
		for ; episodes-60 > 0; episodes = episodes - 60 {
			hours++
		}
		for ; hours-24 > 0; hours = hours - 24 {
			day++
		}
		h := strconv.Itoa(hours)
		e := strconv.Itoa(episodes)
		d := strconv.Itoa(day)
		if d != "0" {
			h = d + " days+\n" + h
		}
		// s.ChannelMessageSend(m.ChannelID, x+" episodes is "+h+" hours and "+e+" minutes")
		embeddedepisode := &discordgo.MessageEmbed{
			Color:  farge,
			Author: &discordgo.MessageEmbedAuthor{Name: h + ":" + e},
			Footer: &discordgo.MessageEmbedFooter{Text: "h : m"},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embeddedepisode)
		return true
	}
	return false
}

//adds emojis to the message
func addemoji(service string, channelid string, mid string, s *discordgo.Session) {
	if service == "kitsu" {
		s.MessageReactionAdd(channelid, mid, "🦊")
	} else if service == "mal" {
		s.MessageReactionAdd(channelid, mid, "🇲")
		s.MessageReactionAdd(channelid, mid, "🇦")
		s.MessageReactionAdd(channelid, mid, "🇱")
	} else if service == "anilist" {
		s.MessageReactionAdd(channelid, mid, "🇦")
		s.MessageReactionAdd(channelid, mid, "🇳")
		s.MessageReactionAdd(channelid, mid, "🇮")
		s.MessageReactionAdd(channelid, mid, "🇱")
		s.MessageReactionAdd(channelid, mid, "ℹ")
		s.MessageReactionAdd(channelid, mid, "🇸")
		s.MessageReactionAdd(channelid, mid, "🇹")
	}
}

func replaceletters(strmap map[string]string, str string) string {
	for i := range strmap {
		str = strings.Replace(str, i, strmap[i], -1)
	}
	return str
}

func replaceword(strmap map[string]string, str string) string {
	arr := []string{" ", " ", ".", "!", "?", "'", ","}
	for i := range strmap {
		for _, suffix := range arr {
			str = strings.Replace(str, " "+i+suffix, " "+strmap[i]+suffix, -1)
		}
	}
	return str
}

func language(lang string, str string) string {

	if lang == "hir" {
		japanese := map[string]string{
			" ":   "",
			"v":   "b",
			"l":   "r",
			"oo":  "oう",
			"ee":  "eい",
			"kk":  "っk",
			"gg":  "っg",
			"ss":  "っs",
			"zz":  "っz",
			"tt":  "っt",
			"dd":  "っd",
			"nn":  "っn",
			"hh":  "っh",
			"bb":  "っb",
			"pp":  "っp",
			"mm":  "っm",
			"rr":  "っr",
			"cc":  "っc",
			"kya": "きゃ",
			"kyu": "きゅ",
			"kyo": "きょ",
			"sha": "しゃ",
			"shu": "しゅ",
			"sho": "しょ",
			"cha": "ちゃ",
			"chu": "ちゅ",
			"cho": "ちょ",
			"nya": "にゃ",
			"nyu": "にゅ",
			"nyo": "にょ",
			"hya": "ひゃ",
			"hyu": "ひゅ",
			"hyo": "ひょ",
			"mya": "みゃ",
			"myu": "みゅ",
			"myo": "みょ",
			"rya": "りゃ",
			"ryu": "りゅ",
			"ryo": "りょ",
			"gya": "ぎゃ",
			"gyu": "ぎゅ",
			"gyo": "ぎょ",
			"bya": "びゃ",
			"byu": "びゅ",
			"byo": "びょ",
			"pya": "ぴゃ",
			"pyu": "ぴゅ",
			"pyo": "ぴょ",
			"jya": "じゃ",
			"jyu": "じゅ",
			"jyo": "じょ",
			"ja":  "じゃ",
			"ju":  "じゅ",
			"jo":  "じょ",
			"ji":  "じ",
			"ka":  "か",
			"ki":  "き",
			"ku":  "く",
			"ke":  "け",
			"ko":  "こ",
			"ga":  "が",
			"gi":  "ぎ",
			"gu":  "ぐ",
			"ge":  "げ",
			"go":  "ご",
			"sa":  "さ",
			"shi": "し",
			"tsu": "つ",
			"tu":  "つ",
			"su":  "す",
			"se":  "せ",
			"so":  "そ",
			"za":  "ざ",
			"zi":  "じ",
			"zu":  "ず",
			"ze":  "ぜ",
			"zo":  "ぞ",
			"ta":  "た",
			"chi": "ち",
			"te":  "て",
			"to":  "と",
			"da":  "だ",
			"di":  "ぢ",
			"du":  "づ",
			"de":  "で",
			"do":  "ど",
			"na":  "な",
			"ni":  "に",
			"nu":  "ぬ",
			"ne":  "ね",
			"no":  "の",
			"ha":  "は",
			"hi":  "ひ",
			"hu":  "ふ",
			"fu":  "ふ",
			"he":  "へ",
			"ho":  "ほ",
			"ba":  "ば",
			"bi":  "び",
			"bu":  "ぶ",
			"be":  "べ",
			"bo":  "ぼ",
			"pa":  "ぱ",
			"pi":  "ぴ",
			"pu":  "ぷ",
			"pe":  "ぺ",
			"po":  "ぽ",
			"ma":  "ま",
			"mi":  "み",
			"mu":  "ぷ",
			"me":  "め",
			"mo":  "も",
			"ya":  "や",
			"yu":  "ゆ",
			"yo":  "よ",
			"ra":  "ら",
			"ri":  "り",
			"ru":  "る",
			"re":  "れ",
			"ro":  "ろ",
			"wa":  "わ",
			"wo":  "を",
			"a":   "あ",
			"i":   "い",
			"u":   "う",
			"e":   "え",
			"o":   "お",
			"n":   "ん",
			".":   "。",
		}
		order := []string{
			" ",
			"v",
			"l",
			"oo",
			"ee",
			"kk",
			"gg",
			"ss",
			"zz",
			"tt",
			"dd",
			"nn",
			"hh",
			"bb",
			"pp",
			"mm",
			"rr",
			"cc",
			"kya",
			"kyu",
			"kyo",
			"sha",
			"shu",
			"sho",
			"cha",
			"chu",
			"cho",
			"nya",
			"nyu",
			"nyo",
			"hya",
			"hyu",
			"hyo",
			"mya",
			"myu",
			"myo",
			"rya",
			"ryu",
			"ryo",
			"gya",
			"gyu",
			"gyo",
			"bya",
			"byu",
			"byo",
			"pya",
			"pyu",
			"pyo",
			"jya",
			"jyu",
			"jyo",
			"ja",
			"ju",
			"jo",
			"ji",
			"ka",
			"ki",
			"ku",
			"ke",
			"ko",
			"ga",
			"gi",
			"gu",
			"ge",
			"go",
			"sa",
			"shi",
			"tsu",
			"tu",
			"su",
			"se",
			"so",
			"za",
			"zi",
			"zu",
			"ze",
			"zo",
			"ta",
			"chi",
			"te",
			"to",
			"da",
			"di",
			"du",
			"de",
			"do",
			"na",
			"ni",
			"nu",
			"ne",
			"no",
			"ha",
			"hi",
			"hu",
			"fu",
			"he",
			"ho",
			"ba",
			"bi",
			"bu",
			"be",
			"bo",
			"pa",
			"pi",
			"pu",
			"pe",
			"po",
			"ma",
			"mi",
			"mu",
			"me",
			"mo",
			"ya",
			"yu",
			"yo",
			"ra",
			"ri",
			"ru",
			"re",
			"ro",
			"wa",
			"wo",
			"a",
			"i",
			"u",
			"e",
			"o",
			"n",
			".",
		}
		for _, v := range order {
			str = strings.Replace(str, v, japanese[v], -1)
		}
		str = strings.ReplaceAll(str, `"`, "「") //open
		str = strings.ReplaceAll(str, `"`, "」") //close
		return str
	}

	str = " " + str

	if lang == "weeb" {
		weebish := map[string]string{
			" i ": " watashi ",
		}
		str = replaceletters(weebish, str)

		weebword := map[string]string{
			"nice to meet you": "Hajimemashite",
			"hello":            "Konnichiwa",
			"good morning":     "ohayou",
			"it":               "kore",
			"amazing":          "sugoi desu",
			"great":            "sugoi desu",
			"tremendus":        "sugoi desu",
			"good":             "sugoi",
			"cool":             "kakkoi",
			"cute":             "kawaii desu",
			"yes":              "hai",
			"your":             "anta no",
			"you":              "anata",
			"my":               "boku no",
			"that":             "sore",
			"what":             "nani",
			"tea":              "ocha",
			"ocean":            "umi",
			"blue":             "aoi",
			"red":              "aka",
			"black":            "kuro",
			"white":            "shiro",
			"green":            "midori",
			"pink":             "pinku",
			"friend":           "tomodachi",
			"comic":            "manga",
			"tv":               "anime",
			"sword":            "katana",
			"summer":           "natsu",
		}

		str = replaceword(weebword, str)

	}

	if lang == "neko" {
		neko := map[string]string{
			"ma":  "mya",
			"na":  "nya",
			"per": "purr",
			"no":  "nyo",
			"me":  "mye",
			"ne":  "nye",
			"mu":  "myu",
			"nu":  "nyu",
			"ka":  "kya",
			"ha":  "hya",
			"pa":  "paw",
			"pos": "paw",
		}
		str = replaceletters(neko, str)
		nekoWords := map[string]string{
			"wow": "meow",
		}
		str = replaceword(nekoWords, str)

		specialCases := []string{
			"*meow*",
			"*nya*",
			"*umya*",
			"*nya*",
			"",
			"",
		}
		var strarr = strings.Split(str, "")
		for p, i := range strarr {
			arr := []string{"", "!", "?", ",", "."}
			for _, suffix := range arr {
				strarr[p] = strings.Replace(i, suffix, " "+specialCases[rand.Intn(len(specialCases))]+i, 1)
			}
		}
		strarr = append(strarr, " "+specialCases[rand.Intn(len(specialCases))])
		str = ""
		for _, i := range strarr {
			str = str + i
		}
		str = specialCases[rand.Intn(len(specialCases))] + " " + str
	}

	return str
}

//cleans up mess from the website scraper
func cleanup(s string, mode string, amount int) string {
	var Mapstring string
	var keywords = []string{s}
	for _, keyword := range keywords {
		res, _ := scraper.Scrape(keyword, mode, amount)
		//transforming struct (output from web scraper)
		marshalbyte, err := json.Marshal(res) //marshal (makes byte array out of struct)
		if err != nil {
			panic(err)
		}
		Mapstring = string(marshalbyte) //byte to string
	}
	if strings.Contains(Mapstring, "Desc") {
		Mapstring = strings.Replace(Mapstring, `{`, "", -1)
		Mapstring = strings.Replace(Mapstring, `}`, "", -1)
		Mapstring = strings.Replace(Mapstring, `[`, "", -1)
		Mapstring = strings.Replace(Mapstring, `]`, "", -1)
		Mapstring = strings.Replace(Mapstring, `UNKNOWN`, "", -1)
		Mapstring = strings.Replace(Mapstring, `"`, "", -1)
		Mapstring = strings.Replace(Mapstring, `\n`, "\n", -1)
		Mapstring = strings.Replace(Mapstring, `\`, "", -1)

		return Mapstring
	}
	return "Could not find what your searching for"
}

func stringsIterate(stringtype string, commands []string, m *discordgo.MessageCreate) bool {
	for _, i := range commands {
		if stringtype == "prefix" {
			if strings.HasPrefix(m.Content, i) {
				return true
			}
		}
		if stringtype == "content" {
			if strings.Contains(m.Content, i) {
				return true
			}
		}
		if stringtype == "suffix" {
			if strings.HasSuffix(m.Content, i) {
				return true
			}
		}
		if stringtype == "message" {
			if m.Content == i {
				return true
			}
		}
	}
	return false
}

func checkerror(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

func checkNil(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

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

func createAnimeEmbed(srch config.AnimeSearchStruct) discordgo.MessageEmbed {

	search := srch.Data.Page.Media[searchPage]
	fmt.Println(srch.Data.Page.Media[searchPage].Title.English)
	status := checkNil(search.Status) + "\n*"

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

	var genres string
	for _, i := range search.Genres {
		genres = genres + "\n" + i
	}
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

	var studios string
	if len(search.Studios.Nodes) == 0 {
		studios = "none"
	} else {
		for _, i := range search.Studios.Nodes {
			studios = studios + "\n" + i.Name
		}
	}
	var externalLinks string

	aformat := toLowerCapital(strings.Replace(search.Format, "_", "", 1))
	episodes := "Episodes"
	epInfo := checkNil(strconv.Itoa(search.Episodes))
	duration := "Ep Duration"
	duInfo := checkNil(strconv.Itoa(search.Duration)) + "min"
	if search.Format == "NOVEL" || search.Format == "MANGA" {
		episodes = "Volumes"
		epInfo = checkNil(strconv.Itoa(search.Volumes))
		duration = "Chapters"
		duInfo = checkNil(strconv.Itoa(search.Chapters))
		externalLinks = "none"

	} else {

		if len(search.ExternalLinks) != 0 {
			for _, i := range search.ExternalLinks {
				if i.Site == "Crunchyroll" || i.Site == "Netflix" {
					externalLinks += "\n" + i.Site + ": " + i.URL
				}
				continue
			}
		} else {
			externalLinks = "none"
		}

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

	score := checkNil(strconv.Itoa(search.AverageScore)) + "/100"
	if score == "none/100" {
		score = "none"
	}
	title := search.Title.English + " | " + search.Title.Romaji
	if search.Title.English == "" {
		title = search.Title.Romaji
	}

	animeEmbed := &discordgo.MessageEmbed{
		Color: farge,
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
