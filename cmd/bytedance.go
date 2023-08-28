package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/download"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
	"github.com/spf13/cobra"
)

// bytedanceCmd represents the bytedance command
var bytedanceCmd = &cobra.Command{
	Use:   "bytedance",
	Short: "字节跳动，火山引擎视频云视频加密下载工具",
	Long: `字节跳动，火山引擎视频云视频加密下载工具. 使用示例:
aliyun-m3u8-downloader bytedance -p "PlayAuthToken" -o=/data/example --concurrency 1 -f 文件名`,
	PreRun: func(cmd *cobra.Command, args []string) {
		httpclient.Defaults(httpclient.Map{
			"Accept-Encoding":        "gzip",
			httpclient.OPT_USERAGENT: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		})
		if referer, _ := cmd.Flags().GetString("referer"); referer != "" {
			httpclient.Defaults(httpclient.Map{
				httpclient.OPT_REFERER: referer,
			})
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		playAuth, _ := cmd.Flags().GetString("playAuth")
		filename := viper.GetString("filename")
		output := viper.GetString("output")
		concurrency := viper.GetInt("concurrency")
		if playAuth == "" {
			tool.PanicParameter("playAuth")
		}
		if err := download.Bytedance(output, filename, concurrency, playAuth); err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(bytedanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// aliyunCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bytedanceCmd.Flags().StringP("playAuth", "p", "", "web播放认证信息")
	_ = bytedanceCmd.MarkFlagRequired("playAuth")
}
