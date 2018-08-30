# Бот оповещает о новинках на https://videoigr.net

Установка Docker

    docker build -t vgnet https://github.com/schnack/videoigrnet_discord_bot.git

Запуск 

    docker run -d vgnet -t <token>

Добавление категории:

    [vgnet add https://videoigr.net/index.php?cPath=142_146
			
Просмотр списка категорий:
    
    [vgnet list
			
Удаление категории:
    
    [vgnet del <num>

Запуск уведомлений в текущем канале:
    
    [vgnet start
			
Остановка уведомлений в текущем канале:
    
    [vgnet stop

Посмотреть статус уведомлений в текущем канале:
    
    [vgnet status

