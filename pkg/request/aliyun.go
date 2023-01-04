package request

import (
	"log"

	"github.com/bitly/go-simplejson"
	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

func init() {
	httpclient.Defaults(httpclient.Map{
		"Accept":                 "application/json, text/plain, */*",
		"Accept-Encoding":        "gzip, deflate, br",
		"Accept-Language":        "zh-CN,zh;q=0.9,en;q=0.8",
		httpclient.OPT_USERAGENT: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		//httpclient.OPT_PROXY:      "http://127.0.0.1:8888",
		httpclient.OPT_UNSAFE_TLS: true,
	})
}

func GetVodPlayerInfo(rand, playAuth, fileId string, opts ...aliyun.OptionFunc) (*simplejson.Json, error) {
	rand, _ = tool.EncryptRand([]byte(rand))
	return getVodPlayerInfo(rand, playAuth, fileId, opts...)
}

func getVodPlayerInfo(rand, playAuth, fileId string, opts ...aliyun.OptionFunc) (*simplejson.Json, error) {
	playInfoRequestUrl, err := aliyun.GetPlayInfoRequestUrl(rand, playAuth, fileId, opts...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	resp, err := httpclient.Get(playInfoRequestUrl)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	data, err := resp.ReadAll()
	if err != nil {
		log.Println(resp)
		return nil, err
	}
	return simplejson.NewJson(data)
}
