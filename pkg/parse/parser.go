package parse

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

type LoadKeyFunc func(m3u8Url, keyUrl string) (string, error)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
}

func FromURL(m3u8Url string, loadKeyFunc LoadKeyFunc) (*Result, error) {
	u, err := url.Parse(m3u8Url)
	if err != nil {
		return nil, err
	}
	m3u8Url = u.String()
	resp, err := httpclient.Get(m3u8Url)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %w", err)
	}
	//noinspection GoUnhandledErrorResult
	m3u8, err := parse(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(tool.ResolveURL(u, sf.URI), loadKeyFunc)
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
	}
	for idx, k := range m3u8.Keys {
		switch {
		case k.Method == "" || k.Method == CryptMethodNONE:
			continue
		case k.Method == CryptMethodAES:
			// Request URL to extract decryption key
			keyUrl := tool.ResolveURL(u, k.URI)
			// 加载key
			keyStr, err := loadKeyFunc(m3u8Url, keyUrl)
			if err != nil {
				return nil, err
			}
			m3u8.Keys[idx].Key = keyStr
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", k.Method)
		}
	}
	return result, nil
}
