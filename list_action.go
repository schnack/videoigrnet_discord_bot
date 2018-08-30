package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ListAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel := (&Channel{}).FindByChannel(m.ChannelID)
	if channel == nil {
		s.ChannelMessageSend(m.ChannelID, "Нет отслеживаемых категорий")
		return
	}

	channelsCategories := (&ChannelsCategories{}).FindByChannel(channel)
	if len(channelsCategories) == 0 {
		s.ChannelMessageSend(m.ChannelID, "Нет отслеживаемых категорий")
		return
	}

	out := "Подключены следующие категории:\n\n"

	for i, cp := range channelsCategories {
		out = out + fmt.Sprintf("%d) %s | %s      [ добавил: %s ]\n", i+1, cp.Category.ParentName, cp.Category.Name, cp.Username)
	}

	s.ChannelMessageSend(m.ChannelID, out)
	return
}
