package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send shrug in chat on /shrug
func Shrug(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	return context.Send("¯\\_(ツ)_/¯")
}
