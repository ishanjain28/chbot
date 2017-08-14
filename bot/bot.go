package bot

import (
	"fmt"
	"log"
	"net/http"
	"os"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/chbot/db"
)

var (
	// TOKEN telegram bot token
	TOKEN = os.Getenv("TOKEN")

	// GO_ENV is used to set environment of application, as it will work differently in different environements
	GO_ENV = os.Getenv("GO_ENV")

	// HOST Address
	HOST = os.Getenv("HOST")
)

// Start the bot
func Start(db *db.DB) {

	if TOKEN == "" {
		log.Fatalln("$TOKEN not set")
	}

	bot, err := tbot.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatalln("error in starting bot: %v", err)
	}

	if GO_ENV != "production" {
		bot.Debug = true
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

		// Remove any existing webhook
		bot.RemoveWebhook()

		// Set a new webhook
		_, err := bot.SetWebhook(tbot.NewWebhook(HOST + "/chbot/" + bot.Token))
		if err != nil {
			log.Fatalf("error in setting webhook: %v", err)
		}

		updates := bot.ListenForWebhook("/chbot/" + bot.Token)

		//redirect users visiting "/" to bot's telegram page
		http.HandleFunc("/", redirectToTelegram)

		return updates
	}

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

func redirectToTelegram(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://t.me/cyanidesub_bot", http.StatusTemporaryRedirect)
}
