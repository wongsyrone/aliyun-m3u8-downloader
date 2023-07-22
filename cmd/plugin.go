//go:build plugin
// +build plugin

package cmd

import "github.com/lbbniu/aliyun-m3u8-downloader/plugins/addcmd"

func init() {
	// 插件为定制开发服务
	addcmd.AddCmd(rootCmd)
}
