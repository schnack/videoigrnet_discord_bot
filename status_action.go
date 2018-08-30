package main

import (
	"github.com/bwmarrin/discordgo"
)

func StatusAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel := (&Channel{}).FindByChannel(m.ChannelID)
	if channel == nil {
		channel = &Channel{Channel: m.ChannelID, Status: OFF}
		channel.Save()
	}
	if channel.Status == ON {
		s.ChannelMessageSend(m.ChannelID, "Уведомления включены для этого канала")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Уведомления отключены для этого канала")
	}

}
