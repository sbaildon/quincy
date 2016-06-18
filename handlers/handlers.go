package handlers

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/sbaildon/bot/commands"
)

func SetupHandlers(session *discordgo.Session) {
	session.AddHandler(onMessageCreate)
	session.AddHandler(onGuildCreate)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
	message := m.Message

	commands.HandleCommand(s, message)
}

func onGuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			s.ChannelMessageSend(channel.ID, "Yawk yawk yawk")
			return
		}
	}
}
