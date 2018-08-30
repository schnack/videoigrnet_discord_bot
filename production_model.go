package main

import (
	"fmt"
	"log"
	"time"
)

const (
	NEW int64 = iota
	DELETE
	EXIST
)

type Production struct {
	Id        int64
	Name      string
	Category  *Category
	BuyStatus int64
	Status    int64
	CreatedAt int64
	UpdatedAt int64
}

// Equal проверяет структуры на идентичность
func (p *Production) Equal(np *Production) bool {
	if np != nil && p.Id == np.Id && p.Name == np.Name && p.Category.Equal(np.Category) && p.BuyStatus == np.BuyStatus {
		return true
	}
	return false
}

// Save сохраняет состояние структуры в базу данных
func (p *Production) Save() error {
	var err error
	tmp := (&Production{}).FindById(p.Id)
	if tmp != nil {
		if p.Equal(tmp) {
			p.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE productions SET status = ?, updated_at = ? WHERE id = ?", EXIST, p.UpdatedAt, p.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в productions %s", err)
			}
			return nil
		} else {
			p.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE productions SET name = ?, category_id = ?, buy_status = ?, status = ?, updated_at = ? WHERE id = ?", p.Name, p.Category.Id, p.BuyStatus, EXIST, p.UpdatedAt, p.Id)
			if err != nil {
				return fmt.Errorf("ошибка обновления строки в productions %s", err)
			}
		}
	} else {
		p.UpdatedAt = time.Now().Unix()
		p.CreatedAt = p.UpdatedAt
		res, err := DB.Exec("INSERT INTO productions(id, name, category_id, buy_status, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.Category.Id, p.BuyStatus, NEW, p.CreatedAt, p.UpdatedAt)
		if err != nil {
			return fmt.Errorf("ошибка добавления строки в productions %s", err)
		}
		p.Id, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("ошибка получения id новой записи productions %s", err)
		}
	}
	return nil
}

// FindBy поиск по id
func (p *Production) FindById(v int64) *Production {
	var cId int64
	err := DB.QueryRow("SELECT id, name, category_id, buy_status, status, created_at, updated_at FROM productions WHERE id = ?", v).Scan(&p.Id, &p.Name, &cId, &p.BuyStatus, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil
	}
	p.Category = (&Category{}).FindById(cId)
	return p
}

func (*Production) FindByStatusNewDelete() []*Production {
	products := make([]*Production, 0)
	rows, err := DB.Query("SELECT id, name, category_id, buy_status, status, created_at, updated_at FROM productions WHERE status = ? OR status = ?", NEW, DELETE)
	if err != nil {
		log.Println("Не удалось найти продукты со статусом ")
		return products
	}
	defer rows.Close()
	for rows.Next() {
		p_tmp := Production{}
		var cId int64
		err := rows.Scan(&p_tmp.Id, &p_tmp.Name, &cId, &p_tmp.BuyStatus, &p_tmp.Status, &p_tmp.CreatedAt, &p_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		p_tmp.Category = (&Category{}).FindById(cId)
		products = append(products, &p_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Printf("ошибка получения данных из таблицы productions: %s", err)
	}

	return products
}

func (*Production) UpdateAllStatusDel() error {
	_, err := DB.Exec("UPDATE productions SET status = ?", DELETE)
	if err != nil {
		return fmt.Errorf("ошибка обновления строки в productions %s", err)
	}
	return nil
}
