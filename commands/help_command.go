package commands

import (
	embed "github.com/Clinet/discordgo-embed"
	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
	"github.com/tony-zeidan/SapphireBot"
	"strconv"
	"time"
)

// respond to the user asking for help with the bots commands by sending a list of available commands
func helpCommand(s *discordgo.Session, data *main.CommandData) {

	paginated := dgwidgets.NewPaginator(s, data.ChannelID)

	embedded := embed.NewEmbed()

	//we don't want the bot to print aliases (just first trigger)
	for i, v := range validCommands {

		if (i+1)%4 == 0 {
			paginated.Add(embedded.MessageEmbed)
			embedded = embed.NewEmbed()
		}
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

	paginated.NextPage()
}
