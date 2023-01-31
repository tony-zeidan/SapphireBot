package commands

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
)

// respond to the user asking for help with the bots commands by sending a list of available commands
func CompleteCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := MakeOptionMapping(i)
	query := options["query"].StringValue()
	tokens := options["tokens"].IntValue()
	resp, err := GPT3Client.Completion(CTX, gpt3.CompletionRequest{
		Prompt:    []string{query},
		MaxTokens: gpt3.IntPtr(int(tokens)),
	})
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an issue with completing your input query.",
			},
		})
		if err != nil {
			return
		}
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
