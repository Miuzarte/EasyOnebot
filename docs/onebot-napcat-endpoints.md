# NapCat OneBot 11 端点索引

- 来源: `NapNeko/NapCatDocs` 仓库 `src/api/4.18.19/openapi.json` (由 napcat-schema 从源码 TypeBox 定义自动生成)
- OpenAPI 版本: `3.1.0`, 接口总数: **175**, 组件 schema: 39
- EasyOnebot 覆盖: 见每行末尾标记, 汇总见 [onebot-impl-napcat-openapi.md](./onebot-impl-napcat-openapi.md)

标记说明: `已封装` = EasyOnebot `api/` 里已有 `NewReq("<name>")`; 空 = 尚未封装。

## 核心接口 (6)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `cancel_group_todo` | 取消群待办 |  |
| `complete_group_todo` | 完成群待办 |  |
| `friend_poke` | 发送戳一戳 | 已封装 |
| `group_poke` | 发送戳一戳 | 已封装 |
| `send_poke` | 发送戳一戳 |  |
| `set_group_todo` | 设置群待办 |  |

## 系统接口 (13)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `can_send_image` | 是否可以发送图片 | 已封装 |
| `can_send_record` | 是否可以发送语音 | 已封装 |
| `clean_cache` | 清理缓存 | 已封装 |
| `get_credentials` | 获取登录凭证 | 已封装 |
| `get_csrf_token` | 获取 CSRF Token | 已封装 |
| `get_doubt_friends_add_request` | 获取可疑好友申请 |  |
| `get_group_system_msg` | 获取群系统消息 |  |
| `get_login_info` | 获取登录号信息 | 已封装 |
| `get_status` | 获取运行状态 | 已封装 |
| `get_version_info` | 获取版本信息 | 已封装 |
| `nc_get_packet_status` | 获取Packet状态 |  |
| `set_doubt_friends_add_request` | 处理可疑好友申请 |  |
| `set_restart` | 重启服务 | 已封装 |

## 消息接口 (10)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `_mark_all_as_read` | 标记所有消息已读 |  |
| `delete_msg` | 撤回消息 | 已封装 |
| `forward_friend_single_msg` | 转发单条消息 |  |
| `forward_group_single_msg` | 转发单条消息 |  |
| `get_msg` | 获取消息 | 已封装 |
| `mark_group_msg_as_read` | 标记群聊已读 |  |
| `mark_msg_as_read` | 标记消息已读 (Go-CQHTTP) | 已封装 |
| `mark_private_msg_as_read` | 标记私聊已读 |  |
| `send_msg` | 发送消息 | 已封装 |
| `send_private_msg` | 发送私聊消息 | 已封装 |

## 群组接口 (25)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `_del_group_notice` | 删除群公告 | 已封装 |
| `_get_group_notice` | 获取群公告 | 已封装 |
| `delete_essence_msg` | 移出精华消息 | 已封装 |
| `get_essence_msg_list` | 获取群精华消息 | 已封装 |
| `get_group_detail_info` | 获取群详细信息 |  |
| `get_group_ignore_add_request` | 获取群被忽略的加群请求 |  |
| `get_group_ignored_notifies` | 获取群忽略通知 |  |
| `get_group_info` | 获取群信息 | 已封装 |
| `get_group_list` | 获取群列表 | 已封装 |
| `get_group_member_info` | 获取群成员信息 | 已封装 |
| `get_group_member_list` | 获取群成员列表 | 已封装 |
| `get_group_shut_list` | 获取群禁言列表 |  |
| `send_group_msg` | 发送群消息 | 已封装 |
| `set_essence_msg` | 设置精华消息 | 已封装 |
| `set_group_add_request` | 处理加群请求 | 已封装 |
| `set_group_admin` | 设置群管理员 | 已封装 |
| `set_group_ban` | 群组禁言 | 已封装 |
| `set_group_card` | 设置群名片 | 已封装 |
| `set_group_kick` | 群组踢人 | 已封装 |
| `set_group_leave` | 退出群组 | 已封装 |
| `set_group_member_invite_policy` | 设置群成员邀请策略 |  |
| `set_group_member_permissions` | 设置群成员功能权限 |  |
| `set_group_name` | 设置群名称 | 已封装 |
| `set_group_new_member_history_visibility` | 设置新成员历史消息可见性 |  |
| `set_group_whole_ban` | 全员禁言 | 已封装 |

## 用户接口 (6)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `get_cookies` | 获取 Cookies | 已封装 |
| `get_friend_list` | 获取好友列表 | 已封装 |
| `get_recent_contact` | 获取最近会话 |  |
| `send_like` | 点赞 | 已封装 |
| `set_friend_add_request` | 处理加好友请求 | 已封装 |
| `set_friend_remark` | 设置好友备注 |  |

## 文件接口 (5)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `get_file` | 获取文件 |  |
| `get_group_file_url` | 获取群文件URL | 已封装 |
| `get_image` | 获取图片 | 已封装 |
| `get_private_file_url` | 获取私聊文件URL | 已封装 |
| `get_record` | 获取语音 | 已封装 |

## 频道接口 (2)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `get_guild_list` | 获取频道列表 |  |
| `get_guild_service_profile` | 获取频道个人信息 |  |

## 系统扩展 (16)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `add_custom_face` | 添加自定义表情 |  |
| `bot_exit` | 退出登录 |  |
| `delete_custom_face` | 删除自定义表情 |  |
| `fetch_custom_face` | 获取自定义表情 | 已封装 |
| `fetch_custom_face_detail` | 获取自定义表情详情 |  |
| `get_collection_list` | 获取收藏列表 |  |
| `get_mini_app_ark` | 获取小程序 Ark |  |
| `get_rkey` | 获取扩展 RKey | 已封装 |
| `get_rkey_server` | 获取 RKey 服务器 |  |
| `get_robot_uin_range` | 获取机器人 UIN 范围 |  |
| `nc_get_rkey` | 获取 RKey |  |
| `nc_get_user_status` | 获取用户在线状态 |  |
| `send_packet` | 发送原始数据包 |  |
| `set_custom_face_desc` | 修改自定义表情描述 |  |
| `set_input_status` | 设置输入状态 |  |
| `set_online_status` | 设置在线状态 |  |

## 消息扩展 (9)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `ArkShareGroup` | 分享群 (Ark) |  |
| `ArkSharePeer` | 分享用户 (Ark) |  |
| `click_inline_keyboard_button` | 点击内联键盘按钮 |  |
| `fetch_emoji_like` | 获取表情点赞详情 |  |
| `fetch_ptt_text` | 获取语音转文字结果 |  |
| `get_emoji_likes` | 获取消息表情点赞列表 |  |
| `send_ark_share` | 分享用户 (Ark) |  |
| `send_group_ark_share` | 分享群 (Ark) |  |
| `set_msg_emoji_like` | 设置消息表情点赞 |  |

## 群组扩展 (15)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `cancel_group_album_media_like` | 取消点赞群相册媒体 |  |
| `del_group_album_media` | 删除群相册媒体 |  |
| `do_group_album_comment` | 发表群相册评论 |  |
| `get_group_album_media_list` | 获取群相册媒体列表 |  |
| `get_group_info_ex` | 获取群详细信息 (扩展) |  |
| `get_group_signed_list` | 获取群组今日打卡列表 |  |
| `get_qun_album_list` | 获取群相册列表 |  |
| `send_group_sign` | 群打卡 |  |
| `set_group_add_option` | 设置群加群选项 |  |
| `set_group_album_media_like` | 点赞群相册媒体 |  |
| `set_group_remark` | 设置群备注 |  |
| `set_group_robot_add_option` | 设置群机器人加群选项 |  |
| `set_group_search` | 设置群搜索选项 |  |
| `set_group_sign` | 群打卡 |  |
| `upload_image_to_qun_album` | 上传图片到群相册 |  |

## 用户扩展 (4)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `get_friends_with_category` | 获取带分组的好友列表 |  |
| `get_profile_like` | 获取资料点赞 |  |
| `get_unidirectional_friend_list` | 获取单向好友列表 |  |
| `set_diy_online_status` | 设置自定义在线状态 |  |

## 文件扩展 (17)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `cancel_online_file` | 取消在线文件 |  |
| `create_flash_task` | 创建闪传任务 |  |
| `download_fileset` | 下载文件集 |  |
| `get_fileset_id` | 获取文件集 ID |  |
| `get_fileset_info` | 获取文件集信息 |  |
| `get_flash_file_list` | 获取闪传文件列表 |  |
| `get_flash_file_url` | 获取闪传文件链接 |  |
| `get_online_file_msg` | 获取在线文件消息 |  |
| `get_share_link` | 获取文件分享链接 |  |
| `move_group_file` | 移动群文件 | 已封装 |
| `receive_online_file` | 接收在线文件 |  |
| `refuse_online_file` | 拒绝在线文件 |  |
| `rename_group_file` | 重命名群文件 |  |
| `send_flash_msg` | 发送闪传消息 |  |
| `send_online_file` | 发送在线文件 |  |
| `send_online_folder` | 发送在线文件夹 |  |
| `trans_group_file` | 传输群文件 |  |

## Go-CQHTTP (27)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `.handle_quick_operation` | 处理快速操作 |  |
| `_get_model_show` | 获取机型显示 |  |
| `_send_group_notice` | 发送群公告 | 已封装 |
| `_set_model_show` | 设置机型 |  |
| `check_url_safely` | 检查URL安全性 |  |
| `create_group_file_folder` | 创建群文件目录 | 已封装 |
| `delete_friend` | 删除好友 | 已封装 |
| `delete_group_file` | 删除群文件 | 已封装 |
| `delete_group_folder` | 删除群文件目录 |  |
| `download_file` | 下载文件 |  |
| `get_forward_msg` | 获取合并转发消息 | 已封装 |
| `get_friend_msg_history` | 获取好友历史消息 | 已封装 |
| `get_group_at_all_remain` | 获取群艾特全体剩余次数 |  |
| `get_group_file_system_info` | 获取群文件系统信息 |  |
| `get_group_files_by_folder` | 获取群文件夹文件列表 | 已封装 |
| `get_group_honor_info` | 获取群荣誉信息 | 已封装 |
| `get_group_msg_history` | 获取群历史消息 | 已封装 |
| `get_group_root_files` | 获取群根目录文件列表 | 已封装 |
| `get_online_clients` | 获取在线客户端 |  |
| `get_stranger_info` | 获取陌生人信息 | 已封装 |
| `send_forward_msg` | 发送合并转发消息 | 已封装 |
| `send_group_forward_msg` | 发送群合并转发消息 | 已封装 |
| `send_private_forward_msg` | 发送私聊合并转发消息 | 已封装 |
| `set_group_portrait` | 设置群头像 | 已封装 |
| `set_qq_profile` | 设置QQ资料 |  |
| `upload_group_file` | 上传群文件 | 已封装 |
| `upload_private_file` | 上传私聊文件 | 已封装 |

## 扩展接口 (12)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `.ocr_image` | 图片 OCR 识别 (内部) |  |
| `create_collection` | 创建收藏 |  |
| `delete_qzone_msg` | 删除QQ空间说说 |  |
| `get_ai_characters` | 获取AI角色列表 | 已封装 |
| `get_clientkey` | 获取ClientKey |  |
| `ocr_image` | 图片 OCR 识别 | 已封装 |
| `send_qzone_msg` | 发表QQ空间说说 |  |
| `set_group_kick_members` | 批量踢出群成员 |  |
| `set_group_special_title` | 设置专属头衔 | 已封装 |
| `set_qq_avatar` | 设置QQ头像 | 已封装 |
| `set_self_longnick` | 设置个性签名 |  |
| `translate_en2zh` | 英文单词翻译 |  |

## 流式接口 (2)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `download_file_stream` | 下载文件流 |  |
| `upload_file_stream` | 上传文件流 |  |

## 流式传输扩展 (4)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `clean_stream_temp_file` | 清理流式传输临时文件 |  |
| `download_file_image_stream` | 下载图片文件流 |  |
| `download_file_record_stream` | 下载语音文件流 |  |
| `test_download_stream` | 测试下载流 |  |

## AI 扩展 (2)

| 端点 | 说明 | EasyOnebot |
|---|---|---|
| `get_ai_record` | 获取 AI 语音 | 已封装 |
| `send_group_ai_record` | 发送群 AI 语音 | 已封装 |
