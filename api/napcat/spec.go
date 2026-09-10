package napcat

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// specFile 指向当前生效的、已注入 operationId 的 spec。
//
// 同步新版本 spec 后需要同步改这里 (以及 generate.go 里的 go:generate 指令)。
const specFile = "spec/openapi-4.18.19.json"

type specInfo struct {
	Version string `json:"version"`
	Latest  string `json:"latest"`
	Paths   int    `json:"paths"`
}

// readVersionMeta 读取同步脚本写下的版本指纹 (spec/version.json)。
func readVersionMeta() (specInfo, error) {
	var info specInfo
	b, err := os.ReadFile(filepath.Join(specDir(), "version.json"))
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(b, &info)
	return info, err
}

func specDir() string { return "spec" }

// SpecFile 返回当前生效的 spec 路径 (相对本包目录)。
func SpecFile() string { return specFile }

// SpecEndpoints 返回 spec 里定义的全部端点名 (path 去掉前导 '/')。
//
// 这里直接解析 spec JSON, 不依赖生成代码, 因此可以独立作为"端点名是否存在"的判据。
func SpecEndpoints() ([]string, error) {
	b, err := os.ReadFile(specFile)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(doc.Paths))
	for p := range doc.Paths {
		out = append(out, trimPath(p))
	}
	sort.Strings(out)
	return out, nil
}

// SpecOperationIDs 返回 spec 里全部 operationId, 顺序与 [SpecEndpoints] 一致。
func SpecOperationIDs() ([]string, error) {
	b, err := os.ReadFile(specFile)
	if err != nil {
		return nil, err
	}
	var doc struct {
		Paths map[string]map[string]struct {
			OperationID string `json:"operationId"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	out := make([]string, 0, len(doc.Paths))
	for _, ops := range doc.Paths {
		for _, op := range ops {
			out = append(out, op.OperationID)
		}
	}
	sort.Strings(out)
	return out, nil
}

func trimPath(p string) string {
	for len(p) > 0 && p[0] == '/' {
		p = p[1:]
	}
	return p
}

// endpointRe 匹配手写 API 层里的 NewReq("<endpoint>", ...) 调用。
var endpointRe = regexp.MustCompile(`NewReq\(\s*"([a-zA-Z0-9_.]+)"`)

// HandwrittenEndpoints 扫描 EasyOnebot 手写的 api 包, 返回其中出现的全部端点名。
func HandwrittenEndpoints(apiDir string) (map[string][]string, error) {
	entries, err := os.ReadDir(apiDir)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".go" {
			continue
		}
		b, err := os.ReadFile(filepath.Join(apiDir, e.Name()))
		if err != nil {
			return nil, err
		}
		for _, m := range endpointRe.FindAllStringSubmatch(string(b), -1) {
			out[m[1]] = append(out[m[1]], e.Name())
		}
	}
	return out, nil
}
