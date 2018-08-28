package main

import (
	"fmt"
	"log"
	"time"
)

type ChannelsProducts struct {
	Id        int
	Product   *Product
	Channel   *Channel
	Username  string
	CreatedAt int
	UpdatedAt int
}

func (cp *ChannelsProducts) Save() error {
	var id int
	err := DB.QueryRow("SELECT id FROM channels_products WHERE id = ?", cp.Id).Scan(&id)
	if err != nil {
		cp.UpdatedAt = int(time.Now().Unix())
		cp.CreatedAt = cp.UpdatedAt
		res, err := DB.Exec("INSERT INTO channels_products(product_id, channel_id, username, created_at, updated_at) VALUES(?, ?, ?, ?, ?)", cp.Product.Id, cp.Channel.Id, cp.Username, cp.CreatedAt, cp.UpdatedAt)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("Попытка получить id добавленой записи в таблице ChannelsProducts %s", err)
		}
		cp.Id = int(id)
	} else {
		cp.UpdatedAt = int(time.Now().Unix())
		_, err = DB.Exec("UPDATE channels_products SET product_id = ?, channel_id = ?, username = ?, updated_at = ? WHERE id = ?", cp.Product.Id, cp.Channel.Id, cp.Username, cp.UpdatedAt, cp.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cp *ChannelsProducts) Destroy() error {
	_, err := DB.Exec("DELETE FROM channels_products WHERE id = ?", cp.Id)
	if err != nil {
		return err
	}
	return nil
}

func (cp *ChannelsProducts) FindLink(c *Channel, p *Product) *ChannelsProducts {
	err := DB.QueryRow("SELECT id, username, created_at, updated_at FROM channels_products WHERE channel_id = ? AND product_id = ?", c.Id, p.Id).Scan(&cp.Id, &cp.Username, &cp.CreatedAt, &cp.UpdatedAt)
	if err != nil {
		return nil
	}
	cp.Product = p
	cp.Channel = c
	return cp
}

func (*ChannelsProducts) FindChannels(c *Channel) []*ChannelsProducts {
	channelsProducts := make([]*ChannelsProducts, 0)
	rows, err := DB.Query("SELECT id, product_id, username, created_at, updated_at FROM channels_products WHERE channel_id = ?", c.Id)
	if err != nil {
		return channelsProducts
	}
	defer rows.Close()
	for rows.Next() {
		cp_tmp := ChannelsProducts{}
		var productId int
		err := rows.Scan(&cp_tmp.Id, &productId, &cp_tmp.Username, &cp_tmp.CreatedAt, &cp_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		product := (&Product{}).FindId(productId)
		cp_tmp.Product = product
		cp_tmp.Channel = c
		channelsProducts = append(channelsProducts, &cp_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return channelsProducts
}

func (*ChannelsProducts) FindProducts(p *Product) []*ChannelsProducts {
	channelsProducts := make([]*ChannelsProducts, 0)
	rows, err := DB.Query("SELECT id, product_id, username, created_at, updated_at FROM channels_products WHERE product_id = ?", p.Id)
	if err != nil {
		return channelsProducts
	}
	defer rows.Close()
	for rows.Next() {
		cp_tmp := ChannelsProducts{}
		var channelId int
		err := rows.Scan(&cp_tmp.Id, &channelId, &cp_tmp.Username, &cp_tmp.CreatedAt, &cp_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		c := (&Channel{}).FindId(channelId)
		cp_tmp.Product = p
		cp_tmp.Channel = c
		channelsProducts = append(channelsProducts, &cp_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return channelsProducts
}
