package main

import (
	"fmt"
	"log"
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
	if ncc != nil && cc.Category.Equal(ncc.Category) && cc.Channel.Equal(ncc.Channel) {
		return true
	}
	return false
}

// Save сохраняет состояние структуры в базу данных
func (cc *ChannelsCategories) Save() error {
	var err error
	tmp := (&ChannelsCategories{}).FindByChannelCategory(cc.Channel, cc.Category)
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
func (cc *ChannelsCategories) FindByChannelCategory(ch *Channel, ca *Category) *ChannelsCategories {
	err := DB.QueryRow("SELECT id, username, created_at, updated_at FROM channels_categories WHERE category_id = ? AND channel_id = ?", ca.Id, ch.Id).Scan(&cc.Id, &cc.Username, &cc.CreatedAt, &cc.UpdatedAt)
	if err != nil {
		return nil
	}
	cc.Category = ca
	cc.Channel = ch
	return cc
}

// FindBy поиск по id
func (cc *ChannelsCategories) FindById(v int64) *ChannelsCategories {
	var caId, chId int64
	err := DB.QueryRow("SELECT id, category_id, channel_id, username, created_at, updated_at FROM channels_categories WHERE id = ?", v).Scan(&cc.Id, &caId, &chId, &cc.Username, &cc.CreatedAt, &cc.UpdatedAt)
	if err != nil {
		return nil
	}
	cc.Category = (&Category{}).FindById(caId)
	cc.Channel = (&Channel{}).FindById(chId)
	return cc
}

// FindByCategory поиск все каналов по категории
func (*ChannelsCategories) FindByCategory(c *Category) []*ChannelsCategories {
	cc := make([]*ChannelsCategories, 0)
	if c == nil {
		return cc
	}
	rows, err := DB.Query("SELECT id, category_id, channel_id, username, created_at, updated_at FROM channels_categories WHERE category_id = ? ", c.Id)
	if err != nil {
		log.Printf("Не удалось найти связи channels_categories %s", err)
		return cc
	}
	defer rows.Close()
	for rows.Next() {
		p_tmp := ChannelsCategories{}
		var caId int64
		var chId int64
		err := rows.Scan(&p_tmp.Id, &caId, &chId, &p_tmp.Username, &p_tmp.CreatedAt, &p_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		p_tmp.Category = (&Category{}).FindById(caId)
		p_tmp.Channel = (&Channel{}).FindById(chId)
		cc = append(cc, &p_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("ошибка получения данных из таблицы channels_categories: %s", err)
	}

	return cc
}

// FindByChannel поиск всех категорий по каналау
func (*ChannelsCategories) FindByChannel(c *Channel) []*ChannelsCategories {
	cc := make([]*ChannelsCategories, 0)
	if c == nil {
		return cc
	}
	rows, err := DB.Query("SELECT id, category_id, channel_id, username, created_at, updated_at FROM channels_categories WHERE channel_id = ? ", c.Id)
	if err != nil {
		log.Printf("Не удалось найти связи channels_categories %s", err)
		return cc
	}
	defer rows.Close()
	for rows.Next() {
		p_tmp := ChannelsCategories{}
		var caId int64
		var chId int64
		err := rows.Scan(&p_tmp.Id, &caId, &chId, &p_tmp.Username, &p_tmp.CreatedAt, &p_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		p_tmp.Category = (&Category{}).FindById(caId)
		p_tmp.Channel = (&Channel{}).FindById(chId)
		cc = append(cc, &p_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("ошибка получения данных из таблицы channels_categories: %s", err)
	}

	return cc
}

// Destroy уничтожение связи
func (cp *ChannelsCategories) Destroy() error {
	_, err := DB.Exec("DELETE FROM channels_categories WHERE id = ?", cp.Id)
	if err != nil {
		return err
	}
	return nil
}
