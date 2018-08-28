package main

import (
	"github.com/bwmarrin/discordgo"
)

func StopAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel := (&Channel{}).FindChannel(m.ChannelID)
	if channel == nil {
		channel = &Channel{Channel: m.ChannelID, Status: OFF}
		channel.Save()
	}

	channel.Status = OFF
	channel.Save()
	s.ChannelMessageSend(m.ChannelID, "Уведомления отключены для этого канала")
}
