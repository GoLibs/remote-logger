package main

import (
	"fmt"

	remote_logger "github.com/golibs/remote-logger"
)

func main() {
	l := remote_logger.NewLogger()
	err := l.AppendTextToDefaultLogFile("This is a test")
	fmt.Println(err)
	groupLog := l.NewGroupLog()
	err = groupLog.AddText("This is first line").AddText("This is second line").AddText("And here ends").AppendTextToDefaultLogFile()
	fmt.Println(err)
}
