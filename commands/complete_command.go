package commands

import (
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
	"log"
)

// respond to the user asking for help with the bots commands by sending a list of available commands
func CompleteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := MakeOptionMapping(i)
	fmt.Println(options)
	fmt.Println(i.ApplicationCommandData().Options[0])
	query := options["query"].StringValue()

	var tokens = 1000
	if _, ok := options["tokens"]; ok {
		tokens = int(options["tokens"].IntValue())
	}

	resp, err := GPT3Client.Completion(CTX, gpt3.CompletionRequest{
		Prompt:    []string{query},
		MaxTokens: gpt3.IntPtr(tokens),
	})
	log.Println(err)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an issue with completing your input query.",
			},
		})
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: resp.Choices[0].Text,
		},
	})
	if err != nil {
		_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Something went wrong! I'm sorry...",
		})
		if err != nil {
			return
		}
	}
}
