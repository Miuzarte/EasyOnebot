package EasyOnebot

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/event"

	"github.com/go-viper/mapstructure/v2"
	// "github.com/gorilla/websocket"
	"github.com/coder/websocket"
)

// var wsDialer = &websocket.Dialer{
// 	Proxy:            nil, // no proxy
// 	HandshakeTimeout: 45 * time.Second,
// }

var dialOp = websocket.DialOptions{
	HTTPClient: &http.Client{
		Transport: &http.Transport{
			Proxy:                 nil, // disable proxy
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	},
}

// connection 连接
type connection struct {
	ws           *websocket.Conn // websocket.Conn
	url          string          // websocket url
	token        string          // access token
	mu           sync.Mutex      // conn mu
	dialCount    int             // 连接次数
	lastDialTime time.Time       // 上次连接时间
}

func (c *connection) dial(ctx context.Context, token string) (err error) {
	if token != "" {
		dialOp.HTTPHeader = make(http.Header, 1)
		dialOp.HTTPHeader.Set("Authorization", "Bearer "+token)
	}
	c.ws, _, err = websocket.Dial(ctx, c.url, &dialOp)
	// c.ws, _, err = wsDialer.Dial(c.url, reqHeader)
	if err != nil {
		return err
	}
	c.ws.SetReadLimit(-1)
	c.lastDialTime = time.Now()
	c.dialCount++
	return nil
}

// heartbeat 心跳
type heartbeat struct {
	signal   chan int // passing interval
	running  bool     // 监听
	interval int      // 间隔 ms
	count    int      // 计数
	lost     int      // 丢失计数
}

// 接收方准备 echo 对应通道并 put 入, 由通信 handler get 出 echo 对应通道并放入响应
type apiRespChanPool struct {
	timeout time.Duration                 // 超时时间
	m       map[string]chan *api.Response // 响应 echo: chan
	mu      sync.Mutex
}

func (a *apiRespChanPool) put(echo string, ch chan *api.Response) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.m[echo] = ch
}

func (a *apiRespChanPool) del(echo string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.m, echo)
}

func (a *apiRespChanPool) get(echo string) (chan *api.Response, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	ch, ok := a.m[echo]
	if ok {
		delete(a.m, echo)
	}
	return ch, ok
}

func (b *Bot) Run() *Bot {
	if b.conn.url == "" {
		b.log.Fatal().Msg("Bot.Run: empty websocket url")
	}

	b.ctx, b.cancel = context.WithCancel(context.Background())

	var tryCount int
	var err error
	b.log.Debug().Msgf("开始连接至 %s ...", b.conn.url)
	for {
		tryCount++
		err = b.conn.dial(b.ctx, b.conn.token)
		if err != nil {
			b.log.Error().Err(err).Msg("连接失败")
			<-time.After(time.Second * 5)
			continue
		}
		break
	}
	if tryCount == 1 {
		b.log.Info().Msgf("成功连接至 %s", b.conn.url)
	} else {
		b.log.Info().Msgf("成功连接至 %s #%d", b.conn.url, tryCount)
	}

	go b.listen()

	n, u, err := b.initSelfInfo()
	if err == nil {
		b.log.Info().Msgf("账号信息: %s (%d)", n, u)
	} else {
		b.log.Error().Err(err).Msg("初始化账号信息失败")
	}

	n, v, err := b.initImplementationInfo()
	if err == nil {
		b.log.Info().Msgf("OneBot 实现: %s (%s)", n, v)
	} else {
		b.log.Error().Err(err).Msg("初始化 OneBot 实现信息失败")
	}

	if b.settings.onlineNotify {
		go func() {
			var err error
			if tryCount == 1 {
				err = b.Log2Sus.Info("已上线")
			} else {
				err = b.Log2Sus.Infof("已上线 #%d", tryCount)
			}
			if err != nil {
				b.log.Error().Err(err).Msg("向超级用户发送消息失败")
			}
		}()
	}

	return b
}

func (b *Bot) listen() {
	var err error
	wait := make(chan struct{})

	go func() {
		defer func() {
			b.cancel()
			wait <- struct{}{}
		}()

		for {
			// select {
			// case <-b.ctx.Done():
			// 	return
			// default:
			// }

			var data []byte
			_, data, err = b.conn.ws.Read(b.ctx)
			// _, data, err = b.conn.ws.ReadMessage()
			if err != nil {
				return
			}
			b.statistics.receive++
			b.log.Trace().Msgf("receive: %s", data)

			m := map[string]any{}
			err := json.Unmarshal(data, &m)
			if err != nil {
				b.log.Panic().Msgf(
					"failed to unmarshal: %v\ndata: %s",
					err, string(data),
				)
				continue
			}

			if _, ok := m["echo"]; ok {
				apiResp := &api.Response{}
				err = mapstructure.Decode(m, apiResp)
				if err != nil {
					b.log.Panic().Msgf("failed to parse api response: %v\ndata: %s\nunmarshaled: %+v",
						err, data, m)
					continue
				}
				b.log.Trace().Msgf("decoded api response: %+v", apiResp)
				go b.handleApiResp(apiResp)

			} else {
				event := &event.Event{RawEvent: m}
				err = mapstructure.Decode(m, event)
				if err != nil {
					b.log.Panic().Msgf("failed to parse event: %v\ndata: %s\nunmarshaled: %+v",
						err, data, m)
					continue
				}
				b.log.Trace().Msgf("decoded event: %+v", event)
				go b.handleEvent(event)

			}
		}
	}()

	<-wait // 等待协程结束, 获取完错误
	if err == nil ||
		websocket.CloseStatus(err) == websocket.StatusNormalClosure {
		// websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		// 正常断开
		return
	}
	b.log.Error().Err(err).Msg("连接意外断开")
	// 重连
	go b.Run()
}

func (b *Bot) Stop() *Bot {
	if b.conn.ws == nil {
		return b
	}

	if b.settings.offlineNotify {
		err := b.Log2Sus.Info("已下线")
		if err != nil {
			b.log.Error().Err(err).Msg("向超级用户发送消息失败")
		}
	}

	b.conn.mu.Lock()
	defer b.conn.mu.Unlock()
	b.conn.ws.Close(websocket.StatusNormalClosure, "Byebye from EasyOnebot")
	// b.conn.ws.Close() // 手动关闭返回的错误无法判断
	// b.conn.ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	b.cancel()
	b.conn.ws = nil

	return b
}

// Statistics 统计信息
func (b *Bot) Statistics() (connDuration time.Duration, dialCount, hbCount, hbLost int) {
	return time.Since(b.conn.lastDialTime),
		b.conn.dialCount,
		b.hb.count,
		b.hb.lost
}

// initSelfInfo 初始化账号信息, 修改 log scope
func (b *Bot) initSelfInfo() (nickname string, userId int, err error) {
	resp, err := b.Call().Std.GetLoginInfo()
	if err != nil {
		return "", 0, err
	}
	b.loginInfo = resp
	b.selfId = resp.UserId
	b.setScope(resp.Nickname)
	b.alias = append([]string{itoa(resp.UserId), strings.ToLower(strings.TrimSpace(resp.Nickname))}, b.alias...) // 用于识别假at
	return resp.Nickname, resp.UserId, nil
}

func (b *Bot) initImplementationInfo() (appName, appVersion string, err error) {
	resp, err := b.Call().Std.GetVersionInfo()
	if err != nil {
		return "", "", err
	}
	b.versionInfo = resp
	return resp.AppName, resp.AppVersion, nil
}

// post 发送数据
func (b *Bot) post(data []byte) error {
	if b.conn.ws == nil {
		return ErrNoWsConn
	}
	if len(data) == 0 {
		return ErrEmptyData
	}

	b.log.Trace().Msgf("posting: %s", data)
	b.conn.mu.Lock()
	defer b.conn.mu.Unlock()
	err := b.conn.ws.Write(b.ctx, websocket.MessageText, data)
	// err := b.conn.ws.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		b.log.Error().Err(err).Msg("failed to post data")
		return err
	}
	return nil
}

// heartbeatLoop 心跳监听
//
// 在接收到心跳包之前不会运行
func (b *Bot) heartbeatLoop() {
	defer func() {
		b.hb.running = false
	}()

	if b.hb.interval == 0 {
		b.hb.interval = 9000
	}

	b.log.Debug().Msg("heartbeatLoop started")
	for {
		select {
		case interval := <-b.hb.signal:
			b.hb.count++
			b.hb.interval = interval

		// 超时
		case <-time.After(time.Millisecond * time.Duration(b.hb.interval+1000)):
			b.hb.lost++
			b.log.Warn().Msgf("心跳超时#%d", b.hb.lost)

		// 结束
		case <-b.ctx.Done():
			return

		}
	}
}
