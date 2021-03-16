package main

import (
	"fmt"

	remote_logger "github.com/golibs/remote-logger"
)

func main() {
	l := remote_logger.NewLogger()
	err := l.SetTelegramBotToken("REPLACE HERE")
	if err != nil {
		fmt.Println(err)
		return
	}
	chatId := int64(0) // REPLACE HERE
	l.SetDefaultTelegramChatID(chatId)
	err = l.SendTextToDefaultTelegramChat("This is a test")
	if err != nil {
		fmt.Println(err)
		return
	}
	groupLog := l.NewGroupLog()
	err = groupLog.AddText("This is first line").AddText("This is second line").AddText("And here ends").SendToDefaultTelegramChat()
	fmt.Println(err)
}
