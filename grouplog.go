package remote_logger

import (
	"strings"
	"sync"
)

type GroupLog struct {
	p     *RemoteLogger
	m     sync.Mutex
	texts []string
	files []*file
}

func newGroupLog(p *RemoteLogger) *GroupLog {
	return &GroupLog{p: p}
}

func (gp *GroupLog) AddFileBytes(fileName, caption string, fileBytes []byte) *GroupLog {
	defer gp.m.Unlock()
	gp.m.Lock()
	gp.files = append(gp.files, newFileBytes(fileName, caption, fileBytes))
	return gp
}

func (gp *GroupLog) reset() {
	defer gp.m.Unlock()
	gp.m.Lock()
	gp.texts = nil
	gp.files = nil
}

func (gp *GroupLog) AddText(text string) *GroupLog {
	defer gp.m.Unlock()
	gp.m.Lock()
	gp.texts = append(gp.texts, text)
	return gp
}

func (gp *GroupLog) AppendTextToDefaultLogFile() error {
	defer gp.m.Unlock()
	gp.m.Lock()
	if len(gp.texts) > 0 {
		return gp.p.AppendTextToDefaultLogFile(strings.Join(gp.texts, "\r\n"))
	}
	return nil
}

func (gp *GroupLog) SendToDefaultTelegramChat() error {
	defer gp.m.Unlock()
	gp.m.Lock()
	if len(gp.texts) > 0 {
		err := gp.p.SendTextToDefaultTelegramChat(strings.Join(gp.texts, "\r\n"))
		if err != nil {
			return err
		}
	}
	if len(gp.files) > 0 {
		for _, f := range gp.files {
			err := gp.p.SendFileBytesToDefaultTelegramChat(f.name, f.caption, f.bytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (gp *GroupLog) SendToDefaultDiscordWebhook() error {
	defer gp.m.Unlock()
	gp.m.Lock()
	if len(gp.texts) > 0 {
		err := gp.p.SendTextToDefaultDiscordWebhook(strings.Join(gp.texts, "\r\n"))
		if err != nil {
			return err
		}
	}
	if len(gp.files) > 0 {
		for _, f := range gp.files {
			err := gp.p.SendFileBytesToDefaultDiscordWebhook(f.name, f.bytes)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
