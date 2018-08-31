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

type ProductImport struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	CategoryId         string `json:"cat_id"`
	CategoryName       string `json:"cat_name"`
	CategoryParentId   string `json:"cat_parent_id"`
	CategoryParentName string `json:"cat_parent_name"`
	BuyStatus          string `json:"buy_status"`
}

func (pi *ProductImport) Import() *Production {
	Id, err := strconv.ParseInt(pi.Id, 10, 0)
	if err != nil {
		log.Printf("Не удалось конвертировать Id: %v", pi.Id)
		return nil
	}
	CategoryId, err := strconv.ParseInt(pi.CategoryId, 10, 0)
	if err != nil {
		log.Printf("Не удалось конвертировать CategoryId: %v", pi.Id)
		return nil
	}
	CategoryParentId, err := strconv.ParseInt(pi.CategoryParentId, 10, 0)
	if err != nil {
		log.Printf("Не удалось конвертировать CategoryParentId: %v", pi.CategoryParentId)
		return nil
	}
	BuyStatus, err := strconv.ParseInt(pi.BuyStatus, 10, 0)
	if err != nil {
		log.Printf("Не удалось конвертировать BuyStatus: %v", pi.BuyStatus)
		return nil
	}
	category := &Category{Id: CategoryId, Name: pi.CategoryName, ParentId: CategoryParentId, ParentName: pi.CategoryParentName}
	category.Save()
	production := &Production{Id: Id, Name: pi.Name, Category: category, BuyStatus: BuyStatus}
	production.Save()
	return production
}

func scanVideoigrNet(done <-chan struct{}) {
	for {
		select {
		case <-time.After(time.Second * 300):
			log.Println("Сканирую videoigr.net")
			updateDB()
			log.Println("Обновление БД завершено")
			notify()
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

	(&Production{}).UpdateAllStatusDel()
	for _, p := range result {
		p.Import()
	}
}

func notify() {
	dispatch := make(map[string]string)

	productions := (&Production{}).FindByStatusNewDel()

	for _, product := range productions {
		links := product.Category.FindChannels()
		for _, link := range links {
			// Пропускаем каналы на которых отключены уведомление
			if link.Channel.Status == OFF {
				continue
			}
			// Собираем сообщения
			if _, ok := dispatch[link.Channel.Channel]; !ok {
				dispatch[link.Channel.Channel] = "Лабудабудабтап!:mega: "
			}

			if product.Status == NEW {
				dispatch[link.Channel.Channel] = dispatch[link.Channel.Channel] + formatMessageNew(product)
			} else {
				dispatch[link.Channel.Channel] = dispatch[link.Channel.Channel] + formatMessageDel(product)
			}

		}
	}

	for ch, mess := range dispatch {
		log.Println(ch, mess)
		DG.ChannelMessageSend(ch, mess)
	}
}

func formatMessageNew(p *Production) string {
	return fmt.Sprintf(":fire: :heavy_plus_sign:%s | %s\n%s\nhttps://videoigr.net/product_info.php?products_id=%d\n\n", p.Category.ParentName, p.Category.Name, p.Name, p.Id)
}

func formatMessageDel(p *Production) string {
	return fmt.Sprintf(":poop: :heavy_minus_sign: %s | %s\n%s\n\n", p.Category.ParentName, p.Category.Name, p.Name)
}
