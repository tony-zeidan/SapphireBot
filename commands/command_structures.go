package commands

import "github.com/bwmarrin/discordgo"

type CommandMapping struct {
	Name    string
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
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
	ValidCommandHandlers = []*CommandMapping{
		{
			Name:    "help",
			Handler: HelpCommand,
		},
		{
			Name:    "complete",
			Handler: CompleteCommand,
		},
		{
			Name:    "codehelp",
			Handler: CodeHelpCommand,
		},
	}
)
