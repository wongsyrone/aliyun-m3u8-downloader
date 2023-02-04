package download

import (
	"errors"
	"fmt"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/request/bytedance"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

func Bytedance(output, saveFilename string, chanSize int, playAuthToken string) error {
	playInfoResp, err := bytedance.GetPlayInfo(playAuthToken)
	if err != nil {
		return err
	}
	playInfoList := playInfoResp.Result.Data.PlayInfoList
	n := len(playInfoList)
	if n == 0 {
		return errors.New("后去播放信息视频，playAuth 可能无效")
	}
	playInfo := playInfoList[n-1] // 高清
	if playInfo.Format != "hls" {
		return fmt.Errorf("不是支持 %s 格式视频下载", playInfo.Format)
	}
	tool.PrintJson(playInfo)
	if saveFilename == "" {
		saveFilename = playInfo.FileID + ".mp4"
	}
	key := tool.PlayAuthDecrypt(playInfo.PlayAuth) // 解密 PlayAuth 为 m3u8 视频解密key
	downloader, err := NewDownloader(playInfo.MainPlayURL, WithOutput(output), WithKey(key), WithFilename(saveFilename))
	if err != nil {
		panic(err)
	}
	if err = downloader.Start(chanSize); err != nil {
		panic(err)
	}
	return nil
}
