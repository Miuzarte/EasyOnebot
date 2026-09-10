package EasyOnebot

import (
	"encoding/json"
	"testing"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/message"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jinzhu/copier"
)

func TestCopier(t *testing.T) {
	m := &event.Event{
		Time:        123456789,
		SelfId:      987654321,
		PostType:    "message",
		MessageType: "group",
		SubType:     "normal",
		MessageId:   -555555555,
		UserId:      234567891,
		Message: message.SegmentArray{
			message.Segment{
				Type: "text",
				Data: map[string]any{
					"text": "这是群消息",
				},
			},
			message.Segment{
				Type: "image",
				Data: map[string]any{
					"file": "https://baidu.com/1.jpg",
				},
			},
		},
		RawMessage: "这是群消息[CQ:image:file=https://baidu.com/1.jpg]",
		Font:       233,
		GroupId:    345678912,
		Anonymous:  &event.Anonymous{},
		Sender: event.GroupSender{
			Sender: event.Sender{
				UserId:   234567891,
				Nickname: "any name",
				Sex:      "unknown",
				Age:      0,
			},
			Card:  "cardName",
			Area:  "not asia",
			Level: "100",
			Role:  "owner",
			Title: "群主",
		},
	}
	a := &event.MessageGroup{}
	err := copier.Copy(a, m)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", a)
}

func FromEventMS(input *event.Event, output any) error {
	return mapstructure.Decode(input, output)
}

func TestStructToStruct(t *testing.T) {
	m := &event.Event{
		Time:        123456789,
		SelfId:      987654321,
		PostType:    "message",
		MessageType: "group",
		SubType:     "normal",
		MessageId:   -555555555,
		UserId:      234567891,
		Message: message.SegmentArray{
			message.Segment{
				Type: "text",
				Data: map[string]any{
					"text": "这是群消息",
				},
			},
			message.Segment{
				Type: "image",
				Data: map[string]any{
					"file": "https://baidu.com/1.jpg",
				},
			},
		},
		RawMessage: "这是群消息[CQ:image:file=https://baidu.com/1.jpg]",
		Font:       233,
		GroupId:    345678912,
		Anonymous:  &event.Anonymous{},
		Sender: event.GroupSender{
			Sender: event.Sender{
				UserId:   234567891,
				Nickname: "any name",
				Sex:      "unknown",
				Age:      0,
			},
			Card:  "cardName",
			Area:  "not asia",
			Level: "100",
			Role:  "owner",
			Title: "群主",
		},
	}
	a := &event.MessageGroup{}
	err := FromEventMS(m, a)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", a)
}

func TestDecodeFromJson(t *testing.T) {
	input := any(
		map[string]any{
			"message_id": 234,
		},
	)
	output := &api.SendPrivateMsgResp{}

	err := mapstructure.Decode(input, output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)
}

func TestDecodeFromJsonArray(t *testing.T) {
	input := any(
		[]map[string]any{
			{
				"type": "reply",
				"data": map[string]any{
					"id": "123456",
				},
			},
			{
				"type": "text",
				"data": map[string]any{
					"text": "纯文本内容",
				},
			},
			{
				"type": "image",
				"data": map[string]any{
					"file": "https://baidu.com/1.jpg",
				},
			},
		},
	)
	output := &message.SegmentArray{}

	err := mapstructure.Decode(input, output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)
}

func TestDecodeFromStringArray(t *testing.T) {
	input := any(
		[]string{
			"123", "456", "789",
		},
	)
	output := &[]string{}

	err := mapstructure.Decode(input, output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)
}

const myJson = `{
	"message_type": "group",
	"sub_type": "normal",
	"message_id": 907599783,
	"group_id": 612645549,
	"user_id": 982809597,
	"anonymous": null,
	"message": "\u8FD9\u662F\u7FA4\u6D88\u606F",
	"raw_message": "\u8FD9\u662F\u7FA4\u6D88\u606F",
	"font": 0,
	"sender": {
		"user_id": 982809597,
		"nickname": "\u8B2C\u7D17\u7279 \u309A\u309A\u309A\u309A\u309A\u309A\u309A\u309A",
		"card": "",
		"sex": "unknown",
		"age": 0,
		"area": "",
		"level": "100",
		"role": "owner",
		"title": ""
	},
	"time": 1714386431,
	"self_id": 2393827810,
	"post_type": "message"
}`

func TestDecodeFromBytes(t *testing.T) {
	// mapstructure 不认 []byte, 要先过一遍 json
	input := map[string]any{}
	if err := json.Unmarshal([]byte(myJson), &input); err != nil {
		t.Fatal(err)
	}
	output := map[string]any{}
	if err := mapstructure.Decode(input, &output); err != nil {
		t.Error(err)
	}
	t.Logf("%#v\n", output)
}

func TestJsonUnmarshalToMap(t *testing.T) {
	input := []byte(myJson)
	output := map[string]any{}

	err := json.Unmarshal(input, &output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)
}

const myJson2 = `{
"number1": 123456789,
"number2": 123.456
}`

func TestJsonUnmarshalToMapThenToStruct(t *testing.T) {
	input := []byte(myJson2)
	output := map[string]any{}

	err := json.Unmarshal(input, &output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)

	var output2 struct {
		Number1 int
		Number2 float64
	}

	err = mapstructure.Decode(output, &output2)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output2)
}

func TestPassMapPtrToDecode(t *testing.T) {
	input := map[string]any{}
	input["message_id"] = 234
	output := &api.SendPrivateMsgResp{}

	err := mapstructure.Decode(&input, output)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#v\n", output)
}
