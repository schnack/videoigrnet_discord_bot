package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func DelAction(s *discordgo.Session, m *discordgo.MessageCreate) {

	deleteCategoryId, err := strconv.ParseInt(strings.Trim(m.Content, "[vgnet del "), 10, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Не удалось удалить категорию. В качестве параметра принимается только целое число. Например\n\t [vgnet del 1")
	}

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

	if deleteCategoryId > 0 && len(channelsProducts) >= int(deleteCategoryId) {
		cp := channelsProducts[deleteCategoryId-1]
		err := cp.Destroy()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ошибка! Не удалось удалить категорию: \"%s | %s\"", cp.Product.CategoryParentName, cp.Product.CategoryName))
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Категория: \"%s | %s\" удалена из этого канала", cp.Product.CategoryParentName, cp.Product.CategoryName))
		return
	}

	out := "Вы ввели не верный номер категории. Попробуйте снова.\n\n"

	for i, cp := range channelsProducts {
		out = out + fmt.Sprintf("%d) %s | %s      [ добавил: %s ]\n", i+1, cp.Product.CategoryParentName, cp.Product.CategoryName, cp.Username)
	}

	s.ChannelMessageSend(m.ChannelID, out)
	return

}
