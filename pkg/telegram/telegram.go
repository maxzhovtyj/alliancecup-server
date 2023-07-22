package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/url"
)

type Manager interface {
	Send(message string) error
	SendRaw(message string) error
	SendWithLink(message, link string) error
}

type BotManager struct {
	bot        *tgbotapi.BotAPI
	tgBotToken string
	tgChatID   int64
}

func NewManager(token string, chatID int64) *BotManager {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	return &BotManager{
		bot:        bot,
		tgBotToken: token,
		tgChatID:   chatID,
	}
}

func (tg *BotManager) Send(message string) error {
	return tg.send(message, nil)
}

func (tg *BotManager) SendRaw(message string) error {
	msg := tgbotapi.NewMessage(tg.tgChatID, message)
	_, err := tg.bot.Send(msg)
	return err
}

func (tg *BotManager) SendWithLink(message, link string) error {
	uri, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("error while parsing link %q: %s", link, err)
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: uri.Host,
				URL:  &link,
			},
		),
	)

	return tg.send(message, &markup)
}

func (tg *BotManager) send(message string, markup *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(tg.tgChatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if markup != nil {
		msg.ReplyMarkup = markup
	}
	_, err := tg.bot.Send(msg)
	return err
}
