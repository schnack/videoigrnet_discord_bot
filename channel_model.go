package main

import (
	"fmt"
	"time"
)

const (
	ON int64 = iota
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
	if nc != nil && c.Channel == nc.Channel && c.Status == nc.Status {
		return true
	}
	return false
}

// Save сохраняет состояние структуры в базу данных
func (c *Channel) Save() error {
	var err error
	tmp := (&Channel{}).FindById(c.Id)
	if tmp != nil {
		if c.Equal(tmp) {
			return nil
		} else {
			c.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE channels SET channel = ?, status = ?, updated_at = ? WHERE id = ?", c.Channel, c.Status, c.UpdatedAt, c.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в channels %s", err)
			}
		}
	} else {
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
func (c *Channel) FindById(v int64) *Channel {
	err := DB.QueryRow("SELECT id, channel, status, created_at, updated_at FROM channels WHERE id = ?", v).Scan(&c.Id, &c.Channel, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}

// FindByChannel поиск по внешнему id канала
func (c *Channel) FindByChannel(channel string) *Channel {
	err := DB.QueryRow("SELECT id, channel, status, created_at, updated_at FROM channels WHERE channel = ?", channel).Scan(&c.Id, &c.Channel, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}
