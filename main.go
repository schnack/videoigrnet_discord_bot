package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	Token string
	DB    *sql.DB
	DG    *discordgo.Session
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	var err error

	// Инициализируем БД
	DB, err = initDB("./db/vgnet.sqlite3")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DB.Close()

	// Инициализируем Discord
	DG, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DG.Close()

	// Служит сигналом для завершения go программ
	done := make(chan struct{}, 1)
	// Счетчик ожидания go подпрограмм
	var wg sync.WaitGroup

	// Передаем события Discord маршрутизатору.
	DG.AddHandler(router)

	// Старт бота
	err = DG.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Запуск программы слежения за https://videoigr.net
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanVideoigrNet(done, DG)
	}()

	// Завершаем работу по CTRL-C с корректным заверешением всех подпрограмм.
	fmt.Println("Бот запущен! Нажмите CTRL-C для выхода.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	// Поссылаем сигнал на завершение go программ
	close(done)

	wg.Wait()
	log.Println("Работа завершена")
}
