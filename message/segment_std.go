package message

import (
	"fmt"

	"github.com/Miuzarte/EasyOnebot/internal/utils"
)

type SegmentText Segment

/*
Text 纯文本
*/
func Text(text any) Segment {
	return Segment{
		Type: TYPE_TEXT,
		Data: map[string]any{"text": utils.MarshalString(text)},
	}
}

type SegmentFace Segment

/*
Face QQ 表情
*/
func Face[numT NumberT](id numT) Segment {
	return Segment{
		Type: TYPE_FACE,
		Data: map[string]any{"id": id},
	}
}

type SegmentImage Segment

/*
Image 图片

	file: 图片文件名
		case string: 直接发送
		case []byte: 图片数据 编码为 base64
		default: fmt.Sprint(file)

发送时，`file` 参数除了支持使用收到的图片文件名直接发送外，还支持：

- 绝对路径，例如 `file:///C:\\Users\Richard\Pictures\1.png`，格式使用 [`file` URI](https://tools.ietf.org/html/rfc8089)

- 网络 URL，例如 `http://i1.piimg.com/567571/fdd6e7b6d93f1ef0.jpg`

- Base64 编码，例如 `base64://iVBORw0KGgoAAAANSUhEUgAAABQAAAAVCAIAAADJt1n/AAAAKElEQVQ4EWPk5+RmIBcwkasRpG9UM4mhNxpgowFGMARGEwnBIEJVAAAdBgBNAZf+QAAAAABJRU5ErkJggg==`

	otherParams:
	"type": 图片类型，`flash` 表示闪照，无此参数表示普通图片
	"cache": 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	"proxy": 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	"timeout": 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间，默认不超时
*/
func Image[fileT FileT](file fileT, otherParams ...map[string]any) Segment {
	seg := Segment{Type: TYPE_IMAGE, Data: map[string]any{}}
	seg.Data["file"] = utils.MarshalFile(file)
	for _, p := range otherParams {
		for k, v := range p {
			if v != nil {
				seg.Data[k] = v
			}
		}
	}
	return seg
}

type SegmentRecord Segment

/*
Record 语音

	file: 语音文件名
		case string: 直接发送
		case []byte: 语音数据 编码为 base64
		default: fmt.Sprint(file)

发送时，`file` 参数除了支持使用收到的语音文件名直接发送外，还支持其它形式，参考 [Image]

	otherParams:
	"magic": 发送时可选，默认 `0`，设置为 `1` 表示变声
	"cache": 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	"proxy": 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	"timeout": 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
*/
func Record[fileT FileT](file fileT, otherParams ...map[string]any) Segment {
	seg := Segment{Type: TYPE_RECORD, Data: map[string]any{}}
	seg.Data["file"] = utils.MarshalFile(file)
	for _, p := range otherParams {
		for k, v := range p {
			if v != nil {
				seg.Data[k] = v
			}
		}
	}
	return seg
}

type SegmentVideo Segment

/*
Video 短视频

	file: 视频文件名
		case string: 直接发送
		case []byte: 视频数据 编码为 base64
		default: fmt.Sprint(file)

发送时，`file` 参数除了支持使用收到的视频文件名直接发送外，还支持其它形式，参考 [Image]

	otherParams:
	"cache": 只在通过网络 URL 发送时有效，表示是否使用已缓存的文件，默认 `1`
	"proxy": 只在通过网络 URL 发送时有效，表示是否通过代理下载文件（需通过环境变量或配置文件配置代理），默认 `1`
	"timeout": 只在通过网络 URL 发送时有效，单位秒，表示下载网络文件的超时时间 ，默认不超时
*/
func Video[fileT FileT](file fileT, otherParams ...map[string]any) Segment {
	seg := Segment{Type: TYPE_VIDEO, Data: map[string]any{}}
	seg.Data["file"] = utils.MarshalFile(file)
	for _, p := range otherParams {
		for k, v := range p {
			if v != nil {
				seg.Data[k] = v
			}
		}
	}
	return seg
}

type SegmentAt Segment

/*
At @某人

	qq: @的 QQ 号，`all` 表示全体成员
*/
func At[numT NumberT](qq numT) Segment {
	return Segment{
		Type: TYPE_AT,
		Data: map[string]any{"qq": fmt.Sprint(qq)},
	}
}

type SegmentRps Segment

/*
Rps 猜拳魔法表情
*/
func Rps() Segment {
	return Segment{Type: TYPE_RPS}
}

type SegmentDice Segment

/*
Dice 掷骰子魔法表情
*/
func Dice() Segment {
	return Segment{Type: TYPE_DICE}
}

type SegmentShake Segment

/*
Shake 窗口抖动（戳一戳）
*/
func Shake() Segment {
	return Segment{Type: TYPE_SHAKE}
}

type SegmentPoke Segment

/*
Poke 戳一戳

	typ: 类型
	id: ID

戳一戳
(1, -1)

比心
(2, -1)

点赞
(3, -1)

心碎
(4, -1)

666
(5, -1)

放大招
(6, -1)

宝贝球 (SVIP)
(126, 2011)

玫瑰花 (SVIP)
(126, 2007)

召唤术 (SVIP)
(126, 2006)

让你皮 (SVIP)
(126, 2009)

结印 (SVIP)
(126, 2005)

手雷 (SVIP)
(126, 2004)

勾引
(126, 2003)

抓一下 (SVIP)
(126, 2001)

碎屏 (SVIP)
(126, 2002)

敲门 (SVIP)
(126, 2002)
*/
func Poke[numT NumberT](typ, id numT) Segment {
	return Segment{
		Type: TYPE_POKE,
		Data: map[string]any{"type": typ, "id": id},
	}
}

type SegmentShare Segment

/*
Share 链接分享

	url: URL
	title: 标题

	otherParams:
	"content": 发送时可选，内容描述
	"image": 发送时可选，图片 URL
*/
func Share(url, title string, otherParams ...map[string]any) Segment {
	seg := Segment{Type: TYPE_SHARE, Data: map[string]any{"url": url, "title": title}}
	for _, p := range otherParams {
		for k, v := range p {
			if v != nil {
				seg.Data[k] = v
			}
		}
	}
	return seg
}

type SegmentContact Segment

/*
Contact 推荐好友/群

	typ: 推荐好友/群 | `qq`/`group`
	id: 被推荐人的 QQ 号/被推荐群的群号
*/
func Contact[numT NumberT](typ string, id numT) Segment {
	return Segment{
		Type: TYPE_CONTACT,
		Data: map[string]any{"type": typ, "id": id},
	}
}

type SegmentLocation Segment

/*
Location 位置

	lat: 纬度
	lon: 经度
*/
func Location[numT NumberT](lat, lon numT) Segment {
	return Segment{
		Type: TYPE_LOCATION,
		Data: map[string]any{"lat": lat, "lon": lon},
	}
}

type SegmentMusic Segment

/*
Music 音乐分享

	typ: 分别表示使用 QQ 音乐、网易云音乐、虾米音乐、自定义分享 | `qq` `163` `xm` `custom`
	id: 歌曲 ID
*/
func Music[numT NumberT](typ string, id numT) Segment {
	return Segment{
		Type: TYPE_MUSIC,
		Data: map[string]any{"type": typ, "id": id},
	}
}

type SegmentReply Segment

/*
Reply 回复

id: 回复时引用的消息 ID
*/
func Reply[numT NumberT](id numT) Segment {
	return Segment{
		Type: TYPE_REPLY,
		Data: map[string]any{"id": fmt.Sprint(id)},
	}
}

type SegmentNode Segment

/*
Node1 合并转发节点

	id: 转发的消息 ID
*/
func Node1[numT NumberT](id numT) Segment {
	return Segment{
		Type: TYPE_NODE,
		Data: map[string]any{"id": id},
	}
}

/*
Node3 合并转发自定义节点

	user_id: 发送者 QQ 号
	nickname: 发送者昵称
	content: 消息内容，支持发送消息时的 `message` 数据类型
*/
func Node3[msgT MessageT, numT NumberT](userId numT, nickname string, content msgT) Segment {
	return Segment{
		Type: TYPE_NODE,
		Data: map[string]any{
			"user_id":  numTToString(userId), // 需要字符串类型: NapCat 的合并转发节点要求 uin 是字符串
			"nickname": nickname,
			"content":  content,
		},
	}
}

// Node3Any 合并转发自定义节点
//
// 用于发送 [event] 中的 `Message any` 数据
func Node3Any[numT NumberT](userId numT, nickname string, content any) Segment {
	return Segment{
		Type: TYPE_NODE,
		Data: map[string]any{
			"user_id":  fmt.Sprint(userId),
			"nickname": nickname,
			"content":  content,
		},
	}
}

type SegmentXml Segment

/*
Xml XML 消息

	data: XML 内容
*/
func Xml(data string) Segment {
	return Segment{
		Type: TYPE_XML,
		Data: map[string]any{"data": data},
	}
}

type SegmentJson Segment

/*
Json JSON 消息

	data: JSON 内容
*/
func Json(data string) Segment {
	return Segment{
		Type: TYPE_JSON,
		Data: map[string]any{"data": data},
	}
}
