package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/walkure/discord-unfurler/handler"
)

func main() {
	os.Exit(run())
}

func run() int {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	if botToken == "" {
		fmt.Println("DISCORD_BOT_TOKEN environment variable is not set")
		return -1
	}
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return -2
	}
	//dg.Debug = true

	// allow the bot to receive messages
	dg.AddHandler(handler.HandleMessageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return -3
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot Started")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()

	fmt.Println("Bot Exiting....")

	// Cleanly close down the Discord session.
	dg.Close()

	return 0
}
