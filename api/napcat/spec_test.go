package napcat

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// legacyEndpoints 是 EasyOnebot 手写层里保留、但 NapCat spec 中已经不存在的端点。
//
// 每次同步新 spec 后, 这个列表 SHOULD 只减不增; 新增项必须说明原因。
var legacyEndpoints = map[string]string{
	"delete_group_file_folder": "NapCat 已改名为 delete_group_folder",
	"fetch_mface_key":          "NapCat 未实现",
	".join_friend_emoji_chain": "NapCat 未实现 (Lagrange 的 emoji 链接口)",
	".join_group_emoji_chain":  "NapCat 未实现 (Lagrange 的 emoji 链接口)",
	"rename_group_file_folder": "NapCat 未实现",
	"send_group_bot_callback":  "NapCat 未实现",
	"set_group_anonymous":      "NapCat 未实现 (匿名相关能力已移除)",
	"set_group_anonymous_ban":  "NapCat 未实现 (匿名相关能力已移除)",
	"set_group_bot_status":     "NapCat 未实现",
	"set_group_reaction":       "NapCat 未实现",
	"upload_image":             "NapCat 只提供 upload_image_to_qun_album",
}

// TestHandwrittenEndpointsExistInSpec 保证手写 API 层的端点名都能在 spec 里找到,
// 否则说明端点名拼错、或 NapCat 已移除该端点。
func TestHandwrittenEndpointsExistInSpec(t *testing.T) {
	specEps, err := SpecEndpoints()
	if err != nil {
		t.Fatalf("读取 spec 失败: %v", err)
	}
	inSpec := make(map[string]bool, len(specEps))
	for _, e := range specEps {
		inSpec[e] = true
	}

	handwritten, err := HandwrittenEndpoints("..")
	if err != nil {
		t.Fatalf("扫描手写 API 失败: %v", err)
	}
	if len(handwritten) == 0 {
		t.Fatal("没有扫描到任何 NewReq 调用, 检查 endpointRe 或目录")
	}

	var unknown []string
	for ep, files := range handwritten {
		if inSpec[ep] {
			continue
		}
		if _, ok := legacyEndpoints[ep]; ok {
			continue
		}
		unknown = append(unknown, ep+" ("+files[0]+")")
	}
	slices.Sort(unknown)
	if len(unknown) > 0 {
		t.Errorf("有 %d 个端点既不在 spec 里也没登记为 legacy:\n  %v", len(unknown), unknown)
	}

	// legacy 列表不能过期: spec 里重新出现同名端点时就该从列表里删掉
	var stale []string
	for ep := range legacyEndpoints {
		if inSpec[ep] {
			stale = append(stale, ep)
		}
	}
	slices.Sort(stale)
	if len(stale) > 0 {
		t.Errorf("这些端点已回到 spec 里, 请从 legacyEndpoints 移除:\n  %v", len(stale))
	}
}

// TestSpecOperationIDs 校验 operationId 注入结果: 数量与端点一致且无重名,
// 这是生成代码里类型名与端点一一对应的前提。
func TestSpecOperationIDs(t *testing.T) {
	eps, err := SpecEndpoints()
	if err != nil {
		t.Fatalf("读取 spec 失败: %v", err)
	}
	oids, err := SpecOperationIDs()
	if err != nil {
		t.Fatalf("读取 operationId 失败: %v", err)
	}
	if len(eps) != len(oids) {
		t.Errorf("端点 %d 个但 operationId %d 个", len(eps), len(oids))
	}

	seen := map[string]bool{}
	for _, o := range oids {
		if o == "" {
			t.Error("存在空 operationId, 同步脚本没注入成功")
			continue
		}
		if seen[o] {
			t.Errorf("operationId 重复: %s", o)
		}
		seen[o] = true
	}
}

// TestGeneratedCodeIsFresh 校验生成代码覆盖了全部端点, 防止忘记 go generate。
//
// 这里不逐字比对 Go 名字: oapi-codegen 的名字规范化会处理 initialism
// (check_url_safely -> CheckURLSafely, get_fileset_id -> GetFilesetID),
// 规则细节不值得在测试里复刻。改为抽掉非字母数字后做大小写无关的签名比对。
func TestGeneratedCodeIsFresh(t *testing.T) {
	oids, err := SpecOperationIDs()
	if err != nil {
		t.Fatalf("读取 operationId 失败: %v", err)
	}
	src, err := os.ReadFile("openapi.gen.go")
	if err != nil {
		t.Fatalf("读取生成代码失败: %v", err)
	}

	// 生成文件里所有标识符的名字签名集合
	generated := map[string]bool{}
	for _, m := range identRe.FindAllString(string(src), -1) {
		generated[signature(m)] = true
	}
	if len(generated) == 0 {
		t.Fatal("生成代码里没解析出任何标识符, 检查 identRe")
	}

	var missing []string
	for _, o := range oids {
		if !generated[signature(o)] {
			missing = append(missing, o)
		}
	}
	if len(missing) > 0 {
		slices.Sort(missing)
		show := missing[:min(len(missing), 10)]
		t.Errorf("生成代码里找不到 %d 个 operationId, 请执行 go generate ./api/napcat:\n  %v",
			len(missing), show)
	}
}

// identRe 匹配 Go 标识符 (够用: 抓类型名与常量名)
var identRe = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]*`)

// signature 抽掉非字母数字并转小写, 用于忽略命名风格差异的比对
func signature(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteString(strings.ToLower(string(r)))
		}
	}
	return b.String()
}

// TestVersionMetaMatchesSpec 校验同步脚本写下的指纹与当前 spec 一致。
func TestVersionMetaMatchesSpec(t *testing.T) {
	version, _, metaPaths, err := SpecVersionMeta()
	if err != nil {
		t.Fatalf("读取 version.json 失败: %v", err)
	}
	eps, err := SpecEndpoints()
	if err != nil {
		t.Fatalf("读取 spec 失败: %v", err)
	}
	if metaPaths != len(eps) {
		t.Errorf("version.json 记 %d 个端点, 实际 %d 个, spec 可能被手改过", metaPaths, len(eps))
	}
	if version == "" {
		t.Error("version.json 缺少 version 字段")
	}
	if !strings.Contains(specFile, version) {
		t.Errorf("specFile (%s) 与 version.json (%s) 版本不一致", specFile, version)
	}
}
