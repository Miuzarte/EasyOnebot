package EasyOnebot

import (
	"testing"
	"time"

	"github.com/Miuzarte/EasyOnebot/event"
)

// TestNormalizePoke 固定 NapCat 与 OneBot 11 标准在 poke 通知上的字段差异。
//
// 实测 NapCat 真实推送 (4.18.19, bot 主动戳好友时):
//
//	{"notice_type":"notify","sub_type":"poke",
//	 "target_id":<好友>, "user_id":<好友>, "sender_id":<bot>, "self_id":<bot>}
//
// 即 user_id 与 target_id 都是被戳方, 发起方在 sender_id —— 与 OneBot 11 标准
// "user_id 是发起方" 不同。归一化之后 SenderId 一定是发起方。
func TestNormalizePoke(t *testing.T) {
	cases := []struct {
		name     string
		in       event.NoticeNotifyPoke
		wantFrom int
	}{
		{
			name: "NapCat: bot 戳好友 (NapCat 回推的事件)",
			in: event.NoticeNotifyPoke{
				NoticeNotify: event.NoticeNotify{
					UserId:   982809597, // 被戳方
					TargetId: 982809597, // 也是被戳方
				},
				NoticeNotifyPoke_Nc: event.NoticeNotifyPoke_Nc{SenderId: 2393827810},
			},
			wantFrom: 2393827810,
		},
		{
			name: "NapCat: 好友戳 bot",
			in: event.NoticeNotifyPoke{
				NoticeNotify: event.NoticeNotify{
					UserId:   2393827810, // 被戳方 = bot
					TargetId: 2393827810,
				},
				NoticeNotifyPoke_Nc: event.NoticeNotifyPoke_Nc{SenderId: 982809597},
			},
			wantFrom: 982809597,
		},
		{
			name: "标准实现: 无 sender_id, user_id 即发起方",
			in: event.NoticeNotifyPoke{
				NoticeNotify: event.NoticeNotify{
					UserId:   982809597,
					TargetId: 2393827810,
				},
			},
			wantFrom: 982809597,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := c.in
			normalizePoke(&p)
			if p.SenderId != c.wantFrom {
				t.Errorf("SenderId = %d, 期望 %d", p.SenderId, c.wantFrom)
			}
		})
	}

	normalizePoke(nil) // 不能因为 nil 崩
}

// TestClassifyPokeEvent 走真实分发路径: 造一个 NapCat 风格的 poke 事件,
// 断言回调收到的 SenderId 是发起方, 而不是被戳方。
func TestClassifyPokeEvent(t *testing.T) {
	const (
		botId    = 2393827810
		friendId = 982809597
	)

	newBot := func() *Bot {
		b := New()
		b.msgProcessEnabled.Store(true)
		// 关掉事件日志输出: 它会把事件格式化 (FormatEvent) 并查 Redis
		b.settings.msgLogOut = false
		return b
	}

	// 场景一: bot 自己发起 poke, NapCat 回推 —— 回调应能看到 sender=bot
	got := make(chan int, 4)
	b := newBot()
	b.OnNoticeNotifyPoke(func(p *event.NoticeNotifyPoke) {
		got <- p.SenderId
	})

	e := &event.Event{
		Time:       1,
		SelfId:     botId,
		PostType:   string(event.TYPE_L1_NOTICE),
		NoticeType: string(event.TYPE_L2_NOTICE_NOTIFY),
		SubType:    string(event.TYPE_L3_NOTICE_NOTIFY_POKE),
		UserId:     friendId, // NapCat: 被戳方
		TargetId:   friendId,
		SenderId:   botId, // NapCat: 发起方
	}
	e.InitFields()
	if e.TypeL1 != event.TYPE_L1_NOTICE || e.TypeL2 != event.TYPE_L2_NOTICE_NOTIFY || e.TypeL3 != event.TYPE_L3_NOTICE_NOTIFY_POKE {
		t.Fatalf("事件类型判定不对: L1=%s L2=%s L3=%s", e.TypeL1, e.TypeL2, e.TypeL3)
	}
	b.classifyEvent(e)

	select {
	case from := <-got:
		if from != botId {
			t.Errorf("SenderId = %d, 期望 %d (发起方是 bot)", from, botId)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("poke 回调没有被触发")
	}
}
