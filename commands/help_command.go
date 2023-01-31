package commands

import (
	embed "github.com/Clinet/discordgo-embed"
	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"time"
)

func HelpCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	HelpWrapper(s, i)
}

// respond to the user asking for help with the bots commands by sending a list of available commands
func HelpWrapper(s *discordgo.Session, i *discordgo.InteractionCreate) {

	paginated := dgwidgets.NewPaginator(s, i.ChannelID)

	embedded := embed.NewEmbed()

	//we don't want the bot to print aliases (just first trigger)
	paginated.Add(embedded.MessageEmbed)
	j := 1
	for i, v := range ValidCommands {
		if (i+1)%4 == 0 {
			j++
			paginated.Add(embedded.MessageEmbed)
			embedded = embed.NewEmbed()
		}
		embedded.SetTitle("Page " + strconv.Itoa(j))
		embedded.AddField(v.Name, v.Description)
	}

	// Sets the footers of all added pages to their page numbers.
	paginated.SetPageFooters()

	// When the paginator is done listening set the colour to yellow
	paginated.ColourWhenDone = 0xffff

	// Stop listening for reaction events after five minutes
	paginated.Widget.Timeout = time.Minute * 5

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "Please see a list of commands below!",
		},
	})
	if err != nil {
		_, err2 := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong! I'm sorry...",
		})
		if err2 != nil {
			return
		}
	}

	err = paginated.Spawn()
	if err != nil {
		log.Fatal(err)
	}
	paginated.NextPage()
}
