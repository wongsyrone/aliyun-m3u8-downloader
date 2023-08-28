package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/download"
	"github.com/spf13/cobra"
)

// normalCmd represents the normal command
var normalCmd = &cobra.Command{
	Use:   "normal",
	Short: "普通m3u8 或 标准AES-128加密 下载",
	Long: `普通m3u8 或 标准AES-128加密 下载. 使用示例:
aliyun-m3u8-downloader normal -u=https://www.lbbniu.com/index.m3u8 -o=/data/example --concurrency 1`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		filename := viper.GetString("filename")
		output := viper.GetString("output")
		concurrency := viper.GetInt("concurrency")
		downloader, err := download.NewDownloader(download.WithUrl(url), download.WithFilename(filename), download.WithOutput(output))
		if err != nil {
			panic(err)
		}
		if err := downloader.Start(concurrency); err != nil {
			panic(err)
		}
		fmt.Println("Done!")
	},
}

func init() {
	rootCmd.AddCommand(normalCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// normalCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// normalCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	normalCmd.Flags().StringP("url", "u", "", "m3u8 地址")
	_ = normalCmd.MarkFlagRequired("url")
}
