/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/TarsCloud/TarsGo/tars/util/rogger"

	"github.com/lbbniu/aliyun-m3u8-downloader/cmd"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/log"
)

func main() {
	defer rogger.FlushLogger()
	log.Init()
	cmd.Execute()
}
