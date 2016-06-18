package roles

import (
	"strings"
	"fmt"

	"github.com/sbaildon/bot/commands"
	"github.com/sbaildon/bot/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

func init() {
	commands.AddCommand("roles", &Command{})
}

type Command struct {
}

func (c *Command) Execute(session *discordgo.Session, message *discordgo.Message) {
	input := strings.Split(message.Content, " ")
	games := input[1:]

	channelId := message.ChannelID
	userId := message.Author.ID
	channel, _ := session.Channel(channelId)
	roles, _ := session.GuildRoles(channel.GuildID)

	//session.GuildRoleCreate(channel.GuildID)
	var roles_s []string
	for _, role := range roles {
		roles_s = append(roles_s, role.Name)
		fmt.Printf("%s\n", role.Name)
	}

	member, _ := session.GuildMember(channel.GuildID, userId)
	memb_roles := member.Roles

	var unrecognised []string
	for _, game := range games {
		results := fuzzy.Find(game, roles_s)
		if len(results) == 1 {
			for _, role  := range roles {
				if results[0] == role.Name {
					fmt.Printf("adding %s\n", role.Name)
					memb_roles = append(memb_roles, role.ID)
				}
			}
		} else {
			unrecognised = append(unrecognised, game)
		}
	}

	if len(unrecognised) > 0 {
		response := "The roles ["
		for _,  un := range unrecognised {
			response = response + un + ", "
		}
		response = response[:len(response)-2] + "] were unrecognised"

		helpers.PrivateMessage(session, message.Author, response)
	}

	session.GuildMemberEdit(channel.GuildID, userId, memb_roles)
}


