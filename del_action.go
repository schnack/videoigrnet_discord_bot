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

	if deleteCategoryId > 0 && len(channelsCategories) >= int(deleteCategoryId) {
		cp := channelsCategories[deleteCategoryId-1]
		err := cp.Destroy()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ошибка! Не удалось удалить категорию: \"%s | %s\"", cp.Category.ParentName, cp.Category.Name))
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Категория: \"%s | %s\" удалена", cp.Category.ParentName, cp.Category.Name))
		return
	}

	out := "Вы ввели не верный порядковый номер категории. Попробуйте снова.\n\n"

	for i, cp := range channelsCategories {
		out = out + fmt.Sprintf("%d) %s | %s      [ добавил: %s ]\n", i+1, cp.Category.ParentName, cp.Category.Name, cp.Username)
	}

	s.ChannelMessageSend(m.ChannelID, out)
	return
}
