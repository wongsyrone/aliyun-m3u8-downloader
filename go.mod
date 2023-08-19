module github.com/lbbniu/aliyun-m3u8-downloader

go 1.18

require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/TarsCloud/TarsGo v1.4.4
	github.com/avast/retry-go/v4 v4.3.4
	github.com/bitly/go-simplejson v0.5.0
	github.com/bytecodealliance/wasmtime-go/v11 v11.0.0
	github.com/ddliu/go-httpclient v0.6.9
	github.com/google/uuid v1.3.0
	github.com/modfy/fluent-ffmpeg v0.1.0
	github.com/moontrade/wavm-go v0.3.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/shirou/gopsutil/v3 v3.23.5
	github.com/spf13/cobra v1.3.0
	github.com/stretchr/testify v1.8.4
	github.com/tdewolff/minify v2.3.6+incompatible
	github.com/wasmerio/wasmer-go v1.0.4
	github.com/yapingcat/gomedia v0.0.0-20230809172329-1ca40e5ed176
	k8s.io/klog/v2 v2.100.1
)

replace (
	github.com/TarsCloud/TarsGo v1.4.4 => github.com/TarsCloud/TarsGo v1.4.5-rc1.0.20230728020733-66bb6e14019d
	github.com/wasmerio/wasmer-go v1.0.4 => github.com/lbbniu/wasmer-go v0.0.0-20230717095824-bb9598b7bd12
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tdewolff/parse v2.3.4+incompatible // indirect
	github.com/tdewolff/test v1.0.9 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.uber.org/automaxprocs v1.5.2 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
