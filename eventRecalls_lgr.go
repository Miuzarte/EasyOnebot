package EasyOnebot

import (
	"runtime/debug"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/internal/utils"
)

type eventRecalls_Lgr struct {
	noticeBotOffline  []func(*event.NoticeBotOffline)
	noticeBotOnline   []func(*event.NoticeBotOnline)
	noticeEssence     []func(*event.NoticeEssence)
	noticeGroupName   []func(*event.NoticeGroupName)
	noticeReaction    []func(*event.NoticeReaction)
	noticeOfflineFile []func(*event.NoticeOfflineFile)
}

func (b *Bot) OnBotOffline(handler func(*event.NoticeBotOffline)) *Bot {
	b.eventRecalls.noticeBotOffline = utils.OnDemandAppend(b.eventRecalls.noticeBotOffline, handler)
	return b
}

func (b *Bot) OnBotOnline(handler func(*event.NoticeBotOnline)) *Bot {
	b.eventRecalls.noticeBotOnline = utils.OnDemandAppend(b.eventRecalls.noticeBotOnline, handler)
	return b
}

func (b *Bot) OnEssence(handler func(*event.NoticeEssence)) *Bot {
	b.eventRecalls.noticeEssence = utils.OnDemandAppend(b.eventRecalls.noticeEssence, handler)
	return b
}

func (b *Bot) OnGroupName(handler func(*event.NoticeGroupName)) *Bot {
	b.eventRecalls.noticeGroupName = utils.OnDemandAppend(b.eventRecalls.noticeGroupName, handler)
	return b
}

func (b *Bot) OnReaction(handler func(*event.NoticeReaction)) *Bot {
	b.eventRecalls.noticeReaction = utils.OnDemandAppend(b.eventRecalls.noticeReaction, handler)
	return b
}

func (b *Bot) OnOfflineFile(handler func(*event.NoticeOfflineFile)) *Bot {
	b.eventRecalls.noticeOfflineFile = utils.OnDemandAppend(b.eventRecalls.noticeOfflineFile, handler)
	return b
}

func (b *Bot) doEventRecall_Lgr(typedEvent any) {
	i := new(int)
	defer func() {
		if err := recover(); err != nil {
			b.log.Error().Msgf("panic while handling %T[%d]: %v", typedEvent, *i, err)
			debug.PrintStack()
		}
	}()
	switch e := typedEvent.(type) {
	case *event.NoticeBotOffline:
		for j, f := range b.eventRecalls.noticeBotOffline {
			*i = j
			f(e)
		}
	case *event.NoticeBotOnline:
		for j, f := range b.eventRecalls.noticeBotOnline {
			*i = j
			f(e)
		}
	case *event.NoticeEssence:
		for j, f := range b.eventRecalls.noticeEssence {
			*i = j
			f(e)
		}
	case *event.NoticeGroupName:
		for j, f := range b.eventRecalls.noticeGroupName {
			*i = j
			f(e)
		}
	case *event.NoticeReaction:
		for j, f := range b.eventRecalls.noticeReaction {
			*i = j
			f(e)
		}
	case *event.NoticeOfflineFile:
		for j, f := range b.eventRecalls.noticeOfflineFile {
			*i = j
			f(e)
		}
	}
}
