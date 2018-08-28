package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(pathDB string) error {

	var err error
	// Устанавливаем подключение к БД
	DB, err = sql.Open("sqlite3", pathDB)
	if err != nil {
		return err
	}

	// Проверяем рабочие таблицы в БД
	rows, err := DB.Query("SELECT name FROM sqlite_master WHERE type='table' AND name IN('products', 'channels', 'channels_products', 'settings')")
	if err != nil {
		return err
	}
	defer rows.Close()

	currentTables := 0
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		currentTables++
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// Инициализируем базу если ее не существовало
	if currentTables != 4 {
		err := createTables()
		if err != nil {
			return err
		}
		err = initSettings()
		if err != nil {
			return err
		}
	}
	return nil
}

// Создание рабочий таблиц
func createTables() error {
	_, err := DB.Exec(`
CREATE TABLE products
(
    id INTEGER PRIMARY KEY,
    name TEXT,
    category_id INTEGER,
    category_name TEXT,
    category_parent_id INTEGER,
    category_parent_name TEXT,
    buy_status INTEGER,
    status INTEGER,
    created_at INTEGER,
    updated_at INTEGER
);
CREATE INDEX products_status_index ON products (status);
CREATE INDEX products_category_id_index ON products (category_id);
CREATE INDEX products_buy_status_index ON products (buy_status);
CREATE INDEX products_category_parent_id_index ON products (category_parent_id);

CREATE TABLE channels
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    channel TEXT,
	status INTEGER,
    created_at INTEGER,
    updated_at INTEGER
);
CREATE UNIQUE INDEX channels_id_uindex ON channels (id);
CREATE INDEX channels_channel_index ON channels (channel);

CREATE TABLE channels_products
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    product_id INTEGER,
    channel_id INTEGER,
    username TEXT,
    created_at INTEGER,
    updated_at INTEGER,
    CONSTRAINT channels_products_products_id_fk FOREIGN KEY (product_id) REFERENCES products (id),
    CONSTRAINT channels_products_channels_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id)
);
CREATE UNIQUE INDEX channels_products_id_uindex ON channels_products (id);
CREATE INDEX channels_products_product_id_index ON channels_products (product_id);
CREATE INDEX channels_products_channel_id_index ON channels_products (channel_id);

CREATE TABLE settings
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    name TEXT NOT NULL,
    value TEXT NOT NULL,
	created_at INTEGER,
    updated_at INTEGER
);
CREATE UNIQUE INDEX settings_id_uindex ON settings (id);
CREATE UNIQUE INDEX settings_name_uindex ON settings (name);
`)
	if err != nil {
		return fmt.Errorf("Первичная инициализация БД: %s", err)
	}
	return nil
}

func initSettings() error {
	setting := Setting{Name: "sync_timeout", Value: "60"}
	return setting.Save()
}
