package download

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ddliu/go-httpclient"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/parse/aliyun"
	"github.com/lbbniu/aliyun-m3u8-downloader/pkg/tool"
)

const (
	tsFolderName     = "ts"
	tsTempFileSuffix = "_tmp"
	progressWidth    = 40
)

type Downloader struct {
	lock            sync.Mutex
	queue           []int
	folder          string
	tsFolder        string
	finish          int32
	segLen          int
	mergeTSFilename string

	result *parse.Result

	output      string
	filename    string
	key         string
	loadKeyFunc parse.LoadKeyFunc
	// mp4 下载
	mp4    bool
	mp4Url string
}

type DownloaderOption func(*Downloader)

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

func loadKeyFunc(_, keyUrl string) (string, error) {
	resp, err := httpclient.Get(keyUrl)
	if err != nil {
		return "", fmt.Errorf("extract key failed: %w", err)
	}
	keyStr, err := resp.ToString()
	if err != nil {
		return "", err
	}
	//fmt.Println("decryption key: ", keyStr)
	return keyStr, err
}

// NewDownloader returns a Task instance
func NewDownloader(url string, opts ...DownloaderOption) (*Downloader, error) {
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
		return nil, fmt.Errorf("create storage folder failed: %s", err.Error())
	}

	// 解析合并的最终文件名
	if d.filename != "" {
		d.mergeTSFilename = strings.TrimSuffix(d.filename, ".mp4") + ".mp4"
	} else {
		d.mergeTSFilename = tsFilename(url) + ".mp4"
	}

	if d.mp4 {
		d.mp4Url = url
	} else {
		// 构造ts文件目录
		d.tsFolder = filepath.Join(d.folder, tsFolderName)
		// 创建ts文件目录
		if err := os.MkdirAll(d.tsFolder, os.ModePerm); err != nil {
			return nil, fmt.Errorf("create ts folder '[%s]' failed: %s", d.tsFolder, err.Error())
		}

		// 解析m3u8文件内容
		var err error
		d.result, err = parse.FromURL(url, d.loadKeyFunc)
		if err != nil {
			return nil, err
		}
		d.segLen = len(d.result.M3u8.Segments)
		d.queue = genSlice(d.segLen)
	}
	return d, nil
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
		go func(idx int) {
			defer wg.Done()
			if er := d.download(idx); er != nil {
				// Back into the queue, retry request
				fmt.Printf("[failed] %v\n", er)
				if er = d.back(idx); er != nil {
					fmt.Println(er)
				}
			}
			<-limitChan
		}(tsIdx)
		limitChan <- struct{}{}
	}
	wg.Wait()
	if err := d.mergeHsToMp4(); err != nil {
		return err
	}
	return nil
}

func (d *Downloader) downloadMp4(mp4Url string) error {
	resp, err := httpclient.Get(mp4Url)
	if err != nil {
		return fmt.Errorf("request %s, err: %w", mp4Url, err)
	}
	// Create a mp4 file
	mFilePath := filepath.Join(d.folder, d.mergeTSFilename)
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("create mp4 file failed：%w", err)
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()
	defer resp.Body.Close()
	_, err = io.Copy(mFile, resp.Body)
	if err != nil {
		return fmt.Errorf("write mp4 file failed：%w", err)
	}
	return nil
}

func (d *Downloader) download(segIndex int) error {
	tsUrl := d.tsURL(segIndex)
	resp, err := httpclient.Get(tsUrl)
	if err != nil {
		return fmt.Errorf("request %s, %s", tsUrl, err.Error())
	}
	filename := tsFilename(tsUrl)
	//noinspection GoUnhandledErrorResult
	fPath := filepath.Join(d.tsFolder, filename)
	fTemp := fPath + tsTempFileSuffix
	f, err := os.Create(fTemp)
	if err != nil {
		return fmt.Errorf("create file: %s, %s", filename, err.Error())
	}
	tsData, err := resp.ReadAll()
	if err != nil {
		return fmt.Errorf("read bytes: %s, %s", tsUrl, err.Error())
	}
	sf := d.result.M3u8.Segments[segIndex]
	if sf == nil {
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	keyInfo, ok := d.result.M3u8.Keys[sf.KeyIndex]
	if ok && keyInfo.Key != "" {
		// 是否阿里云私有加密
		if keyInfo.AliyunVoDEncryption {
			tsParser := aliyun.NewTSParser(tsData, keyInfo.Key)
			tsData = tsParser.Decrypt()
		} else {
			tsData, err = tool.AES128Decrypt(tsData, []byte(keyInfo.Key), []byte(keyInfo.IV))
			if err != nil {
				return fmt.Errorf("decryt: %s, err: %w", tsUrl, err)
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
		return fmt.Errorf("write to %s: %s", fTemp, err.Error())
	}
	// Release file resource to rename file
	_ = f.Close()
	if err = os.Rename(fTemp, fPath); err != nil {
		return err
	}
	// Maybe it will be safer in this way...
	atomic.AddInt32(&d.finish, 1)
	//tool.DrawProgressBar("Downloading", float32(d.finish)/float32(d.segLen), progressWidth)
	fmt.Printf("[download %6.2f%%] %s\n", float32(d.finish)/float32(d.segLen)*100, tsUrl)
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
		return fmt.Errorf("invalid segment index: %d", segIndex)
	}
	d.queue = append(d.queue, segIndex)
	return nil
}

func (d *Downloader) mergeHsToMp4() error {
	// In fact, the number of downloaded segments should be equal to number of m3u8 segments
	missingCount := 0
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		f := filepath.Join(d.tsFolder, tsFilename(d.tsURL(segIndex)))
		if _, err := os.Stat(f); err != nil {
			missingCount++
		}
	}
	if missingCount > 0 {
		fmt.Printf("[warning] %d files missing\n", missingCount)
	}

	// Create a TS file for merging, all segment files will be written to this file.
	mFilePath := filepath.Join(d.folder, d.mergeTSFilename)
	mFile, err := os.Create(mFilePath)
	if err != nil {
		return fmt.Errorf("create main TS file failed：%s", err.Error())
	}
	//noinspection GoUnhandledErrorResult
	defer mFile.Close()

	writer := bufio.NewWriter(mFile)
	mergedCount := 0
	for segIndex := 0; segIndex < d.segLen; segIndex++ {
		bytes, err := ioutil.ReadFile(filepath.Join(d.tsFolder, tsFilename(d.tsURL(segIndex))))
		_, err = writer.Write(bytes)
		if err != nil {
			continue
		}
		mergedCount++
		tool.DrawProgressBar("merge", float32(mergedCount)/float32(d.segLen), progressWidth)
	}
	_ = writer.Flush()
	// Remove `ts` folder
	_ = os.RemoveAll(d.tsFolder)

	if mergedCount != d.segLen {
		fmt.Printf("[warning] \n%d files merge failed", d.segLen-mergedCount)
	}

	fmt.Printf("\n[output] %s\n", mFilePath)

	return nil
}

func (d *Downloader) tsURL(segIndex int) string {
	seg := d.result.M3u8.Segments[segIndex]
	return tool.ResolveURL(d.result.URL, seg.URI)
}

func tsFilename(tsUrl string) string {
	idx := strings.Index(tsUrl, "?")
	if idx > 0 {
		return path.Base(tsUrl[:idx])
	}
	return path.Base(tsUrl)
}

func genSlice(len int) []int {
	s := make([]int, 0)
	for i := 0; i < len; i++ {
		s = append(s, i)
	}
	return s
}
