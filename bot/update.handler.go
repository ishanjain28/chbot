package bot

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/ishanjain28/chbot/ch"

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

		randomBtn := tbot.NewKeyboardButton("Random")

		k := tbot.NewKeyboardButtonRow(randomBtn)

		keyboard := tbot.NewReplyKeyboard(k)

		msg := tbot.NewMessage(chatID, "This bot can send you the latest  Cyanide and Happiness Comics on a Daily Basic or you can read random C&H comics")
		msg.ReplyMarkup = keyboard
		bot.Send(msg)

		// data, err := d.FetchSubscribers()
		// if err != nil {
		// 	log.Fatalln(err)
		// }
		// fmt.Println(data)
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

	chatID := u.Message.Chat.ID
	switch strings.ToLower(u.Message.Text) {
	case "random":
		// previous := "5"
		// next := "10"

		// rows := tbot.NewInlineKeyboardRow(
		// 	tbot.InlineKeyboardButton{CallbackData: &previous, Text: "<<"},
		// 	tbot.InlineKeyboardButton{CallbackData: &next, Text: ">>"})

		// keyboard := tbot.NewInlineKeyboardMarkup(rows)

		fileLink, err := ch.Random()
		if err != nil {
			sendErrorMessage(bot, chatID)
			return
		}

		// fileURL, err := url.Parse(fileLink)
		// if err != nil {
		// 	sendErrorMessage(bot, chatID)
		// 	return
		// }

		filename, resp, err := ch.Download(fileLink)
		defer resp.Body.Close()
		if err != nil {
			sendErrorMessage(bot, chatID)
			return
		}

		b, _ := ioutil.ReadAll(resp.Body)

		msg := tbot.NewPhotoUpload(chatID, tbot.FileBytes{Bytes: b, Name: filename})
		// msg.ReplyMarkup = keyboard
		bot.Send(msg)

	default:
	}
}

func sendErrorMessage(bot *tbot.BotAPI, chatID int64) {
	msg := tbot.NewMessage(chatID, "Error Occurred, Please retry")
	bot.Send(msg)
}
