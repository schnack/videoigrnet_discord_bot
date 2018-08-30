package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(pathDB string) (*sql.DB, error) {

	// Устанавливаем подключение к БД
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к %s %s", pathDB, err)
	}

	// Проверяем рабочие таблицы в БД
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name IN('category', 'productions', 'channels', 'channels_categories')")
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении списка таблиц %s", err)
	}
	defer rows.Close()

	currentTables := 0
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("ошибка при получении названий таблиц %s", err)
		}
		currentTables++
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении названий таблиц %s", err)
	}

	// Инициализируем базу если ее не существовало
	if currentTables != 4 {
		err := createTables(db)
		if err != nil {
			return nil, fmt.Errorf("ошибка при создании рабочих таблиц %s", err)
		}
	}
	return db, nil
}

// Создание рабочий таблиц
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE category
(
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    parent_id INTEGER NOT NULL,
    parent_name TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
CREATE UNIQUE INDEX category_id_uindex ON category (id);
CREATE INDEX category_parent_id_index ON category (parent_id);

CREATE TABLE productions
(
    id INTEGER PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    category_id INTEGER NOT NULL,
	buy_status INTEGER NOT NULL,
	status INTEGER NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CONSTRAINT productions_category_id_fk FOREIGN KEY (category_id) REFERENCES category (id)
);
CREATE UNIQUE INDEX productions_id_uindex ON productions (id);
CREATE INDEX productions_status_index ON productions (status);
CREATE INDEX productions_category_id_index ON productions (category_id);

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

CREATE TABLE channels_categories
(
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    category_id INTEGER NOT NULL,
    channel_id INTEGER NOT NULL,
    username TEXT,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    CONSTRAINT channels_categories_category_id_fk FOREIGN KEY (category_id) REFERENCES category (id)
	CONSTRAINT channels_categories_channel_id_fk FOREIGN KEY (channel_id) REFERENCES channels (id)
);
CREATE UNIQUE INDEX channels_categories_id_uindex ON channels_categories (id);
CREATE INDEX channels_categories_category_id_index ON channels_categories (category_id);
CREATE INDEX channels_categories_channel_id_index ON channels_categories (channel_id);
`)
	if err != nil {
		return err
	}
	return nil
}
