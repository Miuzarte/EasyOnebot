#!/usr/bin/env python3
"""同步 NapCat 的 OpenAPI spec 并注入 operationId。

上游: NapNeko/NapCatDocs 仓库 src/api/<版本>/openapi.json
      由 NapCat 仓库的 packages/napcat-schema 在每个 release 自动生成并提交。

本脚本做两件事:
  1. 拉取指定版本 (默认最新) 的 spec, 原样写入 spec/openapi-raw-<版本>.json
  2. 注入 operationId (端点名蛇形转大驼峰), 写入 spec/openapi-<版本>.json

第 2 步是必须的: NapCat 的 spec 里 175 个端点一个 operationId 都没有, 不注入的话
oapi-codegen 只能按 HTTP 方法 + path 猜名字 (PostSendGroupMsg 之类), 又长又不好用。

用法:
  python3 sync_spec.py                 # 同步最新版本
  python3 sync_spec.py 4.18.19         # 同步指定版本
  python3 sync_spec.py --check         # 只检查本地是否已是最新, 不写文件
"""

import json
import re
import sys
import urllib.request
from pathlib import Path

API_LIST = "https://api.github.com/repos/NapNeko/NapCatDocs/contents/src/api"
RAW = "https://raw.githubusercontent.com/NapNeko/NapCatDocs/main/src/api/{ver}/openapi.json"

HERE = Path(__file__).resolve().parent
SPEC_DIR = HERE / "spec"

# 上游有两个带前导点的端点 (/.ocr_image, /.handle_quick_operation), 是与公开端点
# 同名同实现的旧路由。点号在 Go 名字规范化时会被吃掉, 会生成重复类型, 所以显式改名。
ALIAS_ENDPOINTS = {
    "/.ocr_image": "OcrImageLegacy",
    "/.handle_quick_operation": "HandleQuickOperationLegacy",
}


def version_key(v: str):
    return [int(x) if x.isdigit() else 0 for x in re.split(r"[.\-+]", v)]


def fetch(url: str) -> bytes:
    req = urllib.request.Request(url, headers={"User-Agent": "EasyOnebot-sync-spec"})
    with urllib.request.urlopen(req, timeout=60) as resp:
        return resp.read()


def latest_version() -> str:
    data = json.loads(fetch(API_LIST))
    versions = [d["name"] for d in data if d["type"] == "dir"]
    if not versions:
        raise SystemExit("上游没有找到任何版本目录")
    return sorted(versions, key=version_key)[-1]


def to_operation_id(path: str) -> str:
    """path -> operationId

    端点名按 '_' 拆词后转大驼峰, 与 oapi-codegen 生成类型名时的规范化保持一致
    (/send_group_forward_msg -> SendGroupForwardMsg, 生成 SendGroupForwardMsgJSONBody),
    这样调用方可以用同一个字符串既做类型名前缀又做查表键。

    少数带前导点的旧路由在 ALIAS_ENDPOINTS 里显式改名。
    """
    if path in ALIAS_ENDPOINTS:
        return ALIAS_ENDPOINTS[path]
    s = path.lstrip("/")
    if not s:
        raise ValueError("空 path")
    return "".join(p[:1].upper() + p[1:] for p in s.split("_") if p)


def inject_operation_ids(spec: dict) -> int:
    n = 0
    seen: dict[str, str] = {}
    for path, ops in spec["paths"].items():
        oid = to_operation_id(path)
        if oid in seen:
            raise SystemExit(f"operationId 冲突: {oid} 同时来自 {seen[oid]} 和 {path}")
        seen[oid] = path
        for op in ops.values():
            if not isinstance(op, dict):
                continue
            if op.get("operationId"):
                continue
            op["operationId"] = oid
            n += 1
    return n


def go_name(s: str) -> str:
    parts = re.split(r"[^0-9A-Za-z]+", s)
    return "".join(p[:1].upper() + p[1:] for p in parts if p)


def schema_of(node: dict, prop: str):
    """取属性的 schema, 解析同文件内的 $ref, 便于比较两个字段是不是同一个东西"""
    p = node["properties"][prop]
    seen = 0
    while isinstance(p, dict) and "$ref" in p and seen < 10:
        ref = p["$ref"]
        if not ref.startswith("#/"):
            return p
        cur = node
        for seg in ref[2:].split("/"):
            seg = seg.replace("~1", "/").replace("~0", "~")
            if isinstance(cur, dict):
                cur = cur.get(seg)
            elif isinstance(cur, list):
                cur = cur[int(seg)]
            else:
                return p
        p = cur
        seen += 1
    return p


def dedupe_properties(spec: dict) -> list[tuple[str, str, str]]:
    """去掉同名字段的双写, 返回 [(位置, Go 名, 被丢弃的 json 名)]

    NapCat 的 spec 里有少量字段同时以下划线与驼峰两种形式出现 (例如
    category_id 与 categoryId), 生成的 Go 结构体会撞名编译不过。
    只有在两个字段的类型/描述完全一致 (确认是同一个字段) 时才丢弃非下划线形式,
    否则直接报错, 逼人工确认 —— 不静默丢字段。
    """
    dropped: list[tuple[str, str, str]] = []

    def walk(node, loc):
        if isinstance(node, dict):
            props = node.get("properties")
            if isinstance(props, dict) and props:
                groups: dict[str, list[str]] = {}
                for k in props:
                    groups.setdefault(go_name(k), []).append(k)
                for g, names in groups.items():
                    if len(names) < 2:
                        continue
                    snake = [n for n in names if "_" in n]
                    if len(snake) != 1:
                        raise SystemExit(f"{loc}: 字段 {g} 撞名且无法判定保留哪个: {names}")
                    keep = snake[0]
                    for other in names:
                        if other == keep:
                            continue
                        a, b = schema_of(node, keep), schema_of(node, other)
                        # description 允许不同 (旧版本兼容字段通常只差一句说明),
                        # 其余部分必须完全一致, 否则视为不同字段, 交人工确认。
                        a_cmp = {k: v for k, v in a.items() if k != "description"}
                        b_cmp = {k: v for k, v in b.items() if k != "description"}
                        if json.dumps(a_cmp, sort_keys=True) != json.dumps(b_cmp, sort_keys=True):
                            raise SystemExit(
                                f"{loc}: 字段 {keep} 与 {other} 撞名且 schema 不同, 需要人工处理\n"
                                f"  {keep}: {json.dumps(a, ensure_ascii=False)}\n"
                                f"  {other}: {json.dumps(b, ensure_ascii=False)}"
                            )
                        props.pop(other)
                        dropped.append((loc, g, other))
            for k, v in node.items():
                walk(v, f"{loc}/{k}")
        elif isinstance(node, list):
            for i, v in enumerate(node):
                walk(v, f"{loc}[{i}]")

    walk(spec, "")
    return dropped


def extract_response_models(spec: dict) -> list[str]:
    """把每个端点成功响应的 data schema 提成 components.schemas 里的命名模型

    oapi-codegen 在只生成 models 时不会为 paths 里内联的响应 schema 产出 Go 类型,
    响应体就成了 map[string]interface{}, 调用方无法严格解码。这里把 data 的内联
    schema 搬进 components (名字固定为 <OperationID>Data, 与上面 x-go-type-name 一致),
    再把原位改成 $ref —— 语义完全等价, 只是让生成器看得见。

    返回新建的组件名列表。
    """
    created: list[str] = []
    components = spec.setdefault("components", {}).setdefault("schemas", {})
    for path, ops in spec["paths"].items():
        for op in ops.values():
            if not isinstance(op, dict):
                continue
            oid = op.get("operationId")
            if not oid:
                continue
            for resp in (op.get("responses") or {}).values():
                if not isinstance(resp, dict):
                    continue
                for content in (resp.get("content") or {}).values():
                    schema = content.get("schema") if isinstance(content, dict) else None
                    if not isinstance(schema, dict):
                        continue
                    for part in schema.get("allOf", []):
                        props = part.get("properties") if isinstance(part, dict) else None
                        data = props.get("data") if isinstance(props, dict) else None
                        if not isinstance(data, dict) or "$ref" in data:
                            continue
                        name = to_go_name(oid) + "Data"
                        existing = components.get(name)
                        if existing is not None and json.dumps(
                            existing, sort_keys=True
                        ) != json.dumps(data, sort_keys=True):
                            raise SystemExit(
                                f"{path}: 组件 {name} 已存在且内容不同, 需要人工处理"
                            )
                        if existing is None:
                            body = dict(data)
                            body.pop("x-go-type-name", None)
                            components[name] = body
                            created.append(name)
                        data.clear()
                        data["$ref"] = f"#/components/schemas/{name}"
    return created


# 上游把 ID / 时间戳一律声明成 `type: number` (无 format), oapi-codegen 会映射成
# float32 —— QQ 号、消息 id、时间戳都远超 float32 的 24 位有效精度, 严格解码会出错。
# 这些字段实际都是整数, 用 x-go-type 钉成 int64。
#
# 口径: 只覆盖"整数语义"的字段。像 age / chunk_size / count 这类安全的小整数不动,
# 真正的浮点字段一个都不动。按字段名匹配 (同名在不同 schema 里语义一致)。
INT64_FIELDS = {
    "message_id",
    "real_id",
    "message_seq",
    "group_id",
    "user_id",
    "target_id",
    "self_id",
    "sender_id",
    "invitor_uin",
    "actor",
    "emoji_package_id",
    "time",
    "join_time",
    "last_sent_time",
    "title_expire_time",
    "shut_up_timestamp",
    "login_days",
    "expired_time",
    "created_at",
    "seq",
}


def apply_int64_overrides(spec: dict) -> int:
    """把 [INT64_FIELDS] 里的 number 字段钉成 int64, 返回改动处数"""
    n = 0

    def walk(node):
        nonlocal n
        if isinstance(node, dict):
            props = node.get("properties")
            if isinstance(props, dict):
                for k, v in props.items():
                    if not isinstance(v, dict):
                        continue
                    if k in INT64_FIELDS and v.get("type") == "number" and "x-go-type" not in v:
                        v["x-go-type"] = "int64"
                        n += 1
            for v in node.values():
                walk(v)
        elif isinstance(node, list):
            for v in node:
                walk(v)

    walk(spec)
    return n


# 上游个别字段的类型与 NapCat 实际返回不一致, 或过于笼统 (只有 description 没有结构),
# 严格解码会失败。这里按 "组件名.字段名" 覆盖成正确 schema, 覆盖项在 version.json 里留痕。
#
# 依据是实测: 对一个真实 NapCat 实例调 get_msg, message 返回消息段数组, sender 返回对象。
SCHEMA_OVERRIDES = {
    "GetMsgData.message": {
        "description": "消息内容 (消息段数组)",
        "type": "array",
        "items": {"$ref": "#/components/schemas/OB11MessageData"},
    },
    "GetMsgData.sender": {
        "description": "发送者",
        "type": "object",
        "properties": {
            "user_id": {"description": "发送者 QQ 号", "type": "number", "x-go-type": "int64"},
            "nickname": {"description": "昵称", "type": "string"},
            "card": {"description": "群名片", "type": "string"},
            "role": {"description": "群角色", "type": "string"},
            "level": {"description": "等级", "type": "string"},
            "title": {"description": "头衔", "type": "string"},
        },
    },
}


def apply_schema_overrides(spec: dict) -> list[str]:
    """把 [SCHEMA_OVERRIDES] 应用到组件属性上, 返回实际生效的键"""
    applied: list[str] = []
    components = spec.get("components", {}).get("schemas", {})
    for key, schema in SCHEMA_OVERRIDES.items():
        comp, _, prop = key.partition(".")
        node = components.get(comp)
        if not isinstance(node, dict):
            raise SystemExit(f"schema 覆盖 {key}: 组件 {comp} 不存在")
        props = node.get("properties")
        if not isinstance(props, dict) or prop not in props:
            raise SystemExit(f"schema 覆盖 {key}: 字段不存在")
        props[prop] = json.loads(json.dumps(schema))
        applied.append(key)
    return applied


def to_go_name(s: str) -> str:
    parts = re.split(r"[^0-9A-Za-z]+", s)
    return "".join(p[:1].upper() + p[1:] for p in parts if p)


def main() -> int:
    args = [a for a in sys.argv[1:] if not a.startswith("--")]
    check_only = "--check" in sys.argv

    version = args[0] if args else latest_version()
    if args and args[0] == "latest":
        version = latest_version()

    latest = latest_version()
    if check_only:
        have = sorted((p.name for p in SPEC_DIR.glob("openapi-*.json")), key=version_key)
        print(f"上游最新: {latest}")
        print(f"本地已有: {have[-1] if have else '(无)'}")
        return 0 if have and have[-1].endswith(f"{latest}.json") else 1

    raw = fetch(RAW.format(ver=version))
    raw_spec = json.loads(raw)
    spec = json.loads(raw)
    injected = inject_operation_ids(spec)
    dropped = dedupe_properties(spec)
    created = extract_response_models(spec)
    overrides = apply_int64_overrides(spec)
    schema_overrides = apply_schema_overrides(spec)

    SPEC_DIR.mkdir(parents=True, exist_ok=True)
    raw_path = SPEC_DIR / f"openapi-raw-{version}.json"
    path = SPEC_DIR / f"openapi-{version}.json"

    # 上游产物的关键指纹, 用于日后核对是否被静默改动
    paths = len(raw_spec.get("paths", {}))
    schemas = len(raw_spec.get("components", {}).get("schemas", {}))
    meta = {
        "version": version,
        "openapi": raw_spec.get("openapi"),
        "title": raw_spec.get("info", {}).get("title"),
        "paths": paths,
        "schemas": schemas,
        "operationIdsInjected": injected,
        "responseModelComponentsCreated": len(created),
        "int64FieldOverrides": overrides,
        "schemaOverrides": schema_overrides,
        "aliasedEndpoints": ALIAS_ENDPOINTS,
        "droppedDuplicateProperties": [
            {"location": loc, "goName": g, "droppedJSONName": name} for loc, g, name in dropped
        ],
        "source": RAW.format(ver=version),
        "latest": latest,
    }

    raw_path.write_bytes(raw)
    path.write_text(json.dumps(spec, ensure_ascii=False, indent=1) + "\n", encoding="utf-8")
    (SPEC_DIR / "version.json").write_text(
        json.dumps(meta, ensure_ascii=False, indent=2) + "\n", encoding="utf-8"
    )

    print(f"上游版本 {version} (最新 {latest})")
    print(f"  {raw_path.name}: 原样保存, {len(raw)} 字节")
    print(f"  {path.name}: 注入 {injected} 个 operationId, 提取 {len(created)} 个响应模型组件, int64 覆盖 {overrides} 处, {paths} 端点 / {schemas} schema")
    for loc, g, name in dropped:
        print(f"    去除重复属性 {name} (Go 名 {g}) @ {loc}")
    print(f"  version.json: 记录来源与指纹")
    return 0


if __name__ == "__main__":
    sys.exit(main())
