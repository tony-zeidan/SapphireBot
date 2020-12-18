package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Clinet/discordgo-embed"
	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"github.com/sanzaru/go-giphy"
)

var (
	validCommands []CommandMapping
	validMap      map[string]CommandMapping
	giphyLib      *libgiphy.Giphy
)

type CommandMapping struct {
	Triggers    []string
	Description string
	Syntax      string
	Executor    interface{}
	SubCommands []CommandData
}

//Data structure for containing command data
type CommandData struct {
	//command arguments
	Args []string
	//API message object
	Message *discordgo.Message
	//author of the command
	Author *discordgo.User
	//channel id from which the command was obtained
	ChannelID string
}

func init() {
	giphyLib = libgiphy.NewGiphy(GIPHY_API_TOKEN)
	validMap = make(map[string]CommandMapping)

	validCommands = []CommandMapping{
		{
			Triggers:    []string{"help", "info"},
			Description: "Obtain information about Sapphire's commands.",
			Executor:    helpCommand},
		{
			Triggers:    []string{"hello", "hi", "greetings"},
			Description: "Your personal way of greeting Sapphire.",
			Executor:    helloCommand},
		{
			Triggers:    []string{"roll", "rand"},
			Description: "Roll a random number between two givens.",
			Executor:    rollCommand},
		{
			Triggers:    []string{"giphy", "gif", "gifsearch"},
			Description: "Search Giphy for any gifs.",
			Executor:    giphySearchCommand},
		{
			Triggers:    []string{"freq", "occurrences"},
			Description: "Sapphire will output the frequency of each word in the message.",
			Executor:    occurrencesCommand},
	}

	for _, v := range validCommands {
		for _, v2 := range v.Triggers {
			validMap[v2] = v
		}
	}
}

//respond to the user asking for help with the bots commands by sending a list of available commands
func helpCommand(s *discordgo.Session, data *CommandData) {

	paginated := dgwidgets.NewPaginator(s, data.ChannelID)

	embedded := embed.NewEmbed()

	//we don't want the bot to print aliases (just first trigger)
	for i, v := range validCommands {

		if (i+1)%2 == 0 {
			paginated.Add(embedded.MessageEmbed)
			embedded = embed.NewEmbed()
		}
		fmt.Println(i)
		embedded.SetTitle("Page " + strconv.Itoa(i+1))
		embedded.AddField(v.Triggers[0], v.Description)
	}

	// Sets the footers of all added pages to their page numbers.
	paginated.SetPageFooters()

	// When the paginator is done listening set the colour to yellow
	paginated.ColourWhenDone = 0xffff

	// Stop listening for reaction events after five minutes
	paginated.Widget.Timeout = time.Minute * 5

	paginated.Spawn()
}

//respond to the user saying hello
func helloCommand(s *discordgo.Session, data *CommandData) {
	_, _ = s.ChannelMessageSend(data.ChannelID, "Hi there "+data.Author.Mention())
}

//respond to the roll command by sending a reply (containing random integer)
func rollCommand(s *discordgo.Session, data *CommandData) {
	args := data.Args
	num1 := 1
	num2 := 6
	if len(args) >= 1 {
		parsed1, err1 := strconv.Atoi(args[0])
		if err1 != nil {
			_, _ = s.ChannelMessageSend(data.ChannelID, "You cannot input a non-numeric value into this command. (Slot 1)")
			return
		} else if parsed1 > RAND_UPPER_LIM {
			s.ChannelMessageSend(data.ChannelID, "Your value was too great. (Slot 1)")
			return
		}
		num1 = parsed1
	}

	if len(args) == 2 {
		parsed2, err2 := strconv.Atoi(args[1])
		if err2 != nil {
			s.ChannelMessageSend(data.ChannelID, "You cannot input a non-numeric value into this command. (Slot 2)")
		} else if parsed2 > RAND_UPPER_LIM {
			s.ChannelMessageSend(data.ChannelID, "Your value was too great. (Slot 2)")
			return
		}
		num2 = parsed2
	}

	if num1 >= num2 {
		s.ChannelMessageSend(data.ChannelID, "You must input two values, the first one smaller than the other (num1>=num2)")
		return
	}

	s.ChannelMessageSend(data.ChannelID, "Your roll is: "+strconv.Itoa(rand.Intn(num2-num1)+num1))
}

//respond to the user asking the bot if it is online
func reportCommand(s *discordgo.Session, data *CommandData) {
	s.ChannelMessageSend(data.ChannelID, "Reporting for duty.")
}

//count the occurrences of words in the users message and send it back
func occurrencesCommand(s *discordgo.Session, data *CommandData) {
	occurMap := make(map[string]int)
	for _, w := range data.Args {
		occurMap[w] += 1
	}

	contentString := "Occurrences:\n```"
	for k, v := range occurMap {
		contentString += "(" + k + "," + strconv.Itoa(v) + ")\n"
	}
	s.ChannelMessageSend(data.ChannelID, contentString+"```")
}

//search the giphy library for either the top 3 trending gifs or the a random one of what the user requested
func giphySearchCommand(s *discordgo.Session, data *CommandData) {

	searchString := strings.Join(data.Args, " ")

	paginated := dgwidgets.NewPaginator(s, data.ChannelID)

	if searchString == "trending" {
		dataSearch, err := giphyLib.GetTrending()
		if err != nil {
			s.ChannelMessageSend(data.ChannelID, "There was an error while attempting a request to the Giphy Library.")
			return
		}
		s.ChannelMessageSend(data.ChannelID, "Here are my top "+strconv.Itoa(GIPHY_PRINT_LIM))
		printLen := GIPHY_PRINT_LIM
		if (len(dataSearch.Data)) < GIPHY_PRINT_LIM {
			printLen = len(dataSearch.Data)
		}

		for i := 0; i < printLen; i++ {
			embedded := embed.NewEmbed()

			embedded.SetTitle("Result " + strconv.Itoa(i+1))
			embedded.SetDescription(dataSearch.Data[i].Caption)
			urlEmbed := "https://media.giphy.com/media/" + dataSearch.Data[i].Id + "/giphy.gif"
			embedded.SetImage(urlEmbed)
			embedded.SetURL(urlEmbed)
			paginated.Add(embedded.MessageEmbed)
		}
		paginated.Spawn()
	} else {
		dataSearch, err := giphyLib.GetSearch(searchString, 5, -1, "", "", false)
		if err != nil {
			s.ChannelMessageSend(data.ChannelID, "There was an error while attempting a request to the Giphy Library.")
			return
		}
		s.ChannelMessageSend(data.ChannelID, "Here are my top "+strconv.Itoa(GIPHY_PRINT_LIM))
		printLen := GIPHY_PRINT_LIM
		if (len(dataSearch.Data)) < GIPHY_PRINT_LIM {
			printLen = len(dataSearch.Data)
		}

		for i := 0; i < printLen; i++ {
			embedded := embed.NewEmbed()

			embedded.SetTitle("Result " + strconv.Itoa(i+1))
			embedded.SetDescription(dataSearch.Data[i].Caption)
			urlEmbed := "https://media.giphy.com/media/" + dataSearch.Data[i].Id + "/giphy.gif"
			embedded.SetImage(urlEmbed)
			embedded.SetURL(urlEmbed)

			paginated.Add(embedded.MessageEmbed)
		}
		err = paginated.Spawn()
	}
}

func giftCommand(s *discordgo.Session, data *CommandData) {

	mentions := data.Message.Mentions

	//check if the arguments were only mentions
	if len(data.Args) != len(mentions) {
		s.ChannelMessageSend(data.ChannelID, "You must only mention members for this command.")
		return
	} else if len(data.Args) > GIFT_MENTION_LIM || len(data.Args) < 1 {
		s.ChannelMessageSend(data.ChannelID, "You input too many or too little arguments for this command.")
		return
	}

	mentionString := "A gift for "
	for _, v := range mentions {
		mentionString += v.Username + " "
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x008000, // Blue
		Description: "A sapphire is more precious than anything.",
		Fields:      []*discordgo.MessageEmbedField{},
		Image: &discordgo.MessageEmbedImage{
			URL: SAPPHIRE_URL,
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: data.Author.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
		Title:     mentionString,
	}
	s.ChannelMessageSendEmbed(data.ChannelID, embed)
}

//respond to the creating of message events by checking for input commands
func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	var content string
	content = m.Content

	args := strings.FieldsFunc(content, func(c rune) bool {
		return unicode.IsSpace(c)
	})
	fmt.Println(args)

	if len(args) == 0 || !strings.HasPrefix(args[0], "s/") {
		return
	}

	commandWord := strings.Split(args[0], "s/")[1]
	data := CommandData{args[1:], m.Message, m.Author, m.ChannelID}
	fmt.Println(data)

	if v, found := validMap[commandWord]; found {
		v.Executor.(func(*discordgo.Session, *CommandData))(s, &data)
	} else {
		s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
	}
}
