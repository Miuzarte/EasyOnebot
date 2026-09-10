package napcat

import (
	"embed"
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

// specFS 把 spec 编进二进制。
//
// 这里只能 embed 本包目录下的文件, 所以 spec/ 放在包内; 用 embed 而不是运行期读文件,
// 是为了让调用方在任何工作目录下都能拿到端点表 (单元测试、go run、别的项目都成立)。
//
//go:embed spec/openapi-*.json spec/version.json
var specFS embed.FS

type specInfo struct {
	Version string `json:"version"`
	Latest  string `json:"latest"`
	Paths   int    `json:"paths"`
}

// SpecVersionMeta 返回同步脚本写下的版本指纹 (spec/version.json)。
func SpecVersionMeta() (version, latest string, paths int, err error) {
	b, err := specFS.ReadFile("spec/version.json")
	if err != nil {
		return "", "", 0, err
	}
	var info specInfo
	if err = json.Unmarshal(b, &info); err != nil {
		return "", "", 0, err
	}
	return info.Version, info.Latest, info.Paths, nil
}

// specDoc 是 spec 里本包用得到的部分
type specDoc struct {
	Paths map[string]map[string]struct {
		OperationID string `json:"operationId"`
	} `json:"paths"`
}

func loadSpec() (specDoc, error) {
	var doc specDoc
	b, err := specFS.ReadFile(specFile)
	if err != nil {
		return doc, err
	}
	err = json.Unmarshal(b, &doc)
	return doc, err
}

// SpecEndpoints 返回 spec 里定义的全部端点名 (path 去掉前导 '/'), 已排序。
func SpecEndpoints() ([]string, error) {
	doc, err := loadSpec()
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(doc.Paths))
	for p := range doc.Paths {
		out = append(out, trimPath(p))
	}
	sort.Strings(out)
	return out, nil
}

// SpecOperationIDs 返回 spec 里全部 operationId, 已排序。
func SpecOperationIDs() ([]string, error) {
	doc, err := loadSpec()
	if err != nil {
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

// EndpointToAction 返回 operationId -> 端点名 (即 WS 请求里的 action 字段) 的映射。
//
// 生成类型用 operationId 命名, 实际调用要发端点名, 两者靠 spec 里的 path 对齐。
func EndpointToAction() (map[string]string, error) {
	doc, err := loadSpec()
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(doc.Paths))
	for p, ops := range doc.Paths {
		action := trimPath(p)
		for _, op := range ops {
			if op.OperationID == "" {
				continue
			}
			out[op.OperationID] = action
		}
	}
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

// HandwrittenEndpoints 扫描 EasyOnebot 手写的 api 包, 返回端点名 -> 出现的文件。
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
