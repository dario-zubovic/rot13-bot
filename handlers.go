package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// invite link: https://discordapp.com/oauth2/authorize?client_id=478616184527388672&scope=bot&permissions=11264

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // ignore messages by bot itself
		return
	}

	if strings.HasPrefix(m.Message.Content, "!rot13 ") {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Println(err)
			return
		}

		str := fmt.Sprintf("%v: %v", m.Author.Mention(), doRot13(m.Message.Content[7:]))

		_, err = s.ChannelMessageSend(m.ChannelID, str)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.Emoji.Name != "spoilme" { // only react to spoiler emoji
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	if msg.Author.ID != s.State.User.ID { // we're only interested in messages sent by the bot
		return
	}

	ch, err := s.UserChannelCreate(m.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	i := strings.Index(msg.Content, ":")
	str := doRot13(msg.Content[i+2:])

	_, err = s.ChannelMessageSend(ch.ID, str)
	if err != nil {
		fmt.Println(err)
		return
	}
}
