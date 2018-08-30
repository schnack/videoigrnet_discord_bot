package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func AddAction(s *discordgo.Session, m *discordgo.MessageCreate) {
	categories := strings.Split(strings.Trim(m.Content, "[vgnet add https://videoigr.net/index.php?cPath="), "_")
	if len(categories) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Плохая ссылка. В ссылке должен обязательно присутствовать параметр cPath")
		return
	}

	categoryId, err := strconv.ParseInt(categories[1], 10, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Плохая ссылка. В ссылке должен обязательно присутствовать параметр cPath=34_34")
		return
	}

	channel := (&Channel{}).FindByChannel(m.ChannelID)
	if channel == nil {
		channel = &Channel{Channel: m.ChannelID, Status: OFF}
		channel.Save()
	}

	category := (&Category{}).FindById(categoryId)
	if category == nil {
		s.ChannelMessageSend(m.ChannelID, "Указанной категории не существует")
		return
	}

	dup := (&ChannelsCategories{}).FindByChannelCategory(channel, category)
	if dup != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("категория \"%s | %s\" уже подключена в этом канале", category.ParentName, category.Name))
		return
	}

	channelsCategories := ChannelsCategories{Category: category, Channel: channel, Username: m.Author.Username}
	err = channelsCategories.Save()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Ошибка при добавлении новой категории \"%s | %s\" в этот канал %s", category.ParentName, category.Name, err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("В текущий канал добавлена новая категория: %s | %s", category.ParentName, category.Name))
}
