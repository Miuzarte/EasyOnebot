package message

import (
	"net/url"
	"testing"
)

func TestParseCQCode(t *testing.T) {
	t.Log(ParseCqCode("[CQ:at,qq=982809597,name=@謬紗特]"))
}

func TestParseCQCodes1(t *testing.T) {
	const msg = `text1[CQ:at,qq=982809597,name=@謬紗特] txt2[CQ:at,qq=982809597,name=@謬紗特] &#91;&#91;&#91;&#91; tt3`
	for _, seg := range ParseCqCodes(msg) {
		t.Logf("%s", seg.ToString(false))
	}
}

func TestParseCQCodes2(t *testing.T) {
	const msg = `[CQ:forward,id=aEGUZVEr+8sX6ZQzrmg14PkJJGZnfNU5LE4fXkEnVlBoEmDBTzyyoY/e/VIWknU6]`
	for _, seg := range ParseCqCodes(msg) {
		t.Logf("%s", seg.ToString(false))
	}
}

func TestParseCQCodes3(t *testing.T) {
	const msg = `text but with  ::: a : lot : of :`
	t.Logf("%s", ParseCqCode(msg))
	for _, seg := range ParseCqCodes(msg) {
		t.Logf("%s", seg.ToString(false))
	}
}

func TestUrlParse(t *testing.T) {
	// 测试 URL 解析
	urls := []string{
		"https://example.com/path?query=123",
		"http://example.com/path?query=123",
		"example.com/path?query=123",
		"ftp://example.com/path?query=123",
		"file:///sdcard/Pictures/mon.jpeg",
		"/sdcard/Pictures/mon.jpeg",
		"file:/\\C:\\Pictures\\mon.jpeg",
		"\\C:\\Pictures\\mon.jpeg",
	}
	for _, u := range urls {
		if parsedURL, err := url.Parse(u); err != nil {
			t.Logf("Failed to parse URL: %s", err)
		} else {
			t.Logf("Parsed URL: %s", parsedURL.String())
		}
	}
}
