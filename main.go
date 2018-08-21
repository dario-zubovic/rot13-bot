package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	spoilEmojiID = "🍉"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please specify bot token in the first argument.")
		return
	}

	fmt.Println("Starting...")

	// create session
	session, err := discordgo.New("Bot " + os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	// register handlers
	session.AddHandler(messageCreate)
	session.AddHandler(messageReactionAdd)

	// start session
	err = session.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	// hang
	fmt.Println("CTRL+C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// close session
	session.Close()
}
