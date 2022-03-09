package handlers

import (
	tele "gopkg.in/telebot.v3"
	"log"
	"nodemon/pkg/entities"
	"nodemon/pkg/messaging"
)

func InitHandlers(bot *tele.Bot, environment *messaging.TelegramBotEnvironment) {
	bot.Handle("/hello", func(c tele.Context) error {
		oldChatID, err := environment.ChatStorage.FindChatID(entities.TelegramPlatform)
		if err != nil {
			log.Printf("failed to insert chat id into db: %v", err)
			return c.Send("An error occurred while finding the chat id in database")
		}
		if oldChatID != nil {
			return c.Send("Hello! I remember this chat.")
		}
		chatID := entities.ChatID(c.Chat().ID)

		err = environment.ChatStorage.InsertChatID(chatID, entities.TelegramPlatform)
		if err != nil {
			log.Printf("failed to insert chat id into db: %v", err)
			return c.Send("I failed to save this chat id")
		}
		return c.Send("Hello! This new chat has been saved for alerting.")
	})

	bot.Handle("/ping", func(c tele.Context) error {
		return c.Send("pong!")
	})

	bot.Handle("/start", func(c tele.Context) error {
		environment.ShutUp = true
		return c.Send("Started working...")
	})

	bot.Handle("/mute", func(c tele.Context) error {
		environment.ShutUp = false
		return c.Send("Say no more..")
	})
}
