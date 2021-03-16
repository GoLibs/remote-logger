package remote_logger

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	go_telegram_bot_api "github.com/GoLibs/telegram-bot-api"
)

type RemoteLogger struct {
	fileMutex                sync.Mutex
	filesMutex               map[string]*sync.Mutex
	telegramBotToken         string
	telegramBotClient        *go_telegram_bot_api.TelegramBot
	defaultDiscordWebhookUrl string
	tz                       *time.Location
	defaultLogDirectory      string
	defaultLogFileName       string
	defaultTelegramChatId    int64
	r                        *resty.Client
}

func NewLogger() *RemoteLogger {
	return &RemoteLogger{defaultLogDirectory: "logs", defaultLogFileName: "logs.log", filesMutex: map[string]*sync.Mutex{}, r: resty.New()}
}

func (rl *RemoteLogger) formatText(text string) string {
	return fmt.Sprintf("%s - %s\n", rl.time().Format("2006-01-02 15:04:05"), text)
}

func (rl *RemoteLogger) mutexForPath(path string) *sync.Mutex {
	val, ok := rl.filesMutex[path]
	if !ok {
		val = &sync.Mutex{}
		rl.filesMutex[path] = val
	}
	return val
}

func (rl *RemoteLogger) time() (t time.Time) {
	l := time.Local
	if rl.tz != nil {
		l = rl.tz
	}
	t = time.Now().In(l)
	return
}

func (rl *RemoteLogger) defaultLogFilePath() string {
	return rl.defaultLogDirectory + string(os.PathSeparator) + rl.defaultLogFileName
}

func (rl *RemoteLogger) NewGroupLog() *GroupLog {
	return newGroupLog(rl)
}

func (rl *RemoteLogger) SetDefaultDirectory(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return
		}
	}
	rl.defaultLogDirectory = dir
	return
}

func (rl *RemoteLogger) SetDefaultLogFileName(name string) {
	rl.defaultLogFileName = name
}

func (rl *RemoteLogger) SetTimezone(loc string) (err error) {
	rl.tz, err = time.LoadLocation(loc)
	return
}

func (rl *RemoteLogger) SetDiscordWebhookUrl(address string) {
	rl.defaultDiscordWebhookUrl = address
}

func (rl *RemoteLogger) SetTelegramBotToken(t string) (err error) {
	rl.telegramBotToken = t
	rl.telegramBotClient, err = go_telegram_bot_api.NewTelegramBot(t)
	if err != nil {
		return
	}
	return
}

func (rl *RemoteLogger) SetDefaultTelegramChatID(chatId int64) (err error) {
	rl.defaultTelegramChatId = chatId
	return
}

func (rl *RemoteLogger) SendTextToDefaultDiscordWebhook(text string) (err error) {
	if rl.defaultDiscordWebhookUrl == "" {
		err = errors.New("discord_webhook_not_set")
		return
	}
	r := rl.r.R()
	_, err = r.SetHeader("Content-Type", "application/json").SetBody(map[string]string{
		"content": text,
	}).Post(rl.defaultDiscordWebhookUrl)
	return
}

func (rl *RemoteLogger) SendTextToDefaultTelegramChat(text string) (err error) {
	if rl.telegramBotClient == nil {
		err = errors.New("telegram_bot_client_not_initialised")
		return
	}
	text = rl.formatText(text)
	_, err = rl.telegramBotClient.Send(rl.telegramBotClient.Message().SetChatId(rl.defaultTelegramChatId).SetText(text))
	return
}

func (rl *RemoteLogger) SendFileBytesToDefaultDiscordWebhook(fileName string, fileBytes []byte) (err error) {
	if rl.defaultDiscordWebhookUrl == "" {
		err = errors.New("discord_webhook_not_set")
		return
	}
	r := rl.r.R()
	_, err = r.SetHeader("Content-Type", "multipart/form-data").SetFileReader("file1", fileName, bytes.NewReader(fileBytes)).Post(rl.defaultDiscordWebhookUrl)
	return
}

func (rl *RemoteLogger) SendFileBytesToDefaultTelegramChat(fileName, caption string, fileBytes []byte) (err error) {
	if rl.telegramBotClient == nil {
		err = errors.New("telegram_bot_client_not_initialised")
		return
	}
	_, err = rl.telegramBotClient.Send(rl.telegramBotClient.Document().SetChatId(rl.defaultTelegramChatId).SetCaption(caption).SetDocumentReader(bytes.NewReader(fileBytes), fileName))
	return
}

func (rl *RemoteLogger) AppendTextToDefaultLogFile(text string) (err error) {
	if _, err = os.Stat(rl.defaultLogDirectory); os.IsNotExist(err) {
		err = os.Mkdir(rl.defaultLogDirectory, os.ModePerm)
		if err != nil {
			return
		}
	}
	text = rl.formatText(text)
	path := rl.defaultLogFilePath()
	rl.fileMutex.Lock()
	defer rl.mutexForPath(path).Unlock()
	rl.mutexForPath(path).Lock()
	rl.fileMutex.Unlock()
	var f *os.File
	f, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	f.Write([]byte(text))
	f.Close()
	return
}
