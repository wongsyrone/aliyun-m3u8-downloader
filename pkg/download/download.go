package download

import (
	"github.com/google/uuid"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/request"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

func AliyunDownload(output, saveFilename string, chanSize int, videoId, playAuth string) error {
	// 随机字符串
	clientRand := uuid.NewString()
	sj, err := request.GetVodPlayerInfo(clientRand, playAuth, videoId)
	if err != nil {
		return err
	}
	//tool.PrintJson(sj)
	playInfoList, err := sj.Get("PlayInfoList").Get("PlayInfo").Array()
	if err != nil {
		return err
	}
	playInfo := sj.Get("PlayInfoList").Get("PlayInfo").GetIndex(len(playInfoList) - 1)
	tool.PrintJson(playInfo)
	if saveFilename == "" {
		saveFilename, _ = sj.Get("VideoBase").Get("Title").String()
	}
	serverRand, _ := playInfo.Get("Rand").String()
	plaintext, _ := playInfo.Get("Plaintext").String()
	playURL, _ := playInfo.Get("PlayURL").String()
	tool.PrintJson(playURL)
	key := tool.DecryptKey(clientRand, serverRand, plaintext)
	downloader, err := NewDownloader(playURL, WithOutput(output), WithAliKey(key), WithFilename(saveFilename))
	if err != nil {
		panic(err)
	}
	if err := downloader.Start(chanSize); err != nil {
		panic(err)
	}
	return nil
}
