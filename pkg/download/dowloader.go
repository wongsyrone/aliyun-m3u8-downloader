package download

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/log"
	fluentffmpeg "github.com/modfy/fluent-ffmpeg"

	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

type MergeTsType int

const (
	Golang MergeTsType = iota
	Ffmpeg
)

const (
	tsFolderName     = "ts"
	tsTempFileSuffix = "_tmp"
	progressWidth    = 40
)

type decryptFunc func(int, string, []byte, *parse.Segment, *parse.KeyInfo) ([]byte, error)

type Downloader struct {
	lock            sync.Mutex
	queue           []int
	folder          string
	tsFolder        string
	finish          int32
	segLen          int
	mergeTSFilename string

	result *parse.Result

	url         string
	output      string
	filename    string
	key         string
	loadKeyFunc parse.LoadKeyFunc
	decryptFunc decryptFunc
	// mp4 下载
	mp4    bool
	mp4Url string

	mergeTsType MergeTsType
	// m3u8Content
	m3u8Content string
}

type DownloaderOption func(*Downloader)

func WithUrl(url string) DownloaderOption {
	return func(d *Downloader) {
		d.url = url
	}
}

func WithOutput(output string) DownloaderOption {
	return func(d *Downloader) {
		d.output = output
	}
}

func WithKey(key string) DownloaderOption {
	return func(d *Downloader) {
		d.key = key
		d.loadKeyFunc = func(_, _ string) (string, error) {
			return d.key, nil
		}
	}
}

func WithLoadKeyFunc(loadKeyFunc parse.LoadKeyFunc) DownloaderOption {
	return func(d *Downloader) {
		d.loadKeyFunc = loadKeyFunc
	}
}

func WithFilename(filename string) DownloaderOption {
	return func(d *Downloader) {
		d.filename = filename
	}
}

func WithMp4(mp4 bool) DownloaderOption {
	return func(d *Downloader) {
		d.mp4 = mp4
	}
}

func WithM3u8Content(m3u8Content string) DownloaderOption {
	return func(d *Downloader) {
		d.m3u8Content = m3u8Content
	}
}

func WithMergeTsType(mergeTsType MergeTsType) DownloaderOption {
	return func(d *Downloader) {
		d.mergeTsType = mergeTsType
	}
}

func WithDecryptFunc(decryptFunc decryptFunc) DownloaderOption {
	return func(d *Downloader) {
		d.decryptFunc = decryptFunc
	}
}

func loadKeyFunc(_, keyUrl string) (string, error) {
	resp, err := httpclient.Get(keyUrl)
	if err != nil {
		return "", fmt.Errorf("download: extract key failed: %w", err)
	}
	keyStr, err := resp.ToString()
	if err != nil {
		return "", fmt.Errorf("download: ToString: %w", err)
	}
	//log.Debugf("decryption key: %s", keyStr)
	return keyStr, err
}

// NewDownloader returns a Task instance
func NewDownloader(opts ...DownloaderOption) (*Downloader, error) {
	d := &Downloader{
		loadKeyFunc: loadKeyFunc,
	}
	for _, opt := range opts {
		opt(d)
	}

	// 处理保存目录
	d.folder = d.output
	// If no output folder specified, use current directory
	if d.output == "" {
		current, err := tool.CurrentDir()
		if err != nil {
			return nil, err
		}
		d.folder = current
	}
	// 创建保存目录
	if err := os.MkdirAll(d.folder, os.ModePerm); err != nil {
		return nil, fmt.Errorf("donwload: create storage folder failed: %s", err.Error())
	}

	// 解析合并的最终文件名
	if d.filename != "" {
		d.mergeTSFilename = strings.TrimSuffix(d.filename, ".mp4") + ".mp4"
	} else {
		d.mergeTSFilename = d.tsFilename(d.url) + ".mp4"
	}
	if d.url == "" && d.m3u8Content == "" {
		return nil, fmt.Errorf("donwload: url: %s and m3u8Content: %s", d.url, d.m3u8Content)
	}

	if d.mp4 {
		d.mp4Url = d.url
	} else {
		// 构造ts文件目录
		d.tsFolder = filepath.Join(d.folder, tsFolderName)
		// 创建ts文件目录
		if err := os.MkdirAll(d.tsFolder, os.ModePerm); err != nil {
			return nil, fmt.Errorf("donwload: create ts folder '[%s]' failed: %s", d.tsFolder, err.Error())
		}

		// 解析m3u8文件内容
		var err error
		if d.m3u8Content != "" {
			d.result, err = parse.FromM3u8Content(d.url, d.m3u8Content, d.loadKeyFunc)
		} else if d.url != "" {
			d.result, err = parse.FromM3u8URL(d.url, d.loadKeyFunc)
		} else {
			return nil, fmt.Errorf("donwload: url: %s and m3u8Content: %s all empty", d.url, d.m3u8Content)
		}
		if err != nil {
			return nil, fmt.Errorf("donwload: parse m3u8 err: %s", err)
		}
		d.segLen = len(d.result.M3u8.Segments)
		d.queue = d.genSlice(d.segLen)
	}
	return d, nil
}

func (d *Downloader) SetDecryptFunc(decryptFunc decryptFunc) {
	d.decryptFunc = decryptFunc
}

// Start runs downloader
func (d *Downloader) Start(concurrency int) error {
	if d.mp4 {
		return d.downloadMp4(d.mp4Url)
	}
	var wg sync.WaitGroup
	// struct{} zero size
	limitChan := make(chan struct{}, concurrency)
	for {
		tsIdx, end, err := d.next()
		if err != nil {
			if end {
				break
			}
			continue
		}
		wg.Add(1)
		limitChan <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			if er := d.download(idx); er != nil {
				// Back into the queue, retry request
				log.Errorf("[failed] %v", er)
				if er = d.back(idx); er != nil {
					log.Error(er)
				}
			}
			<-limitChan
		}(tsIdx)
	}
	wg.Wait()
	if err := d.mergeTsToMp4(); err != nil {
		return fmt.Errorf("download: merge ts to mp4 err: %w", err)
	}
	return nil
}

func (d *Downloader) downloadMp4(mp4Url string) error {
	resp, err := httpclient.Get(mp4Url)
	if err != nil {
		return fmt.Errorf("download: request %s, err: %w", mp4Url, err)
	}
	// Create a mp4 file
	mFilePath := filepath.Join(d.folder, d.mergeTSFilename)
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("download: create mp4 file failed：%w", err)
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()
	defer resp.Body.Close()
	_, err = io.Copy(mFile, resp.Body)
	if err != nil {
		return fmt.Errorf("download: write mp4 file failed：%w", err)
	}
	return nil
}

func (d *Downloader) download(segIndex int) error {
	tsUrl := d.tsURL(segIndex)
	resp, err := httpclient.Get(tsUrl)
	if err != nil {
		return fmt.Errorf("download: request %s, %s", tsUrl, err.Error())
	}
	filename := d.tsFilename(tsUrl)
	//noinspection GoUnhandledErrorResult
	fPath := filepath.Join(d.tsFolder, filename)
	fTemp := fPath + tsTempFileSuffix
	f, err := os.Create(fTemp)
	if err != nil {
		return fmt.Errorf("download: create file: %s, %s", filename, err.Error())
	}
	tsData, err := resp.ReadAll()
	if err != nil {
		return fmt.Errorf("download: read bytes: %s, %s", tsUrl, err.Error())
	}
	sf := d.result.M3u8.Segments[segIndex]
	if sf == nil {
		return fmt.Errorf("download: invalid segment index: %d", segIndex)
	}
	keyInfo, ok := d.result.M3u8.Keys[sf.KeyIndex]
	if ok {
		if d.decryptFunc != nil {
			// 自定义解密函数
			tsData, err = d.decryptFunc(segIndex, fPath, tsData, sf, keyInfo)
			if err != nil {
				// TODO: 如果当前片段时长大于1秒，则认为解密失败
				seg := d.segment(segIndex)
				if seg.Duration > 1 {
					return fmt.Errorf("download: decryptFunc: %s, err: %w", tsUrl, err)
				}
				log.Errorf("segment: %+v, err: %v", seg, tsUrl, err)
				tsData = nil
			}
		} else if keyInfo.AliyunVoDEncryption {
			// 阿里云私有加密
			tsParser := aliyun.NewTSParser(tsData, keyInfo.Key)
			tsData = tsParser.Decrypt()
		} else if keyInfo.Key != "" {
			// 标准AES-128解密
			tsData, err = tool.AES128Decrypt(tsData, []byte(keyInfo.Key), []byte(keyInfo.IV))
			if err != nil {
				return fmt.Errorf("download: decrypt: %s, err: %w", tsUrl, err)
			}
		}
	}

	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) //0x47
	bLen := len(tsData)
	for j := 0; j < bLen; j++ {
		if tsData[j] == syncByte {
			tsData = tsData[j:]
			break
		}
	}

	w := bufio.NewWriter(f)
	if _, err := w.Write(tsData); err != nil {
		return fmt.Errorf("download: write to %s: %s", fTemp, err.Error())
	}
	// Release file resource to rename file
	_ = f.Close()
	if err = os.Rename(fTemp, fPath); err != nil {
		return err
	}
	// Maybe it will be safer in this way...
	atomic.AddInt32(&d.finish, 1)
	tool.DrawProgressBar(fmt.Sprintf("downloading %d/%d", d.finish, d.segLen), float32(d.finish)/float32(d.segLen), progressWidth)
	//log.Infof("[download %6.2f%%] %s", float32(d.finish)/float32(d.segLen)*100, tsUrl)
	return nil
}

func (d *Downloader) next() (segIndex int, end bool, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	if len(d.queue) == 0 {
		err = fmt.Errorf("queue empty")
		if d.finish == int32(d.segLen) {
			end = true
			return
		}
		// Some segment indexes are still running.
		end = false
		return
	}
	segIndex = d.queue[0]
	d.queue = d.queue[1:]
	return
}

func (d *Downloader) back(segIndex int) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if sf := d.result.M3u8.Segments[segIndex]; sf == nil {
		return fmt.Errorf("download: invalid segment index: %d", segIndex)
	}
	d.queue = append(d.queue, segIndex)
	return nil
}

func (d *Downloader) mergeTsToMp4() error {
	// In fact, the number of downloaded segments should be equal to number of m3u8 segments
	missingCount := 0
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		f := filepath.Join(d.tsFolder, d.tsFilename(d.tsURL(segIndex)))
		if _, err := os.Stat(f); err != nil {
			missingCount++
		}
	}
	if missingCount > 0 {
		log.Infof("[warning] %d files missing", missingCount)
	}

	// Create a TS file for merging, all segment files will be written to this file.
	mFilePath := filepath.Join(d.folder, d.mergeTSFilename)
	defer func() {
		_ = os.RemoveAll(d.tsFolder)
	}()
	switch d.mergeTsType {
	case Ffmpeg:
		return d.mergeTsToMp4ByFfmpeg(mFilePath)
	default:
		return d.mergeTsToMp4ByGo(mFilePath)
	}
}

func (d *Downloader) mergeTsToMp4ByFfmpeg(mFilePath string) error {
	log.Infof("%s 开始合并", mFilePath)
	defer func(startTime time.Time) {
		log.Infof("%s 合并结束，耗时: %.0fs", mFilePath, time.Since(startTime).Seconds())
	}(time.Now())
	var tsFiles []string
	// 设置输入文件列表
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		tsFile := filepath.Join(d.tsFolder, d.tsFilename(d.tsURL(segIndex)))
		tsFiles = append(tsFiles, tsFile)
	}
	buf := &bytes.Buffer{}
	ffmpeg := fluentffmpeg.NewCommand("")
	ffmpeg.InputPath(fmt.Sprintf("concat:%s", strings.Join(tsFiles, "|"))).
		OutputPath(mFilePath).
		OutputLogs(buf).
		Overwrite(true).
		FromFormat("mpegts").
		VideoCodec("copy").
		AudioCodec("aac")
	if err := ffmpeg.Run(); err != nil {
		out, _ := io.ReadAll(buf)
		log.Info(string(out))
		return err
	}
	return nil
}

func (d *Downloader) mergeTsToMp4ByGo(mFilePath string) error {
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("download: create main TS file failed：%w", err)
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()

	fmt.Println()
	writer := bufio.NewWriter(mFile)
	mergedCount := 0
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		bytes, err := os.ReadFile(filepath.Join(d.tsFolder, d.tsFilename(d.tsURL(segIndex))))
		if _, err = writer.Write(bytes); err != nil {
			continue
		}
		mergedCount++
		tool.DrawProgressBar(
			fmt.Sprintf("merge       %d/%d", mergedCount, d.segLen),
			float32(mergedCount)/float32(d.segLen),
			progressWidth,
		)
	}
	_ = writer.Flush()
	fmt.Println()
	if mergedCount != d.segLen {
		log.Warnf("[warning] %d files merge failed", d.segLen-mergedCount)
	}
	log.Infof("[output] %s", mFilePath)
	return nil
}

func (d *Downloader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return tool.ResolveURL(d.result.URL, seg.URI)
}

func (d *Downloader) segment(segIndex int) *parse.Segment {
	return d.result.M3u8.Segments[segIndex]
}

func (d *Downloader) tsFilename(tsUrl string) string {
	idx := strings.Index(tsUrl, "?")
	if idx > 0 {
		return path.Base(tsUrl[:idx])
	}
	return path.Base(tsUrl)
}

func (d *Downloader) genSlice(len int) []int {
	s := make([]int, 0)
	for i := 0; i < len; i++ {
		s = append(s, i)
	}
	return s
}
