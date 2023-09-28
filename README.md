# aliyun-m3u8-downloader

aliyun-m3u8-downloader 是一个使用了 Go 语言编写的迷你 M3U8 下载工具, 支持阿里云m3u8私有加密。 该工具就会自动帮你解析 M3U8 文件，并将 TS 片段下载下来合并成一个文件。

本工具只供学习研究，如有侵权请联系删除

## 定制
可定制开发使用以下视频云服务的第三方平台下载器，详细咨询微信：lbbniu-com
- **阿里云私有音视频加密**
- **火山引擎视频云点播**
- **百度智能云视频点播**
- **华为云视频点播**
- **气球云视频点播**
- **[保利威视 Polyv](https://www.polyv.net/)**：支持v1104(算法)、v12(算法)、v13(wasm + libx264 全网最快), 未开源
  - v13 架构：ts、h264解密使用go语言实现
  - h264解码为yuv使用wasmtime-go实现
  - yuv转h264使用libx264实现
  - 最后aac+h264合成ts使用go语言实现

### 联系开发者

![wechat](images/wechat.png)

### 插件
目前支持的闭源全自动批量下载器插件包括：
- [x] [光环国际](https://yun.aura.cn)
- [x] [中公网校](https://www.eoffcn.com)
- [x] [学培课堂](https://www.fhzjedu.com)
- [x] [云上虎](https://www.huohujiaoyu.com)
- [x] [慕课网体系课和实战课](https://www.imooc.com)
- [x] [银成医考](https://wx.yixueks.com)
- [x] [51cto](https://edu.51cto.com)
- [x] [某兽医app](https://www.med126.com/)
- [x] [极客时间训练营](https://time.geekbang.org/)
- [x] [现代卓越](https://remote.chinapm.org/)
- [ ] [好医术](https://www.haoyishu.com/)
- [ ] [知群](https://izhiqun.com/)
- [x] [马士兵](https://www.mashibing.com/)
- [ ] [百战程序员](https://www.itbaizhan.com/)
- [ ] [库课网校](https://www.kuke99.com/)
- [ ] [昭昭医考](https://www.yikao88.com/)
- [ ] [翼狐网](https://www.yiihuu.com/)
- [ ] [cgjoy课堂](https://www.cgjoy.com/h5/pages/course/index/index)
- [x] [Siki学院-免费课程](https://www.sikiedu.com/)
- [x] [短书平台](https://www.duanshu.com/) --> [太阳老师讲数学](https://hlrzp.duanshu.com) 公众号
- 图形化界面下载器，适合无计算机基础用户使用

![main](images/main.png)

## 功能

- 支持阿里云M3U8私有加密解密
- 下载和解析 M3U8（仅限 VOD 类型）
- 下载 TS 失败重试
- 解析 Master playlist
- 解密 TS
- 合并 TS 片段

## 用法

### 源码方式

```bash
# 交叉编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o aliyun-m3u8-downloader
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o aliyun-m3u8-downloader
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o aliyun-m3u8-downloader.exe
# 普通m3u8下载
go run main.go normal -u=https://www.lbbniu.com/index.m3u8 -o=/data/example --chanSize 1
# 阿里云m3u8私有加密
go run main.go aliyun -p "WebPlayAuth" -v 视频id -o=/data/example --chanSize 1
```

### 二进制方式:

Linux 和 MacOS

```
# 普通m3u8下载
./aliyun-m3u8-downloader normal -u https://www.lbbniu.com/index.m3u8 -o=/data/example -c 1
# 阿里云m3u8私有加密
./aliyun-m3u8-downloader aliyun -p "PlayAuth" -o=/data/example -c 1
# 火山引擎视频云视频下载
./aliyun-m3u8-downloader bytedance -p "PlayAuthToken" -o=/data/example -c 1
# 百度智能云视频下载
./aliyun-m3u8-downloader baidu -u m3u8视频地址 -t token  -o=/data/example -c 1
```

### 命令帮助

```shell
 aliyun-m3u8-downloader -h
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  aliyun-m3u8-downloader [command]

Available Commands:
  51cto         51cto课程下载
  aliyun        阿里云私有m3u8加密下载工具
  aura          光环国际视频下载
  baidu         百度智能云视频下载
  baidubce      baidubce视频下载
  bytedance     字节跳动，火山引擎视频云视频加密下载工具
  chinapm       现代卓越视频下载
  completion    Generate the autocompletion script for the specified shell
  eoffcn        中公网校课程下载
  fhzjedu       学培课堂课程下载
  geektime      极客时间训练营下载
  gk            极客时间训练营下载
  help          Help about any command
  huohujiaoyu   云上虎视频下载
  imooc         慕课网体系课/实战课下载
  multi         根据PlayAuth批量输出m3u8地址和解密key
  normal        普通m3u8 或 标准AES-128加密 下载
  polyv         保利威视频下载
  qiqiuyun      气球云视频下载
  veterinaryapp 某兽医app视频下载
  yixueks       银成医考课程下载

Flags:
  -c, --concurrency int     下载并发数 (default 1)
  -f, --filename string     保存文件名
  -h, --help                help for aliyun-m3u8-downloader
  -o, --output string       下载保存位置
  -r, --referer string      referer请求头
      --user-agent string   User-Agent

Use "aliyun-m3u8-downloader [command] --help" for more information about a command.
```

## 下载

[二进制文件](https://github.com/lbbniu/aliyun-m3u8-downloader/releases)

## 参考资料

- [https://github.com/SweetInk/lagou-course-downloader](https://github.com/SweetInk/lagou-course-downloader)
- [https://github.com/oopsguy/m3u8](https://github.com/oopsguy/m3u8)

## License

[MIT License](LICENSE)
