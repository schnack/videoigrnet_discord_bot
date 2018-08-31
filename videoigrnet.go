package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
				dispatch[link.Channel.Channel] = "Лабудабудабтап!:mega:\n\n"
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
	return fmt.Sprintf(":fire: :fast_forward: %s | %s\n%s\n%s\n\nhttps://videoigr.net/product_info.php?products_id=%d\n\n", p.Category.ParentName, p.Category.Name, p.Name, GetPrice(p.Id), p.Id)
}

func formatMessageDel(p *Production) string {
	return fmt.Sprintf(":poop: :rewind: %s | %s\n%s\n\n", p.Category.ParentName, p.Category.Name, p.Name)
}

func GetPrice(id int64) string {
	uri := fmt.Sprintf("https://videoigr.net/product_info.php?products_id=%d", id)
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatalf("получение %s: %s", uri, resp.Status)
	}

	newread, _ := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))

	doc, err := html.Parse(newread)
	res := make([]string, 0)
	res = GetInfoProduct(res, doc)
	out := ""
	for _, tm := range res {
		out = out + tm
	}
	return out
}

func GetInfoProduct(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "script" && n.FirstChild != nil && strings.Contains(n.FirstChild.Data, "pa_self.push") {
		arrayParams := strings.Split(n.FirstChild.Data, ",")
		var special_name string
		var price string
		for _, keyValue := range arrayParams {
			if strings.Contains(keyValue, "special_name") {
				tmp := strings.Split(keyValue, ":")
				if len(tmp) > 1 {
					special_name = strings.Trim(tmp[1], "\"")
				}
			}
			if strings.Contains(keyValue, "price") {
				tmp := strings.Split(keyValue, ":")
				if len(tmp) > 1 {
					price = strings.Trim(tmp[1], "\"")
				}
			}
		}
		links = append(links, fmt.Sprintf("%s\t\t%s\n", special_name, price))
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = GetInfoProduct(links, c)
	}
	return links
}
