package main

import (
	"fmt"
	"time"
)

type Category struct {
	Id         int64
	Name       string
	ParentId   int64
	ParentName string
	CreatedAt  int64
	UpdatedAt  int64
}

// Equal проверяет структуры на идентичность
func (c *Category) Equal(nc *Category) bool {
	return c.Id == nc.Id && c.Name == nc.Name && c.ParentId == nc.ParentId && c.ParentName == nc.ParentName
}

// Save сохраняет состояние структуры в базу данных
func (c *Category) Save() error {
	var err error
	tmp := (&Category{}).FindById(c.Id)
	if tmp != nil {
		if c.Equal(tmp) {
			return nil
		} else {
			c.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE category SET name = ?, parent_id = ?, parent_name = ?, updated_at = ? WHERE id = ?", c.Name, c.ParentId, c.ParentName, c.UpdatedAt, c.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в category %s", err)
			}
		}
	} else {
		c.UpdatedAt = time.Now().Unix()
		c.CreatedAt = c.UpdatedAt
		res, err := DB.Exec("INSERT INTO category(id, name, parent_id, parent_name, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?)", c.Id, c.Name, c.ParentId, c.ParentName, c.CreatedAt, c.UpdatedAt)
		if err != nil {
			return fmt.Errorf("ошика добавления строки в category %s", err)
		}
		c.Id, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("ошибка получения id новой записи category %s", err)
		}
	}
	return nil
}

// FindBy поиск по id
func (*Category) FindById(v int64) *Category {
	c := &Category{}
	err := DB.QueryRow("SELECT id, name, parent_id, parent_name, created_at, updated_at, FROM category WHERE id = ?", v).Scan(&c.Id, &c.Name, &c.ParentId, &c.ParentName, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil
	}
	return c
}
