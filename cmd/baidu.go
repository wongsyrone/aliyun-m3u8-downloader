package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/ddliu/go-httpclient"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/download"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
	"github.com/spf13/cobra"
)

type Encrypted struct {
	VideoKeyID        string `json:"videoKeyId"`
	PlayerID          string `json:"playerId"`
	EncryptedVideoKey string `json:"encryptedVideoKey"`
}

// baidubceCmd represents the veterinary app command
var baidubceCmd = &cobra.Command{
	Use:   "baidu",
	Short: "百度智能云视频下载",
	Long: `百度智能云视频下载. 使用示例:
aliyun-m3u8-downloader baidu -u 视频地址 -t token -o=/data/example -f 文件名 --concurrency 1`,
	PreRun: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			tool.PanicParameter("url")
		}
		chanSize, _ := cmd.Flags().GetInt("chanSize")
		if chanSize <= 0 {
			panic("parameter 'chanSize' must be greater than 0")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		token, _ := cmd.Flags().GetString("token")
		filename := viper.GetString("filename")
		output := viper.GetString("output")
		concurrency := viper.GetInt("concurrency")
		keys := make(map[string]string)
		downloader, err := download.NewDownloader(
			download.WithUrl(url),
			download.WithOutput(output),
			download.WithFilename(filename),
			download.WithLoadKeyFunc(func(m3u8Url, keyUrl string) (string, error) {
				// curl https://drm.media.baidubce.com/v1/tokenVideoKey?videoKeyId=job-nh7q3f8d99ahwr83&playerId=pid-1-5-1&token=87ec6bca6068f0a995a0dfd3d4592789e97d69fc129dd9f85bb1e1cd98790143_7d2195c92f8842a586f4299a8244b1fa_1674071062
				keyUrl = fmt.Sprintf("%s&playerId=pid-1-5-1&token=%s", keyUrl, token)
				if key, ok := keys[keyUrl]; ok {
					return key, nil
				}
				resp, err := httpclient.Get(keyUrl)
				if err != nil {
					return "", fmt.Errorf("extract key failed: %w", err)
				}
				data, err := resp.ReadAll()
				if err != nil {
					return "", err
				}
				var encrypted Encrypted
				if err = json.Unmarshal(data, &encrypted); err != nil {
					return "", err
				}
				keyStr, err := tool.BaiduDecrypt(tool.BaiduKey, encrypted.EncryptedVideoKey)
				if err != nil {
					return "", err
				}
				keys[keyUrl] = keyStr
				return keyStr, nil
			}))
		if err != nil {
			log.Fatalln(err)
		}
		if err := downloader.Start(concurrency); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(baidubceCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	baidubceCmd.Flags().StringP("url", "u", "", "m3u8 地址")
	baidubceCmd.Flags().StringP("token", "t", "", "获取key token")
	baidubceCmd.Flags().StringP("output", "o", "", "下载保存位置")
	baidubceCmd.Flags().StringP("filename", "f", "", "保存文件名")
	baidubceCmd.Flags().IntP("chanSize", "c", 1, "下载并发数")
	_ = baidubceCmd.MarkFlagRequired("url")
	_ = baidubceCmd.MarkFlagRequired("token")
}
