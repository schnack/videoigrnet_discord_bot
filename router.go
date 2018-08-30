package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

// router слушает команды и маршрутизирует запросы
func router(s *discordgo.Session, m *discordgo.MessageCreate) {

	switch {
	case m.Author.ID == s.State.User.ID:
		return

	case strings.HasPrefix(m.Content, "[vgnet add"):
		AddAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet list"):
		ListAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet del"):
		DelAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet start"):
		StartAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet stop"):
		StopAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet status"):
		StatusAction(s, m)
	case strings.HasPrefix(m.Content, "[vgnet"):
		s.ChannelMessageSend(m.ChannelID, "Добавление категории:\n\t[vgnet add <url>\n\n"+
			"Просмотр списка категорий:\n\t[vgnet list\n\n"+
			"Удаление категории:\n\t[vgnet del <num>\n\n"+
			"Запуск уведомлений:\n\t[vgnet start\n\n"+
			"Остановка уведомлений:\n\t[vgnet stop\n\n"+
			"Исходный код вы можете посмотреть по адресу https://github.com/schnack/videoigrnet_discord_bot")
	}
}
