package main

import (
	"fmt"
	"time"
)

type ChannelsCategories struct {
	Id        int64
	Category  *Category
	Channel   *Channel
	Username  string
	CreatedAt int64
	UpdatedAt int64
}

// Equal проверяет структуры на идентичность
func (cc *ChannelsCategories) Equal(ncc *ChannelsCategories) bool {
	return cc.Category.Equal(ncc.Category) && cc.Channel.Equal(ncc.Channel) && cc.Username == ncc.Username
}

// Save сохраняет состояние структуры в базу данных
func (cc *ChannelsCategories) Save() error {
	var err error
	tmp := (&ChannelsCategories{}).FindById(cc.Id)
	if tmp != nil {
		if cc.Equal(tmp) {
			return nil
		} else {
			cc.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE channels_categories SET category_id = ?, channel_id = ?, username = ?,  updated_at = ? WHERE id = ?", cc.Category.Id, cc.Channel.Id, cc.Username, cc.UpdatedAt, cc.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в channels_categories %s", err)
			}
		}
	} else {
		cc.UpdatedAt = time.Now().Unix()
		cc.CreatedAt = cc.UpdatedAt
		res, err := DB.Exec("INSERT INTO channels_categories(category_id, channel_id, username, created_at, updated_at) VALUES(?, ?, ?, ?, ?)", cc.Category.Id, cc.Channel.Id, cc.Username, cc.CreatedAt, cc.UpdatedAt)
		if err != nil {
			return fmt.Errorf("ошика добавления строки в channels_categories %s", err)
		}
		cc.Id, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("ошибка получения id новой записи channels_categories %s", err)
		}
	}
	return nil
}

// FindBy поиск по id
func (*ChannelsCategories) FindById(v int64) *ChannelsCategories {
	cc := &ChannelsCategories{}
	var caId, chId int64
	err := DB.QueryRow("SELECT id, category_id, channel_id, username, created_at, updated_at FROM channels_categories WHERE id = ?", v).Scan(&cc.Id, &caId, &chId, &cc.Username, &cc.CreatedAt, &cc.UpdatedAt)
	if err != nil {
		return nil
	}
	cc.Category = (&Category{}).FindById(caId)
	cc.Channel = (&Channel{}).FindById(chId)
	return cc
}

/*
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
*/
