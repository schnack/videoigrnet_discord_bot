package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func scanVideoigrNet(done <-chan struct{}) {
	log.Println("Сканирую videoigr.net")
	updateDB()
	for {
		syncTimeout, err := strconv.ParseInt((&Setting{}).FindName("sync_timeout").Value, 10, 0)
		if err != nil {
			syncTimeout = 3600
		}

		select {
		case <-time.After(time.Second * time.Duration(syncTimeout)):
			log.Println("Сканирую videoigr.net")
			updateDB()
		case <-done:
			log.Println("Завершаем работу синхронизации")
			return
		}
	}
}

func updateDB() {
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
	notify(DELETE)
	notify(NEW)
}

func notify(productStatus int) {

	dispatch := make(map[string]string)
	products := (&Product{}).FindStatus(productStatus)
	if len(products) == 0 {
		return
	}
	for _, p := range products {
		chanelsProducts := (&ChannelsProducts{}).FindProducts(p)
		if len(chanelsProducts) == 0 {
			return
		}
		for _, cp := range chanelsProducts {
			if _, ok := dispatch[cp.Channel.Channel]; !ok {
				if productStatus == NEW {
					dispatch[cp.Channel.Channel] = "Появились новые игры в отслеживаемом разделе:\n\n"
				} else {
					dispatch[cp.Channel.Channel] = "Распроданные игры:\n\n"
				}
			}
			dispatch[cp.Channel.Channel] = dispatch[cp.Channel.Channel] + fmt.Sprintf("https://videoigr.net/product_info.php?products_id=%d\n\n", cp.Product.Id)
		}
	}

	for ch, mess := range dispatch {
		DG.ChannelMessageSend(ch, mess)
	}

}
