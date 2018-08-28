package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func scanVideoigrNet(done <-chan struct{}, dg *discordgo.Session) {
	log.Println("Сканирую videoigr.net")
	updateDB(dg)
	for {
		syncTimeout, err := strconv.ParseInt((&Setting{}).FindName("sync_timeout").Value, 10, 0)
		if err != nil {
			syncTimeout = 3600
		}

		select {
		case <-time.After(time.Second * time.Duration(syncTimeout)):
			log.Println("Сканирую videoigr.net")
			updateDB(dg)
		case <-done:
			log.Println("Завершаем работу синхронизации")
			return
		}
	}
}

func updateDB(dg *discordgo.Session) {
	uri := "https://videoigr.net/matrix.php"
	resp, err := http.PostForm(uri, url.Values{})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Printf("получение %s: %s", uri, resp.Status)
	}

	var result []ProductImport
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
	}

	MarkProductDelete()
	for _, p := range result {
		p.Conv().Save()
	}
	notify(dg)
}

func notify(dg *discordgo.Session) {
	dispatch := make(map[string]string)
	channel := (&Channel{}).FindAllOn()
	if len(channel) == 0 {
		log.Println("В каналах отключены уведомления")
		return
	}

	for _, c := range channel {
		chanelsProducts := (&ChannelsProducts{}).FindChannels(c)
		for _, cp := range chanelsProducts {
			products := (&Product{}).FindStatusCategory(EXIST, cp.Product.CategoryId)
			if len(products) == 0 {
				log.Println("Нет объектов для уведомления")
				return
			}
			for _, p := range products {
				if _, ok := dispatch[cp.Channel.Channel]; !ok {
					if p.Status == NEW {
						dispatch[cp.Channel.Channel] = "Появились новые игры в отслеживаемом разделе:\n\n"
					} else {
						dispatch[cp.Channel.Channel] = "Распроданные игры:\n\n"
					}
				}
				dispatch[cp.Channel.Channel] = dispatch[cp.Channel.Channel] + fmt.Sprintf("https://videoigr.net/product_info.php?products_id=%d\n\n", p.Id)
			}
		}

	}

	for ch, mess := range dispatch {
		log.Println(ch, mess)
		dg.ChannelMessageSend("482526981049679892", mess)
	}

}
