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

// Variables used for command line parameters
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
	// Инициализируем подключение к discord
	DG, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DG.Close()

	// Инициализируем подключение к БД
	err = initDB("./db/vgnet.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	done := make(chan struct{}, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scanVideoigrNet(done)
	}()

	// Устанавливаем обработчик событий в каналах.
	DG.AddHandler(router)

	// Открываем соединение
	err = DG.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Завершаем работу по CTRL-C с корректным заверешением всех подпрограмм.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	close(done)

	wg.Wait()
	log.Println("Работа завершена")
}
