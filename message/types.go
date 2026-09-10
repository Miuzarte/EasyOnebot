package message

import (
	"fmt"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/Miuzarte/EasyOnebot/internal/utils"

	"github.com/go-viper/mapstructure/v2"
)

// MessageT 消息类型
type MessageT interface {
	string | Segment | SegmentArray
}

// FileT 文件类型
type FileT interface {
	[]byte | // base64 编码后发送
		string // url 直接发送 / base64
}

type IntegerT interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

type FloatPointT interface {
	~float32 | ~float64
}

// NumberT 数字类型
type NumberT interface {
	string | // 由调用方保证字符串内容为数字
		IntegerT | FloatPointT
}

type SegType = string

const (
	TYPE_TEXT     SegType = "text"
	TYPE_FACE     SegType = "face"
	TYPE_IMAGE    SegType = "image"
	TYPE_RECORD   SegType = "record"
	TYPE_VIDEO    SegType = "video"
	TYPE_AT       SegType = "at"
	TYPE_RPS      SegType = "rps"
	TYPE_DICE     SegType = "dice"
	TYPE_SHAKE    SegType = "shake"
	TYPE_POKE     SegType = "poke"
	TYPE_SHARE    SegType = "share"
	TYPE_CONTACT  SegType = "contact"
	TYPE_LOCATION SegType = "location"
	TYPE_MUSIC    SegType = "music"
	TYPE_REPLY    SegType = "reply"
	TYPE_NODE     SegType = "node"
	TYPE_XML      SegType = "xml"
	TYPE_JSON     SegType = "json"
)

// Segment 消息段
type Segment struct {
	Type SegType        `json:"type"`
	Data map[string]any `json:"data"`
}

func (seg Segment) String() string {
	return seg.ToString(false)
}

func (seg Segment) ToString(cutBase64 bool) string {
	if seg.Type == "text" {
		if text, ok := seg.Data["text"].(string); ok {
			return text
		}
	}
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(seg.Type)
	for k, v := range seg.Data {
		var s string
		switch v := v.(type) {
		case string:
			s = v
		case []any:
			segs := SegmentArray{}
			err := mapstructure.Decode(v, &segs)
			if err != nil {
				s = fmt.Sprintf(
					"failed to decode seg data %s(%T): %+v\n%s\n%s",
					k, v, err,
					utils.ToJson(seg.Data),
					debug.Stack(),
				)
				break
			}
			s = segs.ToString(cutBase64)
		default:
			s = fmt.Sprintf(
				"unexpected type of seg data: %T\n%s",
				v,
				utils.ToJson(seg.Data),
			)
		}

		switch k {
		case "file": // 将 base64 结果截断, 简化打印
			if cutBase64 && len(s) > 9+16 &&
				strings.HasPrefix(s, "base64://") {
				s = s[:9+16]
			}
		}
		sb.WriteString(",")
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(s)
	}
	sb.WriteString("]")
	return sb.String()
}

// SegmentArray 消息段数组
type SegmentArray []Segment

func (segs *SegmentArray) Append(ms ...Segment) *SegmentArray {
	if segs == nil {
		segs = &SegmentArray{}
	}
	*segs = append(*segs, ms...)
	return segs
}

func (segs *SegmentArray) DeleteAt(qq ...int) *SegmentArray {
	*segs = slices.DeleteFunc(*segs, func(seg Segment) bool {
		if seg.Type != "at" {
			return false
		}
		qq0, ok := seg.Data["qq"].(int)
		if !ok {
			return false
		}
		return slices.Contains(qq, qq0)
	})
	return segs
}

func (segs *SegmentArray) DeleteType(typ ...SegType) *SegmentArray {
	if len(typ) == 0 {
		return segs
	}
	*segs = slices.DeleteFunc(*segs, CreateCompareFunc(typ...))
	return segs
}

func (segs *SegmentArray) String() string {
	return segs.ToString(false)
}

func (segs *SegmentArray) ToString(cutBase64 bool) string {
	switch len(*segs) {
	case 0:
		return ""
	case 1:
		return (*segs)[0].ToString(cutBase64)
	case 2:
		return (*segs)[0].ToString(cutBase64) + (*segs)[1].ToString(cutBase64)
	default:
		sb := strings.Builder{}
		for _, ms := range *segs {
			sb.WriteString(ms.ToString(cutBase64))
		}
		return sb.String()
	}
}

func (segs *SegmentArray) TryAtoi() {
	for i := range *segs {
		m := (*segs)[i].Data
		for k, va := range m {
			vs, ok := va.(string)
			if ok {
				m[k] = tryAtoi(k, vs)
			}
		}
	}
}

func (segs *SegmentArray) TypeSet() []SegType {
	set := utils.Set[SegType]{}
	for i := range *segs {
		set.Add((*segs)[i].Type)
	}
	return set.Get()
}

func (segs *SegmentArray) GetFirstType(typ SegType) *Segment {
	cf := CreateCompareFunc(typ)
	for i := range *segs {
		if cf((*segs)[i]) {
			return &(*segs)[i]
		}
	}
	return nil
}

func (segs *SegmentArray) GetType(typ ...SegType) (res SegmentArray) {
	if len(typ) == 0 {
		return nil
	}
	cf := CreateCompareFunc(typ...)
	for _, seg := range *segs {
		if cf(seg) {
			res = append(res, seg)
		}
	}
	return
}

func (segs *SegmentArray) WithType(typ ...SegType) bool {
	if len(typ) == 0 {
		return false
	}
	return slices.ContainsFunc(*segs, CreateCompareFunc(typ...))
}

func (segs *SegmentArray) WithoutType(typ ...SegType) bool {
	return !segs.WithType(typ...)
}

func (segs *SegmentArray) OnlyType(typ ...SegType) bool {
	if len(typ) == 0 {
		return false
	}
	typeSet := segs.TypeSet()
	if len(typ) != len(typeSet) {
		return false
	}
	slices.Sort(typ)
	slices.Sort(typeSet)
	return slices.Equal(typ, typeSet)
}

// CreateCompareFunc 创建 [SegType] 比较函数,
// 根据 typ 的数量返回不同的匹配实现
func CreateCompareFunc(typ ...SegType) func(Segment) bool {
	switch len(typ) {
	case 0:
		return func(seg Segment) bool {
			return false
		}
	case 1:
		typ0 := typ[0]
		return func(seg Segment) bool {
			return seg.Type == typ0
		}
	default: // hash map
		set := utils.Set[SegType]{}
		set.Add(typ...)
		return func(seg Segment) bool {
			return set.Ok(seg.Type)
		}
	}
}
