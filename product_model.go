package main

import (
	"fmt"
	"time"
)

const (
	NEW int = iota
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
	return p.Id == np.Id && p.Name == np.Name && p.Category.Equal(np.Category) && p.BuyStatus == np.BuyStatus && p.Status == np.Status
}

// Save сохраняет состояние структуры в базу данных
func (p *Production) Save() error {
	var err error
	tmp := (&Production{}).FindById(p.Id)
	if tmp != nil {
		if p.Equal(tmp) {
			return nil
		} else {
			p.UpdatedAt = time.Now().Unix()
			_, err = DB.Exec("UPDATE productions SET name = ?, category_id = ?, buy_status = ?, status = ?, updated_at = ? WHERE id = ?", p.Name, p.Category.Id, p.BuyStatus, p.Status, p.UpdatedAt, p.Id)
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
func (*Production) FindById(v int64) *Production {
	p := &Production{}
	var cId int64
	err := DB.QueryRow("SELECT id, name, category_id, buy_status, status, created_at, updated_at FROM productions WHERE id = ?", v).Scan(&p.Id, &p.Name, &cId, &p.BuyStatus, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil
	}
	p.Category = (&Category{}).FindById(cId)
	return p
}

/*

type ProductImport struct {
	Id                 string `json:"id"`
	Name               string `json:"name"`
	CategoryId         string `json:"cat_id"`
	CategoryName       string `json:"cat_name"`
	CategoryParentId   string `json:"cat_parent_id"`
	CategoryParentName string `json:"cat_parent_name"`
	BuyStatus          string `json:"buy_status"`
}

func (pi *ProductImport) Conv() *Product {
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
	return &Product{Id: int(Id), Name: pi.Name, CategoryId: int(CategoryId), CategoryName: pi.CategoryName, CategoryParentId: int(CategoryParentId), CategoryParentName: pi.CategoryParentName, BuyStatus: int(BuyStatus)}
}


func (p *Product) Save() error {
	var id int
	err := DB.QueryRow("SELECT id FROM products WHERE id = ?", p.Id).Scan(&id)
	if err != nil {
		p.UpdatedAt = int(time.Now().Unix())
		p.CreatedAt = p.UpdatedAt
		p.Status = NEW
		res, err := DB.Exec("INSERT INTO products(id, name, category_id, category_name, category_parent_id, category_parent_name, buy_status, status, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", p.Id, p.Name, p.CategoryId, p.CategoryName, p.CategoryParentId, p.CategoryParentName, p.BuyStatus, p.Status, p.CreatedAt, p.UpdatedAt)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("Попытка получить id добавленой записи %s", err)
		}
		p.Id = int(id)
	} else {
		p.UpdatedAt = int(time.Now().Unix())
		p.Status = EXIST
		_, err = DB.Exec("UPDATE products SET name = ?, category_id = ?, category_name = ?, category_parent_id = ?, category_parent_name = ?, buy_status = ?, status = ?, updated_at = ? WHERE id = ?", p.Name, p.CategoryId, p.CategoryName, p.CategoryParentId, p.CategoryParentName, p.BuyStatus, p.Status, p.UpdatedAt, p.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Product) FindCategory(categoryId int) *Product {
	err := DB.QueryRow("SELECT id, name, category_id, category_name, category_parent_id, category_parent_name, buy_status, status, created_at, updated_at FROM products WHERE category_id = ?", categoryId).Scan(&p.Id, &p.Name, &p.CategoryId, &p.CategoryName, &p.CategoryParentId, &p.CategoryParentName, &p.BuyStatus, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil
	}
	return p
}

func (p *Product) FindId(productId int) *Product {
	err := DB.QueryRow("SELECT id, name, category_id, category_name, category_parent_id, category_parent_name, buy_status, status, created_at, updated_at FROM products WHERE id = ?", productId).Scan(&p.Id, &p.Name, &p.CategoryId, &p.CategoryName, &p.CategoryParentId, &p.CategoryParentName, &p.BuyStatus, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil
	}
	return p
}

func MarkProductDelete() {
	_, err := DB.Exec("UPDATE products SET status = ?", DELETE)
	if err != nil {
		log.Fatal(err)
	}
}

func (*Product) FindStatus(status int) []*Product {
	products := make([]*Product, 0)
	rows, err := DB.Query("SELECT id, name, category_id, category_name, category_parent_id, category_parent_name, buy_status, status, created_at, updated_at FROM products WHERE status = ? ", status)
	if err != nil {
		log.Println("Не удалось найти продукты со статусом ", status)
		return products
	}
	defer rows.Close()
	for rows.Next() {
		p_tmp := Product{}
		err := rows.Scan(&p_tmp.Id, &p_tmp.Name, &p_tmp.CategoryId, &p_tmp.CategoryName, &p_tmp.CategoryParentId, &p_tmp.CategoryParentName, &p_tmp.BuyStatus, &p_tmp.Status, &p_tmp.CreatedAt, &p_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		products = append(products, &p_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return products
}

func (*Product) FindStatusCategory(status int, category int) []*Product {
	products := make([]*Product, 0)
	rows, err := DB.Query("SELECT id, name, category_id, category_name, category_parent_id, category_parent_name, buy_status, status, created_at, updated_at FROM products WHERE status != ? and category_id = ?", status, category)
	if err != nil {
		log.Println("Не удалось найти продукты со статусом ", status)
		return products
	}
	defer rows.Close()
	for rows.Next() {
		p_tmp := Product{}
		err := rows.Scan(&p_tmp.Id, &p_tmp.Name, &p_tmp.CategoryId, &p_tmp.CategoryName, &p_tmp.CategoryParentId, &p_tmp.CategoryParentName, &p_tmp.BuyStatus, &p_tmp.Status, &p_tmp.CreatedAt, &p_tmp.UpdatedAt)
		if err != nil {
			log.Println("Не удалось восстановить объект")
		}
		products = append(products, &p_tmp)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return products
}
*/
