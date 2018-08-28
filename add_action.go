package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

func AddAction(s *discordgo.Session, m *discordgo.MessageCreate) {
	categories := strings.Split(strings.Trim(m.Content, "[vgnet add https://videoigr.net/index.php?cPath="), "_")
	if len(categories) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Плохая ссылка. В ссылке должен обязательно присуствовать параметр cPath")
		return
	}

	category, err := strconv.ParseInt(categories[1], 10, 0)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Плохая ссылка. В ссылке должен обязательно присуствовать параметр cPath=34_34")
		return
	}

	channel := (&Channel{}).FindChannel(m.ChannelID)
	if channel == nil {
		channel = &Channel{Channel: m.ChannelID, Status: OFF}
		channel.Save()
	}

	product := (&Product{}).FindCategory(int(category))
	if product == nil {
		s.ChannelMessageSend(m.ChannelID, "Указанной категории не существует проверте ссылку")
		return
	}

	channelsProdutcs := (&ChannelsProducts{}).FindLink(channel, product)
	if channelsProdutcs == nil {
		channelsProdutcs = &ChannelsProducts{Channel: channel, Product: product, Username: m.Author.Username}
		err := channelsProdutcs.Save()
		if err != nil {
			log.Fatal(err)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("В текущий канал добавлена новая категория: %s | %s", product.CategoryParentName, product.CategoryName))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("В текущем канале уже существует категория: %s | %s", product.CategoryParentName, product.CategoryName))
	}
}
