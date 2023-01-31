package commands

import (
	"github.com/Clinet/discordgo-embed"
	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var GiphyToken string

const (
	//Random number generation upper limit
	RAND_UPPER_LIM = 100000
	//Giphy number of images limit
	GIPHY_PRINT_LIM = 10
	//Sapphire gif (for gift command)
	SAPPHIRE_URL     = "https://assets.bigcartel.com/product_images/158847679/SAV-201V---75361.gif"
	GIFT_MENTION_LIM = 3
)

// respond to the user saying hello
func helloCommand(s *discordgo.Session, data *CommandData) {
	_, _ = s.ChannelMessageSend(data.ChannelID, "Hi there "+data.Author.Mention())
}

// respond to the roll command by sending a reply (containing random integer)
func rollCommand(s *discordgo.Session, data *CommandData) {
	args := data.Args
	num1 := 1
	num2 := 10
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

// respond to the user asking the bot if it is online
func reportCommand(s *discordgo.Session, data *CommandData) {
	s.ChannelMessageSend(data.ChannelID, "Reporting for duty.")
}

// count the occurrences of words in the users message and send it back
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

func getSearch(s string) (interface{}, error) {
	if s == "trending" {
		return giphyLib.GetTrending()
	}
	return giphyLib.GetSearch(s, 5, -1, "", "", false)
}

// search the giphy library for either the top 3 trending gifs or the a random one of what the user requested
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

func testCommand(s *discordgo.Session, data *CommandData) {

	p := dgwidgets.NewPaginator(s, data.Message.ChannelID)

	// Add embed pages to paginator
	p.Add(&discordgo.MessageEmbed{Description: "Page one"},
		&discordgo.MessageEmbed{Description: "Page two"},
		&discordgo.MessageEmbed{Description: "Page three"})

	// Sets the footers of all added pages to their page numbers.
	p.SetPageFooters()

	// When the paginator is done listening set the colour to yellow
	p.ColourWhenDone = 0xffff

	// Stop listening for reaction events after five minutes
	p.Widget.Timeout = time.Minute * 5

	// Add a custom handler for the gun reaction.
	p.Widget.Handle("🔫", func(w *dgwidgets.Widget, r *discordgo.MessageReaction) {
		s.ChannelMessageSend(data.Message.ChannelID, "Bang!")
	})

	p.Spawn()
}

func pongCommand(s *discordgo.Session, data *CommandData) {
	s.ChannelMessageSend(data.ChannelID, "Pong!")
}
