//go:build plugin
// +build plugin

package cmd

import "github.com/lbbniu/aliyun-m3u8-downloader/plugins"

func init() {
	// 插件为定制开发服务
	plugins.AddCmd(rootCmd)
}
