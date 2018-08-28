package main

import (
	"github.com/bwmarrin/discordgo"
)

func StartAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel := (&Channel{}).FindChannel(m.ChannelID)
	if channel == nil {
		channel = &Channel{Channel: m.ChannelID, Status: ON}
		channel.Save()
	}

	channel.Status = ON
	channel.Save()
	s.ChannelMessageSend(m.ChannelID, "Уведомления включены для этого канала")
}
