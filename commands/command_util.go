package commands

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
	"github.com/tony-zeidan/SapphireBot/config"
)

var (
	GPT3Client = gpt3.NewClient(config.Config.GPT3Token)
	CTX        = context.Background()
)

func MakeOptionMapping(i *discordgo.InteractionCreate) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range optionMap {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
