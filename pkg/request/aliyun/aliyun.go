package aliyun

import (
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/ddliu/go-httpclient"
	"github.com/wongsyrone/aliyun-m3u8-downloader/pkg/tool"
)

func init() {
	httpclient.Defaults(httpclient.Map{
		"Accept":                 "application/json, text/plain, */*",
		"Accept-Encoding":        "gzip, deflate, br, zstd",
		"Accept-Language":        "en-US,en;q=0.9,zh-CN;q=0.8,zh-TW;q=0.7,zh;q=0.6",
		httpclient.OPT_USERAGENT: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		//httpclient.OPT_PROXY:      "http://127.0.0.1:8888",
		httpclient.OPT_UNSAFE_TLS: true,
	})
}

func GetVodPlayerInfo(rand, playAuth string, opts ...OptionFunc) (*simplejson.Json, error) {
	rand, _ = tool.EncryptRand([]byte(rand))
	return getVodPlayerInfo(rand, playAuth, opts...)
}

func getVodPlayerInfo(rand, playAuth string, opts ...OptionFunc) (*simplejson.Json, error) {
	playInfoRequestUrl, err := GetPlayInfoRequestUrl(rand, playAuth, opts...)
	if err != nil {
		return nil, err
	}
	resp, err := httpclient.Get(playInfoRequestUrl)
	if err != nil {
		return nil, fmt.Errorf("getVodPlayerInfo: http get url: %s, err: %w", playInfoRequestUrl, err)
	}
	data, err := resp.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("getVodPlayerInfo: read err: %w", err)
	}
	sj, err := simplejson.NewJson(data)
	if err != nil {
		return nil, fmt.Errorf("getVodPlayerInfo: json decode: %s, err: %w", data, err)
	}
	return sj, nil
}
