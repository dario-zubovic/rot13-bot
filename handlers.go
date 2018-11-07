package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// invite link: https://discordapp.com/oauth2/authorize?client_id={CLIENT-ID-HERE}&scope=bot&permissions=11264

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // ignore messages by bot itself
		return
	}

	var str string
	react := false

	if strings.HasPrefix(strings.ToLower(m.Message.Content), "!rot13 ") {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Println(err)
		} else {
			msgWithoutCmd := m.Message.Content[7:]
			if strings.ToLower(msgWithoutCmd) == "test" {
				str = "Yes, I am working fine."
			} else {
				str = fmt.Sprintf("%v: %v", m.Author.Mention(), doRot13(msgWithoutCmd))
			}
			react = true
		}
	}

	if len(str) == 0 {
		dm, err := comesFromDM(s, m)
		if err != nil {
			fmt.Println(err)
			return
		}

		if dm {
			str = doRot13(m.Message.Content)
		}
	}

	if len(str) > 0 {
		msg, err := s.ChannelMessageSend(m.ChannelID, str)
		if err != nil {
			fmt.Println(err)
			return
		}

		if react {
			err = s.MessageReactionAdd(msg.ChannelID, msg.ID, spoilEmojiID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func messageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID { // only react to reaction emojis posted by other users
		return
	}

	if m.Emoji.Name != spoilEmojiID { // only react to spoiler emoji
		return
	}

	msg, err := s.ChannelMessage(m.ChannelID, m.MessageID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// if msg.Author.ID != s.State.User.ID { // we're only interested in messages sent by the bot
	// 	return
	// }

	ch, err := s.UserChannelCreate(m.UserID)
	if err != nil {
		fmt.Println(err)
		return
	}

	var message, username string

	if msg.Author.ID == s.State.User.ID && msg.Content != "Yes, I am working fine." { // message was posted by bot
		i := strings.Index(msg.Content, ":")

		user, err := s.User(msg.Content[2 : i-1])
		if err == nil {
			username = user.Username
			message = msg.Content[i+2:]
		}
	}

	if len(message) == 0 {
		username = msg.Author.Username
		message = msg.Content
	}

	str := fmt.Sprintf("%v: %v", username, doRot13(message))

	_, err = s.ChannelMessageSend(ch.ID, str)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func comesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}
