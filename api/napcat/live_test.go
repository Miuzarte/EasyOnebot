package napcat_test

// 联调测试: 对着一个真实 NapCat 实例跑一遍类型化调用, 验证严格解码在真实返回上成立。
//
// 默认跳过, 需要显式给地址才跑:
//
//	NAPCAT_WS_URL=ws://127.0.0.1:8081 NAPCAT_TEST_GROUP=612645549 NAPCAT_TEST_USER=982809597 \
//	    go test ./api/napcat/ -run TestLive -v
//
// 会真的发消息 (随后撤回), 所以只对测试用小号 + 测试群开放。

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/api/napcat"
	"github.com/Miuzarte/EasyOnebot/message"
	"github.com/coder/websocket"
)

// wsPoster 用 coder/websocket 直接把 api.Request 发出去, 并靠 echo 关联响应。
//
// 特意不复用 EasyOnebot 的 Bot: 这样可以独立验证 napcat 包, 不受 Bot 生命周期影响。
type wsPoster struct {
	conn *websocket.Conn
	mu   sync.Mutex
	pool map[string]chan *api.Response
}

func newWSPoster(ctx context.Context, url string) (*wsPoster, error) {
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	p := &wsPoster{conn: conn, pool: map[string]chan *api.Response{}}
	go p.readLoop(ctx)
	return p, nil
}

func (p *wsPoster) readLoop(ctx context.Context) {
	for {
		_, data, err := p.conn.Read(ctx)
		if err != nil {
			return
		}
		var raw map[string]json.RawMessage
		if json.Unmarshal(data, &raw) != nil {
			continue
		}
		// 无 echo 的是生命周期/事件推送, 直接忽略
		var echo string
		if e, ok := raw["echo"]; !ok || json.Unmarshal(e, &echo) != nil || echo == "" {
			continue
		}
		var resp api.Response
		if json.Unmarshal(data, &resp) != nil {
			continue
		}
		p.mu.Lock()
		ch, ok := p.pool[echo]
		p.mu.Unlock()
		if ok {
			ch <- &resp
		}
	}
}

func (p *wsPoster) PostReq(req *api.Request) (*api.Response, error) {
	ch := make(chan *api.Response, 1)
	p.mu.Lock()
	p.pool[req.Echo] = ch
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		delete(p.pool, req.Echo)
		p.mu.Unlock()
	}()

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := p.conn.Write(ctx, websocket.MessageText, data); err != nil {
		return nil, err
	}
	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("%s 调用超时", req.Action)
	}
}

// messageText 把纯文本构造成消息联合类型。
//
// 联合类型 (OB11MessageMixType / *_MessageID) 都带 MarshalJSON/UnmarshalJSON,
// 但没有导出的构造函数, 所以用一次 JSON 往返来构造 —— 这是使用生成类型时最省事的写法。
func messageText(t *testing.T, s string) napcat.OB11MessageMixType {
	t.Helper()
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("编码消息失败: %v", err)
	}
	var v napcat.OB11MessageMixType
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("构造消息联合类型失败: %v", err)
	}
	return v
}

func getMsgID(t *testing.T, id int64) napcat.GetMsgJSONBody_MessageID {
	t.Helper()
	b, _ := json.Marshal(id)
	var v napcat.GetMsgJSONBody_MessageID
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("构造 message_id 失败: %v", err)
	}
	return v
}

func delMsgID(t *testing.T, id int64) napcat.DeleteMsgJSONBody_MessageID {
	t.Helper()
	b, _ := json.Marshal(id)
	var v napcat.DeleteMsgJSONBody_MessageID
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatalf("构造 message_id 失败: %v", err)
	}
	return v
}

func TestLive(t *testing.T) {
	url := os.Getenv("NAPCAT_WS_URL")
	if url == "" {
		t.Skip("未设置 NAPCAT_WS_URL, 跳过联调")
	}
	group := os.Getenv("NAPCAT_TEST_GROUP")
	user := os.Getenv("NAPCAT_TEST_USER")
	if group == "" || user == "" {
		t.Fatal("需要同时设置 NAPCAT_TEST_GROUP 与 NAPCAT_TEST_USER")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	p, err := newWSPoster(ctx, url)
	if err != nil {
		t.Fatalf("连接 %s 失败: %v", url, err)
	}
	defer p.conn.CloseNow()
	c := napcat.Caller{Poster: p}

	// 只读端点
	login, err := c.GetLoginInfo()
	if err != nil {
		t.Fatalf("GetLoginInfo: %v", err)
	}
	t.Logf("登录号 %d (%s)", login.UserID, login.Nickname)

	groups, err := c.GetGroupList()
	if err != nil {
		t.Fatalf("GetGroupList: %v", err)
	}
	t.Logf("群 %d 个", len(groups))

	friends, err := c.GetFriendList()
	if err != nil {
		t.Fatalf("GetFriendList: %v", err)
	}
	t.Logf("好友 %d 个", len(friends))

	member, err := c.GetGroupMemberInfo(napcat.GetGroupMemberInfoJSONBody{GroupID: group, UserID: user})
	if err != nil {
		t.Fatalf("GetGroupMemberInfo: %v", err)
	}
	t.Logf("群成员 %d 昵称 %s", member.UserID, member.Nickname)

	// 发 → 查 → 撤回; 这里能验证 message_id / time / message 数组等真实类型
	// 用消息段数组 (ctx 发送路径就是这样装箱的), 覆盖 BoxMessage 的段分支
	segs, err := napcat.BoxMessage(message.SegmentArray{message.Text("NothingBot 类型化客户端联调, 可忽略")})
	if err != nil {
		t.Fatalf("BoxMessage 失败: %v", err)
	}
	sent, err := c.SendGroupMsg(napcat.SendGroupMsgJSONBody{
		GroupID: new(group),
		Message: segs,
	})
	if err != nil {
		t.Fatalf("SendGroupMsg: %v", err)
	}
	if sent.MessageID == 0 {
		t.Error("SendGroupMsg 没拿到 message_id")
	}

	if got, err := c.GetMsg(napcat.GetMsgJSONBody{MessageID: getMsgID(t, sent.MessageID)}); err != nil {
		t.Errorf("GetMsg: %v", err)
	} else {
		t.Logf("GetMsg: real_id=%d time=%d 消息段 %d 个", got.RealID, got.Time, len(got.Message))
		if len(got.Message) == 0 {
			t.Log("提示: get_msg 没返回消息段, 检查实例是否开启了消息缓存")
		}
	}

	if _, err := c.DeleteMsg(napcat.DeleteMsgJSONBody{MessageID: delMsgID(t, sent.MessageID)}); err != nil {
		t.Errorf("DeleteMsg: %v", err)
	}

	// 私聊
	priv, err := c.SendPrivateMsg(napcat.SendPrivateMsgJSONBody{
		UserID:  new(user),
		Message: messageText(t, "NothingBot 类型化客户端联调, 可忽略"),
	})
	if err != nil {
		t.Fatalf("SendPrivateMsg: %v", err)
	}
	if _, err := c.DeleteMsg(napcat.DeleteMsgJSONBody{MessageID: delMsgID(t, priv.MessageID)}); err != nil {
		t.Errorf("DeleteMsg(private): %v", err)
	}

	// 交互
	if _, err := c.FriendPoke(napcat.FriendPokeJSONBody{UserID: user}); err != nil {
		t.Errorf("FriendPoke: %v", err)
	}
	if _, err := c.GroupPoke(napcat.GroupPokeJSONBody{GroupID: new(group), UserID: user}); err != nil {
		t.Errorf("GroupPoke: %v", err)
	}
}

// TestLiveEndpointsCovered 保证联调覆盖了 EasyOnebot 真正在用的端点。
func TestLiveEndpointsCovered(t *testing.T) {
	m, err := napcat.EndpointToAction()
	if err != nil {
		t.Fatalf("读取端点表失败: %v", err)
	}
	required := []string{
		"SendGroupMsg", "SendPrivateMsg", "SendGroupForwardMsg", "SendPrivateForwardMsg",
		"DeleteMsg", "GetMsg", "GetForwardMsg", "GetLoginInfo", "GetGroupList",
		"GetFriendList", "GetGroupMemberInfo", "FriendPoke", "GroupPoke",
	}
	for _, op := range required {
		if _, ok := m[op]; !ok {
			t.Errorf("端点表里缺少 %s", op)
		}
	}
}
