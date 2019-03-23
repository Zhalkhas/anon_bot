# ANON_BOT
Simple Telegram bot written on Go, which allows to resend messages to desired chat anonymously. Bot uses webhooks to get updates and I suggest to host it on Heroku. In order to use this bot, constants should be changed to your values 
``` 
const (       
    token = "<Here goes token of telegram bot>"
    id    = -0000000000 // Here goes ID of chat
    url   = "<Here goes URL of server>"
    )
```
Bot uses [gb](https://github.com/constabulary/gb) for resolving dependencies and requires [gin](https://github.com/gin-gonic/gin) and [tgbotapi](https://github.com/go-telegram-bot-api/telegram-bot-api) installed. 
