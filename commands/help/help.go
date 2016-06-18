package list

import (
	"fmt"
	"strings"

	"github.com/sbaildon/quincy/commands"
	"github.com/sbaildon/quincy/helpers"
	"github.com/bwmarrin/discordgo"
)

func init() {
	commands.AddCommand("help", &Command{})
}

type Command struct {
}

func (c *Command) Execute(session *discordgo.Session, message *discordgo.Message) {
	names := commands.CommandNames()

	response := fmt.Sprintf("Available commands include: %s\n", strings.Join(names, ", "))

	helpers.PrivateMessage(session, message.Author, response)
}


