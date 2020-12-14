package main

import (
	"flag"
	"fmt"
	"github.com/sanzaru/go-giphy"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode"

	"github.com/bwmarrin/discordgo"
)

var (
	Token    string
	validMap map[string]interface{}
	giphyLib *libgiphy.Giphy
)

const (
	//Giphy API token
	GIPHY_API_TOKEN = "CXYNdpzCDL4y8XgaJqgWf75khRNc1goy"
	//Random number generation upper limit
	RAND_UPPER_LIM = 100000
	//Giphy number of images limit
	GIPHY_PRINT_LIM = 3
	//Sapphire gif (for gift command)
	SAPPHIRE_URL     = "https://assets.bigcartel.com/product_images/158847679/SAV-201V---75361.gif"
	GIFT_MENTION_LIM = 3
)

//Run once on initialization
func init() {
	flag.StringVar(&Token, "t", "NjcyNTkwMDE4MzUwNDE1ODc0.XjNsRA.Dr_CmP1J2DI0COuw3z23XNLlkgk", "Bot Token")
	flag.Parse()
	giphyLib = libgiphy.NewGiphy(GIPHY_API_TOKEN)
	validMap = make(map[string]interface{})
	validMap["hello"] = helloCommand
	validMap["roll"] = rollCommand
	validMap["status"] = reportCommand
	validMap["freq"] = occurrencesCommand
	validMap["help"] = helpCommand
	validMap["giphy"] = giphySearchCommand
	validMap["gift"] = giftCommand
}

func main() {
	dg, err := discordgo.New("Bot " + Token)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error creating Discord session for Sapphire Bot,", err)
		return
	}

	fmt.Println("Sapphire is now running.")

	//Ctrl + C to kill
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()

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

//respond to the user saying hello
func helloCommand(s *discordgo.Session, data *CommandData) {
	s.ChannelMessageSend(data.ChannelID, "Hi there "+data.Author.Mention())
}

//respond to the roll command by sending a reply (containing random integer)
func rollCommand(s *discordgo.Session, data *CommandData) {
	args := data.Args
	num1 := 1
	num2 := 6
	if len(args) >= 1 {
		parsed1, err1 := strconv.Atoi(args[0])
		if err1 != nil {
			s.ChannelMessageSend(data.ChannelID, "You cannot input a non-numeric value into this command. (Slot 1)")
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

//respond to the user asking for help with the bots commands by sending a list of available commands
func helpCommand(s *discordgo.Session, data *CommandData) {
	contentString := "list of commands:\n```"
	for k := range validMap {
		contentString += "\t-" + k + "\n"
	}
	s.ChannelMessageSend(data.ChannelID, contentString+"```")
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

	if searchString == "trending" {
		dataSearch, err := giphyLib.GetTrending()
		if err != nil {
			s.ChannelMessageSend(data.ChannelID, "There was an error while attempting a request to the Giphy Library.")
			return
		}
		s.ChannelMessageSend(data.ChannelID, "Here are my top 5")
		printLen := GIPHY_PRINT_LIM
		if (len(dataSearch.Data)) < GIPHY_PRINT_LIM {
			printLen = len(dataSearch.Data)
		}

		for i := 0; i < printLen; i++ {
			s.ChannelMessageSend(data.ChannelID, dataSearch.Data[i].Url)
		}

	} else {
		dataSearch, err := giphyLib.GetRandom(searchString)
		if err != nil {
			s.ChannelMessageSend(data.ChannelID, "There was an error while attempting a request to the Giphy Library.")
			return
		}
		s.ChannelMessageSend(data.ChannelID, dataSearch.Data.Url)
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
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		v.(func(*discordgo.Session, *CommandData))(s, &data)
	} else {
		s.ChannelMessageSend(m.ChannelID, "That was not a valid command.")
	}
}
