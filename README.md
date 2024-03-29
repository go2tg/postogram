# Почтограм / Postogram / по 100 грамм
Отправка почтовых сообщений через  Telegram

![Go](https://github.com/go2tg/postogram/workflows/Go/badge.svg)

## Общее описание идеи
Сервис которое с одной стороны выступает в виде SMTP, понимая команды smtp и все отправленные через него письма пересылает в виде сообщений в предварительно
настроенный телеграм канал. Основное назначение - как один из микросервисов для Docker стека который позволит не заботится о настройке
оповещений из приложений. Достаточно просто сформировать и отправить письмо.

## Модульная схема
![Diagram](https://github.com/go2tg/postogram/blob/main/postogram.png)


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
- [ ] модули для slack, mattermost, skype, rocketchat, smtp relay, mqtt brocker
- [ ] Web UI для управления фильтрами, шаблонами и вых. плагинами
- [ ] механизм хранения ключей и токенов
- [ ] CI/CD
- [ ] Dockerfile
- [ ] Docker service (docker-compose.yml)
- [ ] телеграм бот


## API reference
  - telegram - https://core.tlgr.org/bots/api
  - skype - https://docs.microsoft.com/en-us/skype-sdk/skypeuris/skypeuriapireference
  - slack - https://api.slack.com/
  - MQTT - https://thingsboard.io/docs/reference/mqtt-api/
  - smtp commands - https://blog.mailtrap.io/smtp-commands-and-responses/ , https://www.ibm.com/support/knowledgecenter/en/SSLTBW_2.3.0/com.ibm.zos.v2r3.halu001/smtpcommands.htm


## FAQ

### Почему отличаются русское и английское названия ? 
Из-за игры слов русское название выглядит несколько неоднозначно. Поэтому принято решение немного его изменить так, чтобы лучше отражало цель проекта. Но неофициально внутри команды мы по прежнему используем неизмененное название. Просто так веселей. :) 

## Contributing

If you'd like to contribute new features, enhancements or bug fixes to the code base just follow these steps:

* Create a [GitHub](https://github.com/signup/free) account, if you don't own one already
* Then, [fork](https://help.github.com/articles/fork-a-repo) the [postogram](https://github.com/go2tg/postogram) repository to your account
* Create a new [branch](https://help.github.com/articles/creating-and-deleting-branches-within-your-repository) from the *develop* branch in your forked repository
* Modify the existing code, or add new code to your branch
* When ready, make a [pull request](http://help.github.com/send-pull-requests/) to the main repository

There may be some discussion regarding your contribution to the repository before any code is merged in, so be prepared to provide feedback on your contribution if required.
