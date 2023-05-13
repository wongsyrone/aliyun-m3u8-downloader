package download

import (
	"fmt"

	"github.com/google/uuid"

	parseAliyun "github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/request/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

const (
	AliyunVoDEncryption = "AliyunVoDEncryption"
)

func Aliyun(output, saveFilename string, chanSize int, videoId, playAuth string, parseOpts ...parseAliyun.OptionFunc) error {
	// 随机字符串
	clientRand := uuid.NewString()
	sj, err := aliyun.GetVodPlayerInfo(clientRand, playAuth, videoId, parseOpts...)
	if err != nil {
		return err
	}
	//tool.PrintJson(sj)
	playInfoList, err := sj.Get("PlayInfoList").Get("PlayInfo").Array()
	if err != nil {
		return fmt.Errorf("donwload: get PlayInfo err: %w", err)
	}
	playInfo := sj.Get("PlayInfoList").Get("PlayInfo").GetIndex(len(playInfoList) - 1)
	tool.PrintJson(playInfo)
	if saveFilename == "" {
		saveFilename, _ = sj.Get("VideoBase").Get("Title").String()
	}
	encryptType, _ := playInfo.Get("EncryptType").String()
	playURL, _ := playInfo.Get("PlayURL").String()
	format, _ := playInfo.Get("Format").String()
	tool.PrintJson(playURL)
	opts := []DownloaderOption{WithUrl(playURL), WithOutput(output), WithFilename(saveFilename)}
	if encryptType == AliyunVoDEncryption {
		serverRand, _ := playInfo.Get("Rand").String()
		plaintext, _ := playInfo.Get("Plaintext").String()
		key := tool.DecryptKey(clientRand, serverRand, plaintext)
		opts = append(opts, WithKey(key))
	}
	if format == "mp4" {
		opts = append(opts, WithMp4(true))
	}

	downloader, err := NewDownloader(opts...)
	if err != nil {
		return err
	}
	if err = downloader.Start(chanSize); err != nil {
		return err
	}
	return nil
}
