package main

import (
	"fmt"
	"log"
	"time"
)

const (
	ON int = iota
	OFF
)

type Channel struct {
	Id        int
	Channel   string
	Status    int
	CreatedAt int
	UpdatedAt int
}

func (c *Channel) Save() error {
	var id int
	err := DB.QueryRow("SELECT id FROM channels WHERE id = ?", c.Id).Scan(&id)
	if err != nil {
		c.Status = OFF
		c.UpdatedAt = int(time.Now().Unix())
		c.CreatedAt = c.UpdatedAt
		res, err := DB.Exec("INSERT INTO channels(channel, status, created_at, updated_at) VALUES(?, ?, ?, ?)", c.Channel, c.Status, c.CreatedAt, c.UpdatedAt)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("Попытка получить id добавленой записи %s", err)
		}
		c.Id = int(id)
	} else {
		c.UpdatedAt = int(time.Now().Unix())
		_, err = DB.Exec("UPDATE channels SET channel = ?, status = ?, updated_at = ? WHERE id = ?", c.Channel, c.Status, c.UpdatedAt, c.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

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
