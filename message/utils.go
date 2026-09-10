package message

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Miuzarte/EasyOnebot/internal/utils"
)

func numTToString[numT NumberT](num numT) string {
	switch u := any(num).(type) {
	case string:
		return u
	default:
		return fmt.Sprint(num)
	}
}

// ConstraintMessage 约束消息类型为 [MessageT]
func ConstraintMessage(msg any) any {
	switch m := msg.(type) {
	case string:

	case Segment:
	case []Segment:
	case SegmentArray:

	case *Segment:
		msg = *m
	case *[]Segment:
		msg = *m
	case *SegmentArray:
		msg = *m

	case fmt.Stringer:
		msg = m.String()
	case error:
		msg = "[ERROR] " + m.Error()
	case []byte:
		msg = string(m)
	}
	return msg
}

// UnmarshalMessage 转换消息类型为 [SegmentArray]
func UnmarshalMessage(msg any) (segs SegmentArray) {
	switch m := msg.(type) {
	case string:
		segs = ParseCqCodes(m)

	case Segment:
		segs = SegmentArray{m}
	case []Segment:
		segs = SegmentArray(m)
	case SegmentArray:
		segs = m

	case *Segment:
		segs = SegmentArray{*m}
	case *[]Segment:
		segs = SegmentArray(*m)
	case *SegmentArray:
		segs = *m

	default:
		segs = SegmentArray{Segment{
			Type: "text",
			Data: map[string]any{"text": utils.MarshalString(msg)},
		}}
	}
	return segs
}

func AssertSegmentArray(a any) (segs SegmentArray) {
	switch a := a.(type) {
	case []any:
		for _, v := range a {
			segs = append(segs, AssertSegmentArray(v)...)
		}
	case map[string]any:
		if typ, ok := a["type"].(string); ok {
			if data, ok := a["data"].(map[string]any); ok {
				segs.Append(Segment{
					Type: typ,
					Data: data,
				})
			}
		}
	case string:
		segs = ParseCqCodes(a)
	default:
		panic(fmt.Errorf("unsupported type to assert: %T", a))
	}
	return
}

// ParseCqCodes 从 RawMessage 中解析消息段
func ParseCqCodes(raw string) (segs SegmentArray) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
			panic(raw)
		}
	}()

	stack := utils.Stack[int]{} // 泛型栈 存储 '[' 的索引
	lastIndex := 0              // 记录最后 ']' 的下一个位置

	for i, char := range raw {
		switch char {
		case '[':
			stack.Push(i)

			// 中括号外的内容提取为 [Text]
			if lastIndex < i {
				text := raw[lastIndex:i]
				segs.Append(Segment{
					Type: "text",
					Data: map[string]any{"text": utils.MarshalString(text)},
				})
			}

		case ']':
			start, ok := stack.Pop()
			if !ok { // 非闭合
				continue
			}

			// 带着中括号去解析
			segs.Append(ParseCqCode(raw[start : i+1]))

			lastIndex = i + 1
		}
	}

	if lastIndex < len(raw) {
		// 剩余的文本
		text := raw[lastIndex:]
		if text != "" {
			segs.Append(Segment{
				Type: "text",
				Data: map[string]any{"text": utils.MarshalString(text)},
			})
		}
	}
	if start, ok := stack.Pop(); ok {
		_ = start
		// 不应存在未闭合的中括号 ?
		// panic(fmt.Sprintf("unclosed brackets at %d: %s\nraw: %s", start, raw[start:], raw))
	}

	return
}

// ParseCqCode 从 CQ 码中解析出消息段,
// 必须以中括号包裹, 否则直接返回文本消息段
func ParseCqCode(msg string) Segment {
	// l := strings.HasPrefix(msg, "[")
	// r := strings.HasSuffix(msg, "]")
	l := msg[0] == '['
	r := msg[len(msg)-1] == ']'
	if l && r { // 去除中括号
		msg = msg[1 : len(msg)-1]
	} else {
		return Segment{
			Type: "text",
			Data: map[string]any{"text": utils.MarshalString(utils.UnescapeCqCode(msg))},
		}
	}

	// CQ:image,

	i := strings.Index(msg, ":")
	if i == -1 || i == len(msg)-1 {
		return Segment{
			Type: "text",
			Data: map[string]any{"text": utils.MarshalString(utils.UnescapeCqCode(msg))},
		}
	}

	// 第一个逗号
	j := strings.Index(msg, ",")

	if j == -1 { // 没有参数
		j = len(msg)
	}

	if i > j {
		panic(fmt.Errorf("unexpected message format: %s", msg))
	}

	seg := Segment{
		Type: msg[i+1 : j],
	}

	for param := range strings.SplitSeq(msg[i+1:], ",") {
		if seg.Data == nil {
			seg.Data = map[string]any{} // value 的类型必为 string
		}
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			k := kv[0]
			v := utils.UnescapeCqCode(kv[1])
			seg.Data[k] = tryAtoi(k, v)
		}
	}

	return seg
}

func tryAtoi(k, v string) any {
	var i int
	isNumber := utils.IsNumber(v)
	if isNumber {
		i, _ = strconv.Atoi(v)
	}

	switch k {
	case "qq", "uin", "user_id": // known number value
		if !isNumber {
			if v == "all" {
				return -1 // TODO: constant
			} else {
				panic(fmt.Errorf("[FIXME] unknown \"%s\" value: %s", k, v))
			}
		}
	}

	if isNumber {
		return i
	} else {
		return v
	}
}
