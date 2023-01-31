package commands

import "github.com/bwmarrin/discordgo"

var (
	validCommands []CommandMapping
	validMap      map[string]CommandMapping
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

func init() {
	validMap = make(map[string]CommandMapping)

	validCommands = []CommandMapping{
		{
			Triggers:    []string{"help", "info"},
			Description: "Obtain information about Sapphire's commands.",
			Syntax:      "s/help",
			Executor:    helpCommand},
	}

	for _, v := range validCommands {
		for _, v2 := range v.Triggers {
			validMap[v2] = v
		}
	}
}
