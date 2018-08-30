Бот оповещает о изменениях на сайте videoigr.net.

Установка Docker

docker build -t vgnet https://github.com/schnack/videoigrnet_discord_bot.git

Запуск 

docker run -d vgnet -t <token>

Добавление категории:

    [vgnet add <url>
			
Просмотр списка категорий:
    
    [vgnet list
			
Удаление категории:
    
    [vgnet del <num>

Запуск уведомлений:
    
    [vgnet start
			
Остановка уведомлений
    
    [vgnet stop

Посмотреть статус уведомлений в текущем канале
    
    [vgnet status

