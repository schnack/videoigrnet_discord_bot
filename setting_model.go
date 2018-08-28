package main

import (
	"log"
	"time"
)

type Setting struct {
	Id    int
	Name  string
	Value string
}

func (s *Setting) Save() error {
	var id int
	err := DB.QueryRow("SELECT id FROM settings WHERE id = ?", s.Id).Scan(&id)
	if err != nil {
		_, err = DB.Exec("INSERT INTO settings(name, value, created_at, updated_at) VALUES(?, ?, ?, ?)", s.Name, s.Value, time.Now().Unix(), time.Now().Unix())
		if err != nil {
			return err
		}
	} else {
		_, err = DB.Exec("UPDATE settings SET name = ?, value = ?, updated_at = ? WHERE id = ?", s.Name, s.Value, time.Now().Unix(), s.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Setting) FindName(name string) *Setting {
	err := DB.QueryRow("SELECT id, value FROM settings WHERE name = ?", name).Scan(&s.Id, &s.Value)
	if err != nil {
		log.Fatalln(err)
	}
	return s
}
