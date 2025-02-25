package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

//Kill user on /kill
func Kill(context tele.Context) error {
	command := "/bless"
	if context.Text()[1:5] == "kill" {
		command = "/kill"
	}
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply(fmt.Sprintf("Пример использования: <code>%v {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>%v</code>", command, command))
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	if context.Message().ReplyTo != nil {
		utils.Bot.Delete(context.Message().ReplyTo)
	}
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
		return context.Send(fmt.Sprintf("<code>👻 %v возродился у костра.</code>", utils.UserFullName(&target)))
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(target.ID).First(&duelist)
	if result.RowsAffected == 0 {
		duelist.UserID = target.ID
		duelist.Kills = 0
		duelist.Deaths = 0
	}
	duelist.Deaths++
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&duelist)
	if result.Error != nil {
		return err
	}
	duration := utils.RandInt(1, duelist.Deaths+1)
	prependText := ""
	if utils.RandInt(0, 100) >= 98 {
		duration = duration * 10
		prependText = "критически "
		if command == "/bless" {
			prependText = "очень "
		}
	}
	if duration > 600 {
		duration = 600
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(60*duration)).Unix()
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("💥 %v %vпристрелил %v.\n%v отправился на респавн на %v мин.", utils.UserFullName(context.Sender()), prependText, utils.UserFullName(&target), utils.UserFullName(&target), duration)
	if command == "/bless" {
		text = fmt.Sprintf("🤫 %v %vпопросил %v помолчать %v минут.", utils.UserFullName(context.Sender()), prependText, utils.UserFullName(&target), duration)
	}
	return context.Send(text)
}
