package hello

import (
	"fmt"

	"github.com/sbaildon/bot/commands"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.AddCommand("hello", &Command{})
}

type Command struct {
}

func (c *Command) Execute(session *discordgo.Session, message *discordgo.Message) {
	response := fmt.Sprintf("Hello %s", message.Author.Username)
	session.ChannelMessageSend(message.ChannelID, response)
}


