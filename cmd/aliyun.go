package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/download"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/log"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/request/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

// aliyunCmd represents the aliyun command
var aliyunCmd = &cobra.Command{
	Use:   "aliyun",
	Short: "阿里云私有m3u8加密下载工具",
	Long: `阿里云私有m3u8加密下载工具. 使用示例:
aliyun-m3u8-downloader aliyun -p "WebPlayAuth" -v 视频id -o=/data/example --concurrency 1 -f 文件名`,
	Run: func(cmd *cobra.Command, args []string) {
		playAuth, _ := cmd.Flags().GetString("playAuth")
		filename := viper.GetString("filename")
		output := viper.GetString("output")
		concurrency := viper.GetInt("concurrency")
		if playAuth == "" {
			tool.PanicParameter("playAuth")
		}
		var opts []aliyun.OptionFunc
		videoId, _ := cmd.Flags().GetString("videoId")
		if videoId != "" {
			opts = append(opts, aliyun.WithVideoId(videoId))
		}
		region, _ := cmd.Flags().GetString("region")
		if region != "" {
			opts = append(opts, aliyun.WithRegion(region))
		}
		if err := download.Aliyun(output, filename, concurrency, playAuth, opts...); err != nil {
			log.Errorf("aliyun err: %v", err)
			return
		}
		klog.Info("Done!")
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
	aliyunCmd.Flags().StringP("region", "g", "", "地区，区域，默认值：cn-shanghai，可选值有：cn-beijing/cn-hangzhou/cn-shanghai等")
	_ = aliyunCmd.MarkFlagRequired("playAuth")
}
