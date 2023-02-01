package parse

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
}

func FromURL(link string, key string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	resp, err := httpclient.Get(link)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	m3u8, err := parse(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(tool.ResolveURL(u, sf.URI), key)
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
		case k.AliyunVoDEncryption && k.Method == CryptMethodAES:
			m3u8.Keys[idx].Key = key
		case !k.AliyunVoDEncryption && k.Method == CryptMethodAES:
			// 已知 key 值，直接赋值
			if key != "" {
				m3u8.Keys[idx].Key = key
				m3u8.Keys[idx].IV = "" // TODO: 要不要重置为空
				continue
			}
			// Request URL to extract decryption key
			keyURL := tool.ResolveURL(u, k.URI)
			resp, err = httpclient.Get(keyURL)
			if err != nil {
				return nil, fmt.Errorf("extract key failed: %s", err.Error())
			}
			if keyStr, err := resp.ToString(); err != nil {
				return nil, err
			} else {
				// fmt.Println("decryption key: ", keyStr)
				m3u8.Keys[idx].Key = keyStr
			}
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", k.Method)
		}
	}
	return result, nil
}
