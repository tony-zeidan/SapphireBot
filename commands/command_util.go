package commands

import (
	"github.com/PullRequestInc/go-gpt3"
	"github.com/tony-zeidan/SapphireBot"
)

var (
	gptClient = gpt3.NewClient(main.Config.GPT3Token)
)
