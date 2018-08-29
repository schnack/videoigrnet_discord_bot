package main

import (
	"fmt"
	"time"
)

const (
	ON int = iota
	OFF
)

type Channel struct {
	Id        int64
	Channel   string
	Status    int64
	CreatedAt int64
	UpdatedAt int64
}

// Equal проверяет структуры на идентичность
func (c *Channel) Equal(nc *Channel) bool {
	return c.Channel == nc.Channel && c.Status == nc.Status
}

// Save сохраняет состояние структуры в базу данных
func (c *Channel) Save() error{
	var err error
	tmp := (&Channel{}).FindById(c.Id)
	if tmp != nil {
		if c.Equal(tmp) {
			return nil
		}else{
			c.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE channels SET channel = ?, status = ?, updated_at = ? WHERE id = ?", c.Channel, c.Status, c.UpdatedAt, c.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в channels %s", err)
			}
		}
	}else{
		c.UpdatedAt = time.Now().Unix()
		c.CreatedAt = c.UpdatedAt
		res, err := DB.Exec("INSERT INTO channels(channel, status, created_at, updated_at) VALUES(?, ?, ?, ?)", c.Channel, OFF, c.CreatedAt, c.UpdatedAt)
		if err != nil {
			return fmt.Errorf("ошика добавления строки в channels %s", err)
		}
		c.Id, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("ошибка получения id новой записи channels %s", err)
		}
	}
	return nil
}

// FindBy поиск по id
func (*Channel) FindById(v int64) *Channel {
	c := &Channel{}
	err := DB.QueryRow("SELECT id, channel, status, created_at, updated_at FROM channels WHERE id = ?", v).Scan(&c.Id, &c.Channel, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}
/*

func (c *Channel) FindChannel(channel string) *Channel {
	err := DB.QueryRow("SELECT id, channel, status, created_at, updated_at FROM channels WHERE channel = ?", channel).Scan(&c.Id, &c.Channel, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}

func (c *Channel) FindId(chanelId int) *Channel {
	err := DB.QueryRow("SELECT id, channel, status, created_at, updated_at FROM channels WHERE id = ?", chanelId).Scan(&c.Id, &c.Channel, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}

func (*Channel) FindAllOn() []*Channel {
	channels := make([]*Channel, 0)
	rows, err := DB.Query("SELECT id, channel, status, created_at, updated_at FROM channels WHERE status = ?", ON)
	if err != nil {
		return channels
	}
	defer rows.Close()
	for rows.Next() {
		c_tmp := Channel{}
		err := rows.Scan(&c_tmp.Id, &c_tmp.Channel, &c_tmp.Status, &c_tmp.CreatedAt, &c_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		channels = append(channels, &c_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return channels
}
/*