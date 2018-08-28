package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func ListAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	channel := (&Channel{}).FindChannel(m.ChannelID)
	if channel == nil {
		s.ChannelMessageSend(m.ChannelID, "В текущем канале отсуствуют категории")
		return
	}

	channelsProducts := (&ChannelsProducts{}).FindChannels(channel)
	if len(channelsProducts) == 0 {
		s.ChannelMessageSend(m.ChannelID, "В текущем канале отсуствуют категории")
		return
	}

	out := "В текущем канале подключенные следующие категории:\n\n"

	for i, cp := range channelsProducts {
		out = out + fmt.Sprintf("%d) %s | %s      [ добавил: %s ]\n", i+1, cp.Product.CategoryParentName, cp.Product.CategoryName, cp.Username)
	}

	s.ChannelMessageSend(m.ChannelID, out)
	return
}
