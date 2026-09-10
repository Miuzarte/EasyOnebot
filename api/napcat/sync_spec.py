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

    直接取整条 path 去掉 '/' 后首字母大写, 不按 '_' 拆词:
    NapCat 的 path 本身就是端点名 (/send_group_forward_msg -> SendGroupForwardMsg)。
    少数带前导点的旧路由在 ALIAS_ENDPOINTS 里显式改名。
    """
    if path in ALIAS_ENDPOINTS:
        return ALIAS_ENDPOINTS[path]
    s = path.lstrip("/")
    if not s:
        raise ValueError("空 path")
    return s[:1].upper() + s[1:]


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


def name_response_data(spec: dict) -> list[str]:
    """给每个端点的响应 data 内联 schema 起名字, 返回命名的类型列表

    NapCat 的响应形如 allOf[BaseResponse, {data: <匿名对象>}], 匿名对象在只生成 models 时
    不会产出 Go 类型, 调用方就没法强类型解析 data。用 oapi-codegen 支持的
    x-go-type-name 扩展直接命名, 生成 <OperationID>Data。
    """
    named: list[str] = []
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
                        if not isinstance(data, dict):
                            continue
                        # data 是 $ref 时不用命名 (已经有名字了)
                        if "$ref" in data:
                            continue
                        name = to_go_name(oid) + "Data"
                        data["x-go-type-name"] = name
                        named.append(name)
    return named


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
    named = name_response_data(spec)

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
        "responseDataTypesNamed": len(named),
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
    print(f"  {path.name}: 注入 {injected} 个 operationId, 命名 {len(named)} 个响应 data 类型, {paths} 端点 / {schemas} schema")
    for loc, g, name in dropped:
        print(f"    去除重复属性 {name} (Go 名 {g}) @ {loc}")
    print(f"  version.json: 记录来源与指纹")
    return 0


if __name__ == "__main__":
    sys.exit(main())
