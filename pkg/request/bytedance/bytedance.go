package bytedance

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/bitly/go-simplejson"
	"github.com/ddliu/go-httpclient"
)

type PlayInfo struct {
	ResponseMetadata struct {
		RequestID string `json:"RequestId"`
		Action    string `json:"Action"`
		Version   string `json:"Version"`
		Service   string `json:"Service"`
		Region    string `json:"Region"`
	} `json:"ResponseMetadata"`
	Result struct {
		EncryptKey string `json:"EncryptKey"`
		CipherText string `json:"CipherText"`
		Data       struct {
			Status         int     `json:"Status"`
			VideoID        string  `json:"VideoID"`
			CoverURL       string  `json:"CoverUrl"`
			Duration       float64 `json:"Duration"`
			MediaType      string  `json:"MediaType"`
			EnableAdaptive bool    `json:"EnableAdaptive"`
			PlayInfoList   []struct {
				Bitrate          int     `json:"Bitrate"`
				FileHash         string  `json:"FileHash"`
				Size             int     `json:"Size"`
				Height           int     `json:"Height"`
				Width            int     `json:"Width"`
				Format           string  `json:"Format"`
				Codec            string  `json:"Codec"`
				Logo             string  `json:"Logo"`
				Definition       string  `json:"Definition"`
				Quality          string  `json:"Quality"`
				Duration         float64 `json:"Duration"`
				EncryptionMethod string  `json:"EncryptionMethod"`
				PlayAuth         string  `json:"PlayAuth"`
				PlayAuthID       string  `json:"PlayAuthID"`
				MainPlayURL      string  `json:"MainPlayUrl"`
				BackupPlayURL    string  `json:"BackupPlayUrl"`
				URLExpire        int     `json:"UrlExpire"`
				FileID           string  `json:"FileID"`
				P2PVerifyURL     string  `json:"P2pVerifyURL"`
				PreloadInterval  int     `json:"PreloadInterval"`
				PreloadMaxStep   int     `json:"PreloadMaxStep"`
				PreloadMinStep   int     `json:"PreloadMinStep"`
				PreloadSize      int     `json:"PreloadSize"`
				MediaType        string  `json:"MediaType"`
				CheckInfo        string  `json:"CheckInfo"`
			} `json:"PlayInfoList"`
			TotalCount int `json:"TotalCount"`
		} `json:"Data"`
	} `json:"Result"`
}

func GetPlayInfo(playAuthToken string) (playInfo PlayInfo, err error) {
	data, err := base64.StdEncoding.DecodeString(playAuthToken)
	if err != nil {
		return playInfo, err
	}
	sj, err := simplejson.NewJson(data)
	if err != nil {
		return playInfo, err
	}
	token, err := sj.Get("GetPlayInfoToken").String()
	if err != nil {
		return playInfo, err
	}
	playInfoRequestUrl := fmt.Sprintf("https://vod.bytedanceapi.com/?%s&ssl=true", token)
	resp, err := httpclient.Get(playInfoRequestUrl)
	if err != nil {
		return playInfo, err
	}
	if data, err = resp.ReadAll(); err != nil {
		return playInfo, err
	}
	if err = json.Unmarshal(data, &playInfo); err != nil {
		return playInfo, err
	}
	return playInfo, nil
}
