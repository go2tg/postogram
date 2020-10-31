# Почтограм / Postogram / по 100 грамм
Отправка почтовых сообщений через  Telegram

![Go](https://github.com/go2tg/postogram/workflows/Go/badge.svg)

## Общее описание идеи
Сервис которое с одной стороны выступает в виде SMTP, понимая команды smtp и все отправленные через него письма пересылает в виде сообщений в предварительно
настроенный телеграм канал. Основное назначение - как один из микросервисов для Docker стека который позволит не заботится о настройке
оповещений из приложений. Достаточно просто сформировать и отправить письмо.

## Модульная схема
![Diagram](https://github.com/go2tg/postogram/blob/main/postogram.png)


### отправка текстовых сообщений
### отправка вложений

## Предварительная настройка
Для работы через Telegram API необходиме предварительно получить авторизационный токен и ID чатов
1. Create a new Telegram bot: https://core.telegram.org/bots#creating-a-new-bot.
2. Open that bot account in the Telegram account which should receive the messages, press /start.
3. Retrieve a chat id with curl https://api.telegram.org/bot<BOT_TOKEN>/getUpdates.
4. Repeat steps 2 and 3 for each Telegram account which should receive the messages.



## План
- [x] создание концепта и описания
- [ ] отправка тестовых сообщений
- [ ] отправка файлов в виде вложений
- [ ] шаблоны для отправки сообщения
- [ ] возможность настройки нескольких выходных модулей получателей для одного сообщения
- [ ] отправка одного сообщения в несколько каналов получателей
- [ ] модули для slack, mattermost, skype, rocketchat, smtp relay
- [ ] Web UI для управления фильтрами, шаблонами и вых. плагинами
- [ ] механизм хранения ключей и токенов
- [ ] CI/CD
- [ ] Dockerfile
- [ ] Docker service (docker-compose.yml)
- [ ] телеграм бот
