package bot

import (
	"fmt"
	"log"
	"os"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/chbot/db"
)

// TOKEN telegram bot token
var TOKEN = os.Getenv("TOKEN")

// GO_ENV is used to set environment of application, as it will work differently in different environements
var GO_ENV = os.Getenv("GO_ENV")

// Start the bot
func Start(db *db.DB) {

	if TOKEN == "" {
		log.Fatalln("$TOKEN not set")
	}

	bot, err := tbot.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatalln("error in starting bot: %v", err)
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	updates := fetchUpdates(bot)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil && update.InlineQuery == nil && update.EditedMessage == nil {
			continue
		}

		handleUpdates(bot, db, update)
	}
}

func fetchUpdates(bot *tbot.BotAPI) tbot.UpdatesChannel {
	if GO_ENV == "production" {
		// Use webhook if in production

	} else {
		// Use polling if not in production
		bot.RemoveWebhook()

		u := tbot.NewUpdate(0)
		u.Timeout = 60
		updates, err := bot.GetUpdatesChan(u)
		if err != nil {
			log.Printf("warn: error in fetching updates %v", err)

		}

		return updates
	}

	return nil
}
