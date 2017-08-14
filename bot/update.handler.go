package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ishanjain28/chbot/db"
)

func handleUpdates(bot *tbot.BotAPI, db *db.DB, u tbot.Update) {

	if u.Message != nil && u.Message.IsCommand() {
		handleCommands(bot, db, u)
	}

	if u.Message != nil && u.Message.Text != "" {
		handleTextInput(bot, u)
	}
}

func handleCommands(bot *tbot.BotAPI, d *db.DB, u tbot.Update) {

	chatID := u.Message.Chat.ID

	switch u.Message.Text {
	case "/start":

		fmt.Println(u.Message.Date, time.Now().Unix())
		data, err := d.FetchSubscribers()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(data)
	case "/subscribe":
		user := &db.User{
			ChatID:   chatID,
			Username: u.Message.From.UserName,
		}
		err := d.AddSubscriber(user)
		if err != nil {
			if mgo.IsDup(err) {
				msg := tbot.NewMessage(chatID, "You are already subscribed, Type /unsubscribe to unsubscribe")
				bot.Send(msg)
				return
			}
			sendErrorMessage(bot, chatID)
			log.Printf("error occurred in adding subscriber: %v", err)
			return
		}

		msg := tbot.NewMessage(chatID, "Subscription successful! I'll now send you the latest C&H Comics on a daily basis")
		bot.Send(msg)

	case "/unsubscribe":
		user := &db.User{
			ChatID:   chatID,
			Username: u.Message.From.UserName,
		}
		err := d.RemoveSubscriber(user)
		if err != nil {
			if err == mgo.ErrNotFound {
				msg := tbot.NewMessage(chatID, "You are not on the subscribers List, Type /subscribe to subscribe")
				bot.Send(msg)
				return
			}
			sendErrorMessage(bot, chatID)
			log.Printf("error occurred in unsubscribing: %v", err)
			return
		}

		msg := tbot.NewMessage(chatID, "Unsubscription successful! You'll now receive no updates from me\n If you wish to subscribe again, Use the /subscribe command")
		bot.Send(msg)

	case "/help":
	default:
		msg := tbot.NewMessage(chatID, "Invalid command, Type /help to get help")
		bot.Send(msg)
	}
}

func handleTextInput(bot *tbot.BotAPI, u tbot.Update) {
	switch strings.ToLower(u.Message.Text) {
	case "random":

	default:

	}
}

func sendErrorMessage(bot *tbot.BotAPI, chatID int64) {
	msg := tbot.NewMessage(chatID, "Error Occurred, Please retry")
	bot.Send(msg)
}
