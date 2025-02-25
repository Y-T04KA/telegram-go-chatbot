package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

//Unmute user on /unmute
func Unmute(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/unmute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unmute</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember.CanSendMessages = true
	TargetChatMember.CanSendMedia = true
	TargetChatMember.CanSendPolls = true
	TargetChatMember.CanSendOther = true
	TargetChatMember.CanAddPreviews = true
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Restrict(context.Chat(), TargetChatMember)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка снятия ограничения пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> снова может отправлять сообщения в чат.", target.ID, utils.UserFullName(&target)))
}
