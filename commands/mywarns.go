package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

//Send warning amount on /mywarns
func Mywarns(context tele.Context) error {
	var warn utils.Warn
	result := utils.DB.First(&warn, context.Sender().ID)
	if result.RowsAffected != 0 {
		warn.Amount = warn.Amount - int(time.Since(warn.LastWarn).Hours()/24/7)
		if warn.Amount < 0 {
			warn.Amount = 0
		}
	} else {
		warn.UserID = context.Sender().ID
		warn.LastWarn = time.Unix(0, 0)
		warn.Amount = 0
	}
	warnStrings := []string{"предупреждений", "предупреждение", "предупреждения", "предупреждения"}
	return context.Reply(fmt.Sprintf("У тебя %v %v.", warn.Amount, warnStrings[warn.Amount]))
}
