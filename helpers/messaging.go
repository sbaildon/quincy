package helpers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func PrivateMessage(session *discordgo.Session, user *discordgo.User, message string) {
        userChannel, err := session.UserChannelCreate(user.ID)
        if err != nil {
                fmt.Printf("Couldn't create private channel, err: " + err.Error())
                return
        }

        if userChannel.IsPrivate != true {
                fmt.Printf("Attempted to create private channel, but wasn't private")
        }

        session.ChannelMessageSend(userChannel.ID, message)
        session.ChannelDelete(userChannel.ID)
}

func DeleteMessage(session *discordgo.Session, message *discordgo.Message) {
	err := session.ChannelMessageDelete(message.ChannelID, message.ID)
	if err != nil {
		fmt.Printf("Unable to delete message")
	}
}
