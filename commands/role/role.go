package rolte

import (
	"strings"
	"fmt"

	"github.com/sbaildon/quincy/commands"
	"github.com/sbaildon/quincy/helpers"
	"github.com/bwmarrin/discordgo"
	"github.com/renstrom/fuzzysearch/fuzzy"
)

func init() {
	commands.AddCommand("role", &Command{})
}

type Command struct {
}

func (c *Command) Execute(session *discordgo.Session, message *discordgo.Message) {
	processSubCommand(session, message)
}

func processSubCommand(session *discordgo.Session, message *discordgo.Message) {
	input := strings.Split(message.Content, " ")

	if len(input) < 2 {
		usage(session, message)
		fmt.Println("incorrect usage of role command")
		return
	}

	subcommand := input[1]
	subcommand = strings.ToLower(subcommand)

	switch subcommand {
		case "list": listRoles(session, message)
		case "assign": assign(session, message)
		case "unassign": unassign(session, message)
		case "create": createRole(session, message)
		case "delete": deleteRole(session, message)
		default: usage(session, message)
	}
}

func usage(session *discordgo.Session, message *discordgo.Message) {
	response := "!role [list, assign, unassign, create, delete] [name]"
	helpers.PrivateMessage(session, message.Author, response)
}

func listRoles(session *discordgo.Session, message *discordgo.Message) {
	channelId := message.ChannelID
	channel, _ := session.Channel(channelId)
	roles, _ := session.GuildRoles(channel.GuildID)

	response := "Available roles are ["
	for _, role := range roles {
		response = response + role.Name + ", "
	}
	response = response[:len(response)-2] + "]"

	helpers.PrivateMessage(session, message.Author, response)
	helpers.DeleteMessage(session, message)
}

func assign(session *discordgo.Session, message *discordgo.Message) {
	input := strings.Split(message.Content, " ")
	games := input[2:]

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
	helpers.DeleteMessage(session, message)
}

func unassign(session *discordgo.Session, message *discordgo.Message) {
	input := strings.Split(message.Content, " ")
	games := input[2:]

	channelId := message.ChannelID
	userId := message.Author.ID
	channel, _ := session.Channel(channelId)
	roles, _ := session.GuildRoles(channel.GuildID)

	var roles_s []string
	for _, role := range roles {
		roles_s = append(roles_s, role.Name)
		fmt.Printf("%s\n", role.Name)
	}

	member, _ := session.GuildMember(channel.GuildID, userId)
	memb_roles := member.Roles

	var rolesForRemoval []string
	for _, role := range roles {
		for _, game := range games {
			if role.Name == game {
				rolesForRemoval = append(rolesForRemoval, role.ID)
			}
		}
	}

	for i := 0; i < len(memb_roles); i++ {
		for _, role := range rolesForRemoval {
			if memb_roles[i] == role {
				memb_roles = append(memb_roles[:i], memb_roles[i+1:]...)
			}
		}
	}

	session.GuildMemberEdit(channel.GuildID, userId, memb_roles)
	helpers.DeleteMessage(session, message)
}

func createRole(session *discordgo.Session, message *discordgo.Message) {
	channelId := message.ChannelID
	channel, _ := session.Channel(channelId)

	newGuildRole, err := session.GuildRoleCreate(channel.GuildID)

	if err != nil {
		fmt.Println("error creating new role")
	} else {
		fmt.Println("created new role")
	}


	input := strings.Split(message.Content, " ")
	roleName := input[2]

	_, err = session.GuildRoleEdit(channel.GuildID, newGuildRole.ID, roleName, 0, false, 0)

	if err != nil {
		fmt.Println("error renaming role")
	} else {
		fmt.Println("renamed new role")
	}
}

func deleteRole(session *discordgo.Session, message *discordgo.Message) {
	channelId := message.ChannelID
	channel, _ := session.Channel(channelId)

	input := strings.Split(message.Content, " ")
	roleName := input[2]

	roles, _ := session.GuildRoles(channel.GuildID)

	for _, role := range roles {
		if role.Name == roleName {
			err := session.GuildRoleDelete(channel.GuildID, role.ID)
			if err != nil {
				fmt.Println("error deleting role")
			} else  {
				fmt.Println("deleted role")
			}
		}
	}
}

