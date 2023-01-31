package commands

import "github.com/bwmarrin/discordgo"

var (
	validCommands       []CommandMapping
	ValidCommandMapping map[string]CommandMapping
)

// CommandData Data structure for containing command data
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

type CommandMapping struct {
	Triggers    []string
	Description string
	Syntax      string
	Executor    interface{}
	SubCommands []CommandData
}

var (
	integerOptionMinValue = 1.0
	ValidCommands         = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Gives a basic list of commands for the Sapphire bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "pagenum",
					Description: "Specify the page number to display.",
					MinValue:    &integerOptionMinValue,
					Type:        discordgo.ApplicationCommandOptionInteger,
				},
			},
		},
		{
			Name:        "complete",
			Description: "Uses OpenAI GPT-3 to complete your phrase!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "query",
					Description: "Specify the text to complete.",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
		},
		{
			Name:        "codehelp",
			Description: "Uses OpenAi to help you with code!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "codefile",
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Description: "The code in attachment format.",
				},
				{
					Name:        "codetext",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "The code in text format.",
				},
			},
		},
	}
)
