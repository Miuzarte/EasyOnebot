package message

// https://lagrange-onebot.apifox.cn/236981924e0

const ( // lagrange
	TYPE_FORWARD  SegType = "forward"
	TYPE_LONGMSG  SegType = "longmsg"
	TYPE_MFACE    SegType = "mface"
	TYPE_MARKDOWN SegType = "markdown"
	TYPE_KEYBOARD SegType = "keyboard"
	TYPE_FILE     SegType = "file"
)

// RichTextSeg 富文本消息段 (doc only)
type RichTextSeg struct {
	// [At]
	// [Face]
	// [Image]
	// [Reply]
	// [Text]
}

// SingleSeg 单个消息段 (doc only)
type SingleSeg struct {
	// [Dict]
	// [Forward]
	// [Json]
	// [Location]
	// [Longmsg]
	// [MFace]
	// [Music]
	// [Poke]
	// [Record]
	// [Rps]
	// [Video]
}

// MarkdownSeg Markdown消息段 (doc only)
type MarkdownSeg struct {
	// [Markdown]
	// [Keyboard]
}

// Dict 骰子 dict????
func Dict() Segment {
	return Segment{Type: "dict"}
}

/*
Forward 转发

	id: 转发 ID
*/
func Forward(id string) Segment {
	return Segment{
		Type: TYPE_FORWARD,
		Data: map[string]any{"id": id},
	}
}

/*
Json Json

std: [Json]

	data: Json 数据
*/
// func Json(data string) Segment

/*
Location_Lgr 定位

std: [Location]

	lat: 纬度
	lon: 经度
*/
func Location_Lgr(lat, lon, title, content string) Segment {
	return Segment{
		Type: TYPE_LOCATION,
		Data: map[string]any{
			"lat":     lat,
			"lon":     lon,
			"title":   title,
			"content": content,
		},
	}
}

/*
Longmsg 长消息

	id: 长消息 ID
*/
func Longmsg(id string) Segment {
	return Segment{
		Type: TYPE_LONGMSG,
		Data: map[string]any{"id": id},
	}
}

/*
MFace 商城表情

	emojiPkgId: 表情包 ID
	emojiId: 表情 ID
	key: 表情 Key
	summary: 表情说明
	url(可选): 表情 Url
*/
func MFace(emojiPkgId int, emojiId, key, summary string, url ...string) Segment {
	if len(url) == 0 {
		return Segment{
			Type: TYPE_MFACE,
			Data: map[string]any{
				"emoji_package_id": emojiPkgId,
				"emoji_id":         emojiId,
				"key":              key,
				"summary":          summary,
			},
		}
	} else {
		return Segment{
			Type: TYPE_MFACE,
			Data: map[string]any{
				"emoji_package_id": emojiPkgId,
				"emoji_id":         emojiId,
				"key":              key,
				"summary":          summary,
				"url":              url[0],
			},
		}
	}
}

/*
Music 音乐

std: [Music]

	typ: 音乐类型
	url: 跳转 Url
	audio: 音乐 Url
	title: 标题
	content: 内容
	image: 图片
*/
func Music_Lgr(typ, url, audio, title, content, image string) Segment {
	return Segment{
		Type: TYPE_MUSIC,
		Data: map[string]any{
			"type":    typ,
			"url":     url,
			"audio":   audio,
			"title":   title,
			"content": content,
			"image":   image,
		},
	}
}

/*
Poke_Lgr 戳一戳

std: [Poke]

	typ: 戳一戳类型
	id: 戳一戳 ID
	strength(可选): 戳一戳强度 <默认值: 0>
*/
func Poke_Lgr(typ, id string, strength ...string) Segment {
	if len(strength) == 0 {
		return Segment{
			Type: TYPE_POKE,
			Data: map[string]any{
				"type": typ,
				"id":   id,
			},
		}
	} else {
		return Segment{
			Type: TYPE_POKE,
			Data: map[string]any{
				"type":     typ,
				"strength": strength[0],
				"id":       id,
			},
		}
	}
}

/*
Record 语音

std: [Video]
*/
// func Record[fileT File](file fileT) Segment

/*
Rps 猜拳

std: [Rps]
*/
// func Rps() Segment

/*
Video 视频

std: [Video]
*/
// func Video[fileT File](file fileT) Segment

/*
At_Lgr At

std: [At]

	userId: 用户 Uin
	name(已废弃): 显示的文本
*/
func At_Lgr[numT NumberT](userId numT, name string) Segment {
	return Segment{
		Type: TYPE_AT,
		Data: map[string]any{
			"user_id": userId,
			// "name": name,
		},
	}
}

/*
Face_Lgr 表情

std: [Face]

	id: 表情 ID
	large(可选): 是否大表情
*/
func Face_Lgr[numT NumberT](id numT, large ...bool) Segment {
	if len(large) == 0 {
		return Segment{
			Type: TYPE_FACE,
			Data: map[string]any{
				"id": id,
			},
		}
	} else {
		return Segment{
			Type: TYPE_FACE,
			Data: map[string]any{
				"id":    id,
				"large": large[0],
			},
		}
	}
}

/*
Image_Lgr 图片

std: [Image]

	file: 图片链接, 支持 http/https/file/base64
	filename: 图片名称
	url: 图片链接
	summary: 图片说明
	subType: 图片子类型
*/
func Image_Lgr[fileT FileT](file fileT, filename, url, summary, subType string) Segment {
	return Segment{
		Type: TYPE_IMAGE,
		Data: map[string]any{
			"file":     file,
			"filename": filename,
			"url":      url,
			"summary":  summary,
			"sub_type": subType,
		},
	}
}

/*
Reply 回复

std: [Reply]
*/
// func Reply[numT Number](id numT)

/*
Text 文本

std: [Text]
*/
// func Text(id numT)

/*
Markdown Markdown

	content: 内容
*/
func Markdown(content string) Segment {
	return Segment{
		Type: TYPE_MARKDOWN,
		Data: map[string]any{"content": content},
	}
}

/*
Keyboard 按钮

	content: 内容
*/
func Keyboard(content map[string]any) Segment {
	return Segment{
		Type: TYPE_KEYBOARD,
		Data: map[string]any{"content": content},
	}
}

/*
File 文件

	fileName: 文件名
	fileId: 文件ID
	fileHash: 文件Hash
	url: 下载链接
*/
func File(fileName, fileId, fileHash, url string) Segment {
	return Segment{
		Type: TYPE_FILE,
		Data: map[string]any{
			"file_name": fileName,
			"file_id":   fileId,
			"file_hash": fileHash,
			"url":       url,
		},
	}
}
