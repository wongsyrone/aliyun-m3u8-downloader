package cmd

import (
	"fmt"
	"log"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/request/aliyun"

	"github.com/ddliu/go-httpclient"
	"github.com/spf13/cobra"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/download"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

// aliyunCmd represents the aliyun command
var aliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "阿里云私有m3u8加密下载工具",
	Long: `阿里云私有m3u8加密下载工具. 使用示例:
aliyun-m3u8-downloader aliyun -p "WebPlayAuth" -v 视频id -o=/data/example --chanSize 1 -f 文件名`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if referer, _ := cmd.Flags().GetString("referer"); referer != "" {
			httpclient.Defaults(httpclient.Map{
				httpclient.OPT_REFERER: referer,
			})
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		playAuth, _ := cmd.Flags().GetString("playAuth")
		videoId, _ := cmd.Flags().GetString("videoId")
		filename, _ := cmd.Flags().GetString("filename")
		output, _ := cmd.Flags().GetString("output")
		chanSize, _ := cmd.Flags().GetInt("chanSize")
		region, _ := cmd.Flags().GetString("region")
		if playAuth == "" {
			tool.PanicParameter("playAuth")
		}
		if videoId == "" {
			tool.PanicParameter("videoId")
		}
		if chanSize <= 0 {
			panic("parameter 'chanSize' must be greater than 0")
		}
		var opts []aliyun.OptionFunc
		if region != "" {
			opts = append(opts, aliyun.WithRegion(region))
		}
		if err := download.Aliyun(output, filename, chanSize, videoId, playAuth, opts...); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(aliyunCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aliyunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// aliyunCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	aliyunCmd.Flags().StringP("playAuth", "p", "", "web播放认证信息")
	aliyunCmd.Flags().StringP("videoId", "v", "", "视频id")
	aliyunCmd.Flags().StringP("referer", "r", "", "referer请求头")
	aliyunCmd.Flags().StringP("output", "o", "", "下载保存位置")
	aliyunCmd.Flags().StringP("filename", "f", "", "保存文件名")
	aliyunCmd.Flags().IntP("chanSize", "c", 1, "下载并发数")
	aliyunCmd.Flags().StringP("region", "g", "", "地区，区域，默认值：cn-shanghai，可选值有：cn-beijing/cn-hangzhou/cn-shanghai等")
	_ = aliyunCmd.MarkFlagRequired("videoId")
	_ = aliyunCmd.MarkFlagRequired("playAuth")
}
