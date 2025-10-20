package handler

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/walkure/discord-unfurler/twitter"
	"github.com/walkure/discord-unfurler/util"
)

func HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("Message received:", m.Content)

	targets := util.ExtractHTTPSURLs(m.Content, []string{"x.com", "twitter.com"})

	for _, target := range targets {
		c, err := twitter.HandleExpandContent(target)
		if err != nil {
			fmt.Printf("expand err[%q]:%v\n", target, err)
			continue
		}
		_, err = s.ChannelMessageSendComplex(m.ChannelID, c)
		if err != nil {
			fmt.Printf("send err[%q]:%v\n", target, err)
		}
	}

}
