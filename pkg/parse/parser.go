package parse

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/ddliu/go-httpclient"
	"github.com/wongsyrone/aliyun-m3u8-downloader/pkg/tool"
)

type LoadKeyFunc func(m3u8Url, keyUrl string) (string, error)

type Result struct {
	URL  *url.URL
	M3u8 *M3u8
}

func FromM3u8URL(m3u8Url string, loadKeyFunc LoadKeyFunc) (*Result, error) {
	u, err := url.Parse(m3u8Url)
	if err != nil {
		return nil, err
	}
	m3u8Url = u.String()
	resp, err := httpclient.Get(m3u8Url)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %w", err)
	}
	data, err := resp.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("read m3u8 URL failed: %w", err)
	}
	//noinspection GoUnhandledErrorResult
	m3u8, err := parse(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromM3u8URL(tool.ResolveURL(u, sf.URI), loadKeyFunc)
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
	}
	keys := make(map[string]string, len(m3u8.Keys))
	for idx, k := range m3u8.Keys {
		switch {
		case k.Method == "" || k.Method == CryptMethodNONE:
			continue
		case k.Method == CryptMethodAES:
			keyStr, ok := keys[k.URI]
			if !ok {
				// Request URL to extract decryption key
				keyUrl := tool.ResolveURL(u, k.URI)
				// 加载key
				keyStr, err = loadKeyFunc(m3u8Url, keyUrl)
				if err != nil {
					return nil, err
				}
				keys[k.URI] = keyStr
			}
			m3u8.Keys[idx].Key = keyStr
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", k.Method)
		}
	}
	return result, nil
}

func FromM3u8Content(url, m3u8Str string, loadKeyFunc LoadKeyFunc) (*Result, error) {
	reader := strings.NewReader(m3u8Str)
	//noinspection GoUnhandledErrorResult
	m3u8, err := parse(reader)
	if err != nil {
		return nil, err
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		M3u8: m3u8,
	}
	for idx, k := range m3u8.Keys {
		switch {
		case k.Method == "" || k.Method == CryptMethodNONE:
			continue
		case k.Method == CryptMethodAES:
			// Request URL to extract decryption key
			keyUrl := tool.ResolveURL(nil, k.URI)
			// 加载key
			keyStr, err := loadKeyFunc(url, keyUrl)
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
