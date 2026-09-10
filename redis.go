package EasyOnebot

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/Miuzarte/EasyOnebot/api"
	"github.com/Miuzarte/EasyOnebot/event"
	"github.com/Miuzarte/EasyOnebot/internal/utils"

	"github.com/redis/rueidis"
)

const MAX_MSG_LIST_LEN = 80

type redisClient struct {
	Client rueidis.Client
}

func (r *redisClient) set(ctx context.Context, key string, value any, ex time.Duration) error {
	v, err := anyToString(value)
	if err != nil {
		return err
	}
	return r.Client.Do(ctx, r.Client.B().Set().Key(key).Value(v).Ex(ex).Build()).Error()
}

func (r *redisClient) mSet(ctx context.Context, kvs map[string]any) error {
	if len(kvs) == 0 {
		return nil
	}
	cmd := r.Client.B().Mset().KeyValue()
	for k, v := range kvs {
		value, err := anyToString(v)
		if err != nil {
			return err
		}
		cmd.KeyValue(k, value)
	}
	return r.Client.Do(ctx, cmd.Build()).Error()
}

func (r *redisClient) get(ctx context.Context, key string) redisResult {
	data, err := r.Client.Do(ctx, r.Client.B().Get().Key(key).Build()).ToString()
	return redisResult{Data: data, Err: err}
}

func (r *redisClient) mGet(ctx context.Context, keys ...string) redisResultList {
	result := r.Client.Do(ctx, r.Client.B().Mget().Key(keys...).Build())
	data, err := result.AsStrSlice()
	return redisResultList{Data: data, Err: err}
}

func (r *redisClient) scanGet(ctx context.Context, match string) redisResult {
	var cursor uint64
	for {
		se, err := r.Client.Do(ctx, r.Client.B().Scan().Cursor(cursor).Match(match).Build()).AsScanEntry()
		if err != nil {
			return redisResult{Err: err}
		}
		cursor = se.Cursor
		for _, key := range se.Elements {
			return r.get(ctx, key)
		}
		if cursor == 0 {
			break
		}
	}
	return redisResult{}
}

func (r *redisClient) lPush(ctx context.Context, key string, value any) error {
	var elems []string
	if value == nil {
		return nil
	}
	vValue := reflect.ValueOf(value)
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		l := vValue.Len()
		if l == 0 {
			return nil
		}
		elems = make([]string, l)
		for i := range l {
			s, err := anyToString(vValue.Index(i).Interface())
			if err != nil {
				return err
			}
			elems[i] = s
		}
	default:
		s, err := anyToString(value)
		if err != nil {
			return err
		}
		elems = []string{s}
	}
	return r.Client.Do(ctx, r.Client.B().Lpush().Key(key).Element(elems...).Build()).Error()
}

func (r *redisClient) lTrim(ctx context.Context, key string, start, stop int64) error {
	return r.Client.Do(ctx, r.Client.B().Ltrim().Key(key).Start(start).Stop(stop).Build()).Error()
}

func (r *redisClient) lLen(ctx context.Context, key string) (int64, error) {
	return r.Client.Do(ctx, r.Client.B().Llen().Key(key).Build()).AsInt64()
}

func (r *redisClient) lRange(ctx context.Context, key string, start, stop int64) redisResultList {
	result := r.Client.Do(ctx, r.Client.B().Lrange().Key(key).Start(start).Stop(stop).Build())
	data, err := result.AsStrSlice()
	return redisResultList{Data: data, Err: err}
}

type redisResult struct {
	Data string
	Err  error
}

func (rr *redisResult) IsNil() bool {
	return rr == nil || rueidis.IsRedisNil(rr.Err) || rr.Data == ""
}

func unmarshalResult[T any](rr redisResult) (to *T, err error) {
	if rr.Err != nil {
		return nil, rr.Err
	}
	to = new(T)
	return to, json.Unmarshal([]byte(rr.Data), to)
}

func (rr redisResult) MessageGroup() (*event.MessageGroup, error) {
	return unmarshalResult[event.MessageGroup](rr)
}

func (rr redisResult) MessagePrivate() (*event.MessagePrivate, error) {
	return unmarshalResult[event.MessagePrivate](rr)
}

func (rr redisResult) GroupInfo() (*api.GroupInfo, error) {
	return unmarshalResult[api.GroupInfo](rr)
}

func (rr redisResult) StrangerInfo() (*api.StrangerInfo, error) {
	return unmarshalResult[api.StrangerInfo](rr)
}

func (rr redisResult) GroupCardName() (string, error) {
	return rr.Data, rr.Err
}

type redisResultList struct {
	Data []string
	Err  error
}

func (rr *redisResultList) IsNil() bool {
	return rr == nil || rueidis.IsRedisNil(rr.Err) || len(rr.Data) == 0
}

func unmarshalResultList[T any](results redisResultList, to *[]T) (alsoTo *[]T, err error) {
	if len(results.Data) == 0 {
		return nil, nil
	}
	if results.Err != nil {
		return nil, results.Err
	}
	tos := make([]T, len(results.Data))
	for i, data := range results.Data {
		err := json.Unmarshal([]byte(data), &tos[i])
		if err != nil {
			return nil, err
		}
	}
	*to = tos
	return
}

func (rrl redisResultList) msgKeys() (msgKeys []string, err error) {
	if rrl.IsNil() {
		return nil, nil
	}
	if rrl.Err != nil {
		return nil, rrl.Err
	}
	return utils.DelElements(rrl.Data, ""), nil
}

func (rrl redisResultList) messageGroups() (mgl []event.MessageGroup, err error) {
	if rrl.Err != nil {
		return nil, rrl.Err
	}
	rrl.Data = utils.DelElements(rrl.Data, "")
	_, err = unmarshalResultList(rrl, &mgl)
	return
}

// --------------------------------

func anyToString(value any) (string, error) {
	marshal := func() (string, error) {
		j, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(j), nil
	}

	if value == nil { // reflect.ValueOf(nil) 是无效值, 取 Kind 会 panic
		return "", nil
	}
	vValue := reflect.ValueOf(value)
	switch vValue.Kind() {
	case reflect.String:
		return vValue.String(), nil

	case reflect.Slice:
		if vValue.Type().Elem().Kind() == reflect.Uint8 {
			return string(vValue.Bytes()), nil
		}
		return marshal()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// [reflect.Value.CanInt]
		return strconv.FormatInt(vValue.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		// [reflect.Value.CanUint]
		return strconv.FormatUint(vValue.Uint(), 10), nil

	case reflect.Complex64, reflect.Complex128:
		return strconv.FormatComplex(vValue.Complex(), 'f', -1, 64), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(vValue.Float(), 'f', -1, 64), nil

	default:
		return marshal()
	}
}

// --------------------------------

func (b *Bot) redisSetMessageGroup(msg *event.MessageGroup) error {
	key := "message_group:" + itoa(msg.GroupId) + ":" + itoa(msg.MessageId)
	return b.redis.set(b.ctx, key, msg, b.settings.redisExpireMsg)
}

func (b *Bot) RedisGetMessageGroup(groupId, msgId int) (*event.MessageGroup, error) {
	key := "message_group:" + itoa(groupId) + ":" + itoa(msgId)
	return b.redis.get(b.ctx, key).MessageGroup()
}

func (b *Bot) RedisScanMessageGroup(msgId int) (*event.MessageGroup, error) {
	match := "message_group:*:" + itoa(msgId)
	return b.redis.scanGet(b.ctx, match).MessageGroup()
}

// --------------------------------

func (b *Bot) redisSetMessagePrivate(msg *event.MessagePrivate) error {
	key := "message_private:" + itoa(msg.UserId) + ":" + itoa(msg.MessageId)
	return b.redis.set(b.ctx, key, msg, b.settings.redisExpireMsg)
}

func (b *Bot) RedisGetMessagePrivate(userId, msgId int) (*event.MessagePrivate, error) {
	key := "message_private:" + itoa(userId) + ":" + itoa(msgId)
	return b.redis.get(b.ctx, key).MessagePrivate()
}

func (b *Bot) RedisScanMessagePrivate(msgId int) (*event.MessagePrivate, error) {
	match := "message_private:*:" + itoa(msgId)
	return b.redis.scanGet(b.ctx, match).MessagePrivate()
}

// --------------------------------

func (b *Bot) redisSetGroupInfo(info *api.GroupInfo) error {
	key := "group_info:" + itoa(info.GroupId)
	return b.redis.set(b.ctx, key, info, b.settings.redisExpireInfo)
}

func (b *Bot) RedisGetGroupInfo(groupId int) (*api.GroupInfo, error) {
	key := "group_info:" + itoa(groupId)
	return b.redis.get(b.ctx, key).GroupInfo()
}

// --------------------------------

func (b *Bot) redisSetStrangerInfo(info *api.StrangerInfo) error {
	key := "stranger_info:" + itoa(info.UserId)
	return b.redis.set(b.ctx, key, info, b.settings.redisExpireInfo)
}

func (b *Bot) RedisGetStrangerInfo(userId int) (*api.StrangerInfo, error) {
	key := "stranger_info:" + itoa(userId)
	return b.redis.get(b.ctx, key).StrangerInfo()
}

// --------------------------------

func (b *Bot) redisSetGroupCardName(groupId int, userId int, name string) error {
	key := "group_card_name:" + itoa(groupId) + ":" + itoa(userId)
	return b.redis.set(b.ctx, key, name, b.settings.redisExpireInfo)
}

func (b *Bot) RedisGetGroupCardName(groupId int, userId int) (string, error) {
	key := "group_card_name:" + itoa(groupId) + ":" + itoa(userId)
	return b.redis.get(b.ctx, key).GroupCardName()
}

// --------------------------------

// 群内某个群成员 userId 带 at 的消息 keys,
// 具体 at 了谁需要获取并解析对应消息内容
func (b *Bot) redisPushAt(groupId, userId, messageId int) error {
	key := "at_to:" + itoa(groupId) + ":" + itoa(userId)
	err := b.redis.lTrim(b.ctx, key, 0, MAX_MSG_LIST_LEN-2)
	if err != nil {
		return err
	}
	valueK := "message_group:" + itoa(groupId) + ":" + itoa(messageId)
	return b.redis.lPush(b.ctx, key, valueK)
}

func (b *Bot) RedisRangeAt(groupId, userId int) ([]string, error) {
	key := "at_to:" + itoa(groupId) + ":" + itoa(userId)
	return b.redis.lRange(b.ctx, key, 0, MAX_MSG_LIST_LEN-1).msgKeys()
}

func (b *Bot) RedisRangeGetAt(groupId, userId int) (mgs []event.MessageGroup, err error) {
	msgIds, err := b.RedisRangeAt(groupId, userId)
	if err != nil {
		return nil, err
	}
	return b.redis.mGet(b.ctx, msgIds...).messageGroups()
}

// --------------------------------

// 群内某个群成员被 at 的消息 keys,
// 具体被谁 at 需要获取并解析对应消息发送者
func (b *Bot) redisPushAtBy(groupId, atBy, messageId int) error {
	key := "at_by:" + itoa(groupId) + ":" + itoa(atBy)
	err := b.redis.lTrim(b.ctx, key, 0, MAX_MSG_LIST_LEN-2)
	if err != nil {
		return err
	}
	valueK := "message_group:" + itoa(groupId) + ":" + itoa(messageId)
	return b.redis.lPush(b.ctx, key, valueK)
}

func (b *Bot) RedisRangeAtBy(groupId, userId int) ([]string, error) {
	key := "at_by:" + itoa(groupId) + ":" + itoa(userId)
	return b.redis.lRange(b.ctx, key, 0, MAX_MSG_LIST_LEN-1).msgKeys()
}

func (b *Bot) RedisRangeGetAtBy(groupId, userId int) (mgs []event.MessageGroup, err error) {
	msgIds, err := b.RedisRangeAtBy(groupId, userId)
	if err != nil {
		return nil, err
	}
	return b.redis.mGet(b.ctx, msgIds...).messageGroups()
}

// --------------------------------
// 群内某个群成员撤回的消息ID

func (b *Bot) redisPushRecall(groupId, userId, messageId int) error {
	key := "recall:" + itoa(groupId) + ":" + itoa(userId)
	err := b.redis.lTrim(b.ctx, key, 0, MAX_MSG_LIST_LEN-2)
	if err != nil {
		return err
	}
	vKey := "message_group:" + itoa(groupId) + ":" + itoa(messageId)
	return b.redis.lPush(b.ctx, key, vKey)
}

func (b *Bot) RedisRangeRecall(groupId, userId int) ([]string, error) {
	key := "recall:" + itoa(groupId) + ":" + itoa(userId)
	return b.redis.lRange(b.ctx, key, 0, MAX_MSG_LIST_LEN-1).msgKeys()
}

func (b *Bot) RedisRangeGetRecall(groupId, userId int) (mgs []event.MessageGroup, err error) {
	msgIds, err := b.RedisRangeRecall(groupId, userId)
	if err != nil {
		return nil, err
	}
	return b.redis.mGet(b.ctx, msgIds...).messageGroups()
}

// --------------------------------

// GetMessagePrivateTryCache 尝试从缓存中获取私聊消息,
// 只有 userId 不为 0 的时候才会访问缓存
func (b *Bot) GetMessagePrivateTryCache(userId, messageId int) (*event.MessagePrivate, error) {
	if messageId == 0 {
		return nil, ErrInvalidId
	}
	if userId != 0 {
		if mp, _ := b.RedisGetMessagePrivate(userId, messageId); mp != nil {
			return mp, nil
		}
	}

	resp, err := b.Call().Std.GetMsg(messageId)
	if err != nil {
		return nil, err
	}
	mp := utils.AnyCopy[event.MessagePrivate](resp)
	go func() {
		err := b.redisSetMessagePrivate(mp)
		if err != nil {
			b.log.Error().Err(err).Msg("缓存私聊消息出错")
		}
	}()
	return mp, nil
}

// GetMessageGroupTryCache 尝试从缓存中获取群消息,
// 只有 groupId 不为 0 的时候才会访问缓存
func (b *Bot) GetMessageGroupTryCache(groupId, messageId int) (*event.MessageGroup, error) {
	if messageId == 0 {
		return nil, ErrInvalidId
	}
	if groupId != 0 {
		if mg, _ := b.RedisGetMessageGroup(groupId, messageId); mg != nil {
			return mg, nil
		}
	}

	resp, err := b.Call().Std.GetMsg(messageId)
	if err != nil {
		return nil, err
	}
	mg := utils.AnyCopy[event.MessageGroup](resp)
	go func() {
		err := b.redisSetMessageGroup(mg)
		if err != nil {
			b.log.Error().Err(err).Msg("缓存群消息出错")
		}
	}()
	return mg, nil
}

// GetStrangerInfoTryCache 尝试从缓存中获取陌生人信息
func (b *Bot) GetStrangerInfoTryCache(userId int) (*api.StrangerInfo, error) {
	if userId == 0 {
		return nil, ErrInvalidId
	}
	if si, _ := b.RedisGetStrangerInfo(userId); si != nil {
		return si, nil
	}

	resp, err := b.Call().Std.GetStrangerInfo(userId, false)
	if err != nil {
		return nil, err
	}
	go func() {
		err := b.redisSetStrangerInfo(resp)
		if err != nil {
			b.log.Error().Err(err).Msg("缓存陌生人信息出错")
		}
	}()
	return resp, nil
}

// GetGroupInfoTryCache 尝试从缓存中获取群信息
func (b *Bot) GetGroupInfoTryCache(groupId int) (*api.GroupInfo, error) {
	if groupId == 0 {
		return nil, ErrInvalidId
	}
	if gi, _ := b.RedisGetGroupInfo(groupId); gi != nil {
		return gi, nil
	}

	resp, err := b.Call().Std.GetGroupInfo(groupId, false)
	if err != nil {
		return nil, err
	}
	go func() {
		err := b.redisSetGroupInfo(resp)
		if err != nil {
			b.log.Error().Err(err).Msg("缓存群信息出错")
		}
	}()
	return resp, nil
}

// GetGroupCardNameTryCache 尝试从缓存中获取群名片, 不存在则返回昵称
func (b *Bot) GetGroupCardNameTryCache(groupId, userId int) (string, error) {
	if userId == 0 {
		return "", ErrInvalidId
	}
	if groupId != 0 {
		if cardName, _ := b.RedisGetGroupCardName(groupId, userId); cardName != "" {
			return cardName, nil
		}
	}

	si, err := b.GetStrangerInfoTryCache(userId)
	if err != nil {
		return "", err
	}
	go func() {
		err := b.redisSetGroupCardName(groupId, userId, si.Nickname)
		if err != nil {
			b.log.Error().Err(err).Msg("缓存群名片出错")
		}
	}()
	return si.Nickname, nil
}
