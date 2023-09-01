package cmd

import (
	"os"

	"github.com/TarsCloud/TarsGo/tars/util/rogger"
	"github.com/ddliu/go-httpclient"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/log"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aliyun-m3u8-downloader",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		concurrency := viper.GetInt("concurrency")
		if concurrency <= 0 {
			tool.PanicParameter("parameter 'concurrency' must be greater than 0")
		}

		config := httpclient.Map{}
		if referer := viper.GetString("referer"); referer != "" {
			config[httpclient.OPT_REFERER] = referer
		}
		if ua := viper.GetString("user-agent"); ua != "" {
			config[httpclient.OPT_USERAGENT] = ua
		}
		httpclient.Defaults(config)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.Init()
	var err error
	defer func() {
		rogger.FlushLogger()
		if err != nil {
			os.Exit(1)
		}
	}()
	err = rootCmd.Execute()
}

func init() {
	// cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aliyun-m3u8-downloader.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "", "下载保存位置")
	rootCmd.PersistentFlags().StringP("filename", "f", "", "保存文件名")
	rootCmd.PersistentFlags().IntP("concurrency", "c", 1, "下载并发数")
	rootCmd.PersistentFlags().StringP("referer", "r", "", "referer请求头")
	rootCmd.PersistentFlags().StringP("user-agent", "", "", "User-Agent")

	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("filename", rootCmd.PersistentFlags().Lookup("filename"))
	viper.BindPFlag("concurrency", rootCmd.PersistentFlags().Lookup("concurrency"))
	viper.BindPFlag("referer", rootCmd.PersistentFlags().Lookup("referer"))
	viper.BindPFlag("user-agent", rootCmd.PersistentFlags().Lookup("user-agent"))
	//viper.BindPFlags(rootCmd.PersistentFlags())

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
