package EasyOnebot

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"

	"github.com/redis/rueidis"
)

var rdb redisClient

func init() {
	rueidisClient, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress:  []string{"192.168.1.104:6379"},
		DisableCache: true, // Disable Client-Side Caching
	})
	if err != nil {
		panic(err)
	}
	rdb = redisClient{Client: rueidisClient}
}

func TestDelSample(t *testing.T) {
	const match = "sample:*"
	var cursor uint64
	for {
		cmd := rdb.Client.B().Scan().Cursor(cursor).Match(match).Build()
		se, err := rdb.Client.Do(t.Context(), cmd).AsScanEntry()
		if err != nil {
			t.Logf("Error scanning keys: %v", err)
			return
		}
		t.Log("Deleting key:")
		for _, key := range se.Elements {
			println(key)
		}

		// 删除匹配的键
		if len(se.Elements) > 0 {
			_, err = rdb.Client.Do(context.Background(), rdb.Client.B().Del().Key(se.Elements...).Build()).ToInt64()
			if err != nil {
				t.Logf("Error deleting keys: %v", err)
				return
			}
		}

		// 如果 cursor 为 0，表示已经遍历完所有匹配的键
		if cursor == 0 {
			break
		}
	}
}

func TestSet(t *testing.T) {
	rdb.set(t.Context(), "key", "value", 20*time.Second)
}

var testMsgP = &event.MessagePrivate{
	MessageBase: event.MessageBase{
		Base: event.Base{
			Time:     int(time.Now().Unix()),
			SelfId:   123456789,
			PostType: "message",
		},
		MessageType: "private",
		SubType:     "friend",
		MessageId:   -123456789,
		UserId:      987654321,
		Message: message.SegmentArray{
			{Type: "text", Data: map[string]any{"text": "[第一部分]"}},
			{Type: "image", Data: map[string]any{"file": "123.jpg"}},
			{Type: "text", Data: map[string]any{"text": "图片之后的部分，表情："}},
			{Type: "face", Data: map[string]any{"id": "123"}},
		},
		RawMessage: `&#91;第一部分&#93;[CQ:image,file=123.jpg]图片之后的部分，表情：[CQ:face,id=123]`,
		Font:       233,
	},
	Sender: event.Sender{
		UserId:   987654321,
		Nickname: "nickname",
		Sex:      "unknown",
		Age:      18,
	},
}

func TestJson(t *testing.T) {
	msgBefore := fmt.Sprintf("%+v", testMsgP)
	t.Log(msgBefore)
	data, err := json.Marshal(testMsgP)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	err = rdb.set(t.Context(), "message_private:-123456789", data, time.Minute)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	result := rdb.get(t.Context(), "message_private:-123456789")
	if result.Err != nil {
		t.Error(err)
		t.FailNow()
	}
	mp, err := result.MessagePrivate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	msgAfter := fmt.Sprintf("%+v", mp)
	t.Log(msgAfter)

	if msgBefore != msgAfter {
		t.Error("MessagePrivate not equal")
	}
}

func TestSetMessagePrivate(t *testing.T) {
	before := fmt.Sprintf("%+v", testMsgP)
	key := "message_private:" + itoa(testMsgP.UserId) + ":" + itoa(testMsgP.MessageId)
	err := rdb.set(t.Context(), key, testMsgP, 0)
	rdb.set(t.Context(), "message_private:-123456789", "value", time.Minute)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	msg, err := rdb.get(t.Context(), key).MessagePrivate()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	after := fmt.Sprintf("%+v", msg)
	t.Log(before == after)
}

func TestHSet(t *testing.T) {
	// rdb.Client.HSet(context.Background(), "hash", "key", "value")
	// rdb.Client.Expire(context.Background(), "hash", 20*time.Second)
	rdb.Client.Do(context.Background(), rdb.Client.B().Hset().Key("hash").FieldValue().FieldValue("key", "value").Build())
	rdb.Client.Do(context.Background(), rdb.Client.B().Expire().Key("hash").Seconds(20).Build())
}
