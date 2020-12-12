package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type SCommand struct {
	ContentRaw string
	Channel    string
	Author     *discordgo.User
}

func init() {

}
