/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/lbbniu/aliyun-m3u8-downloader/cmd"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	cmd.Execute()
}
