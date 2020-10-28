# Почтограм
Отправка почтовых сообщений через  Telegram

Приложение которое биндится в виде фейкового SMTP сервера и все отправленные через него письма отсылает в предварительно 
настроенный телеграм канал. Основное назначение - как один из микросервисов для Docker стека который позволит не заботится о настройке 
оповещений из приложений. Достаточно просто сформировать и отправить письмо. 


## План
- [x] создание концепта и описания
- [ ] механизм хранения ключей и токенов
- [ ] CI/CD
- [ ] отправка тестовых сообщений
- [ ] отправка файлов в виде вложений
- [ ] Dockerfile
- [ ] Docker service (docker-compose.yml)
- [ ] Docker swarm stack (docker-stack.yml)
- [ ] отправка в несколько каналов
- [ ] шаблоны для отправки сообщения
- [ ] телеграм бот
