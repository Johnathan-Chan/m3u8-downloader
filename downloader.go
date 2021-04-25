package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	"unsafe"
)

const (
	defaultNumberOfThread = 10
	NoException           = 0
	UrlException          = 1
	IOException           = 2
	NetworkException      = 3
	UnexpectedException   = 9
	WriteIntoCacheAndSaveModel = 1
	SaveAsTsFileAndMergeModel = 2
	SuffixMp4             = ".mp4"
	SuffixTs              = ".ts"
	TestUrl1              = "https://c.mhkuaibo.com/20200305/WzAhqGg8/1200kb/hls/index.m3u8"
	TestUrl2              = "https://v.bdcache.com/vh3/640/3c24f5d470b860351bc5e670aef38446c0a125dc/master.m3u8"
)
//https://xxx.sdhdbd1.com/52av/20210412/A%E6%97%A5%E9%9F%A9%E6%97%A0%E7%A0%81/%E5%88%B6%E6%9C%8D-%E5%B7%A8%E4%B9%B3-%E7%8E%A9%E5%85%B7-3P-%E5%87%BA%E6%9D%A5%E7%9A%84%E7%9A%84%E6%BC%8F%E4%B9%B3%E6%B1%81/SD/playlist.m3u8

var(
	defaultSaveDirectory  = []byte("./video/")
	downloadModelMap = map[DownloadModelType]DownloadModelType{}
	errorMap = map[DownloadExceptionType]error{}
)

// init 初始化函数
func init()  {
	downloadModelMap[WriteIntoCacheAndSaveModel] = WriteIntoCacheAndSaveModel
	downloadModelMap[SaveAsTsFileAndMergeModel] = SaveAsTsFileAndMergeModel
	errorMap[UrlException] = errors.New("[URLException]:Please check you url")
	errorMap[IOException] = errors.New("[IOException]:fuck")
	errorMap[NetworkException] = errors.New("[NetworkException]:fuck")
}

//自定义类型
type (
	IntChannel            chan int
	DownloadExceptionType int
	DownloadModelType int
	DownloadModelFunction func(int,int,[]byte)
)

// M3u8Downloader 下载器接口对象
type M3u8Downloader interface {
	// ParseM3u8File 解析.m3u8文件
	ParseM3u8File(url string)
	// Download 根据配置信息执行下载
	Download() error
	// SetUrl 设置url
	SetUrl(url string)
	// SetIfShowTheBar 是否显示进度条
	SetIfShowTheBar(ifShow bool)
	// SetNumOfThread 设置下载线程的数量
	SetNumOfThread(num int)
	// SetMovieName 设置视频的文件名
	SetMovieName(videoName string)
	// SetSaveDirectory 设置保存目录
	SetSaveDirectory(targetDir string)
	// SetDownloadModel 设置下载模式
	SetDownloadModel(model DownloadModelType)
	// MergeFile 默认合并文件
	MergeFile() error
	// MergeFileToDir 将合并后的视频文件保存到目录dir中
	MergeFileToDir(dir string) error

}

// m3u8downloader 下载器对象结构体
type m3u8downloader struct {
	// config 下载配置
	config *DownloadConfig
	// taskChannel 发布下载任务的管道
	taskChannel IntChannel
	// suffixList 需要下载的url列表
	suffixList [][]byte
	//waitGroup
	waitGroup *sync.WaitGroup
	// buffer 缓冲区
	buffer []StringBuilder
	// cacheMap
	cacheMap map[int][]byte
	//检查下载中的错误
	exception DownloadExceptionType
}

// DownloadConfig 下载配置对象
// 因为最后还有很多拼接字符串的地方，所以我们尽可能使用[]byte和StringBuilder
// 以减少不必要的开销
type DownloadConfig struct {
	// Url 下载链接
	Url []byte
	// noSuffixUrl 经过处理后得到无后缀的url
	noSuffixUrl []byte
	// NumOfThreads 下载的线程数
	NumOfThreads int
	// VideoName 保存的视频名称
	VideoName []byte
	// SaveDirectory 保存视频的目录
	SaveDirectory []byte
	// ifShowBar 是否在控制台显示进度条
	ifShowBar bool
	// errCount 下载出错的数量
	errCount int
	// completeCount 下载完成的数量
	completeCount int
	// TotalNum 总共的下载数量
	TotalNum int
	// DownloadModel 设置下载模式
	DownloadModel DownloadModelType
}

// NewDownloader 创建一个新对象
func NewDownloader() M3u8Downloader {
	return newDownload( &DownloadConfig{
		NumOfThreads:  defaultNumberOfThread,
		ifShowBar:     false,
		SaveDirectory: []byte(defaultSaveDirectory),
		Url:           nil,
	})
}

// NewDownloaderWithConfig 使用自定义配置创建下载器对象
func NewDownloaderWithConfig(config *DownloadConfig) M3u8Downloader {
	return newDownload(config)
}

// newDownload
func newDownload(config *DownloadConfig)M3u8Downloader{
	return &m3u8downloader{
		config:     config,
		buffer:     make([]StringBuilder, defaultNumberOfThread),
		cacheMap:   map[int][]byte{},
		suffixList: make([][]byte, 0),
		exception:  NoException,
		waitGroup: &sync.WaitGroup{},
		taskChannel: make(IntChannel,50),
	}
}


//https://v.bdcache.com/vh3/640/3c24f5d470b860351bc5e670aef38446c0a125dc/master.m3u8
//https://c.mhkuaibo.com/20200305/WzAhqGg8/1200kb/hls/index.m3u8
func (md *m3u8downloader) simplyHttpGet(url string) []byte {
	res, err := http.Get(url)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		if err.Error()[len(err.Error())-12:] == "no such host" {
			md.exception = NetworkException
		}
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		md.exception = IOException
		return nil
	}
	//file not found
	//[102 105 108 101 32 (110)(n) 111 116 32 (102)(f) 111 117 110 100]
	if body[5] == 110 && body[9] == 102 {
		md.exception = UrlException
		return nil
	}
	return body
}

// ParseM3u8File 解析并处理m3u8文件
func (md *m3u8downloader) ParseM3u8File(url string) {
	body := md.simplyHttpGet(url)
	if md.exception != NoException {
		return
	}
	var i, left, bodyLen=0,0,len(body) - 1
	for  i < bodyLen  {
		if body[i] == '/' || body[i] == '\n' {
			i+=1
			left = i
		} else if body[i] == '.' && body[i+1] == 't' {
			i+=3
			md.suffixList = append(md.suffixList, body[left:i])
		}else{
			i++
		}
	}
	//println(len(md.suffixList))
}

// SetUrl 设置需要下载视频的.m3u8文件的url
func (md *m3u8downloader) SetUrl(url string) {
	//设置url
	md.config.Url = []byte(url)
	//查找url中的最后一个‘/’,并作为noSuffixUrl的相关参数来构造noSuffixUrl
	last:=md.reFind(len(md.config.Url)-1,md.config.Url)
	md.config.noSuffixUrl = make([]byte,last)
	md.config.noSuffixUrl = md.config.Url[0:last+1]
}

// getUnixTimeByte 根据加当前时间戳设置为默认名称
func getUnixTimeByte() []byte {
	t1:=time.Now().Unix()
	var temp int64
	var buf  = make([]byte,10)
	var i = 9
	for t1>0 {
		temp = t1%10
		buf[i] = (*(*byte)(unsafe.Pointer(&temp)))+48
		i--
		t1/=10
	}
	return buf//*(*string)(unsafe.Pointer(&buf))
}

// SetIfShowTheBar 是否显示进度条
func (md *m3u8downloader) SetIfShowTheBar(ifShow bool) {
	md.config.ifShowBar = ifShow
}

// SetNumOfThread 设置线程数量
func (md *m3u8downloader) SetNumOfThread(num int) {
	md.config.NumOfThreads = num
	md.buffer = make([]StringBuilder, num)
}

// SetMovieName 设置保存后的视频名称
func (md *m3u8downloader) SetMovieName(videoName string) {
	md.config.VideoName = []byte(videoName)
}

// SetSaveDirectory 设置下载的视频的保存路径
func (md *m3u8downloader) SetSaveDirectory(targetDir string) {
	md.config.SaveDirectory = []byte(targetDir)
}

// SetDownloadModel 设置下载模式
func (md *m3u8downloader) SetDownloadModel(model DownloadModelType) {
	var ok bool
	md.config.DownloadModel,ok = downloadModelMap[model]
	if !ok {
		md.config.DownloadModel = WriteIntoCacheAndSaveModel
	}
}

// showTheBar 显示进度条方法
func (md *m3u8downloader) showTheBar() {
	var i,total int64
	total = int64(md.config.TotalNum)
	bar := NewBar(total)
	bar.Setting().SetShowModel(WindowsCmd)
	for i = 0; i<=total; i++{
		time.Sleep(1000*time.Millisecond)
		bar.Play(i)
	}
	bar.Finish()
}

// reFind 倒叙查找
func (md *m3u8downloader) reFind(startIndex int,str []byte) int {
	var i int
	for i = startIndex;str[i]!='/';i--{}
	return i
}

// ThrowsException 抛出异常已弃用
func (md *m3u8downloader)ThrowsException() error{
	switch md.exception {
	case UrlException:
		return errors.New("[URLException]:Please check you url")
	case IOException:
		return errors.New("[IOException]:fuck")
	case NetworkException:
		return errors.New("[NetworkException]:fuck")
	default :
		return nil
	}
}

// Download 下载任务核心方法
func (md *m3u8downloader) Download() error {
	//首先设置url，然后获取.mu3u8文件并解析，测试网络是否畅通
	//根据.m3u8文件拼接好保存目录，文件名，
	md.ParseM3u8File(*(*string)(unsafe.Pointer(&md.config.Url)))
	if md.exception != NoException {
		return errorMap[md.exception]
	}
	//check if is default file name
	if md.config.VideoName == nil {
		//如果没有设置name,使用时间戳来代替
		md.config.VideoName = getUnixTimeByte()
	}
	//将参数归零
	md.config.errCount = 0
	md.config.completeCount = 0
	md.config.TotalNum = len(md.suffixList)
	md.waitGroup.Add(md.config.NumOfThreads)
	//选择下载模式
	var callBackFunc DownloadModelFunction
	if md.config.DownloadModel == WriteIntoCacheAndSaveModel{
		go md.WriteIntoCacheAndSaveProcessor()
		callBackFunc = md.WriteIntoCacheAndSave
	}else{
		callBackFunc = md.SaveAsTsFileAndMerge
	}
	//发布下载任务
	go md.publisher()
	//开启下载线程
	for i := 0; i < md.config.NumOfThreads; i++ {
		go md.download(i,callBackFunc)
	}
	//显示进度条
	if md.config.ifShowBar {
		go md.showTheBar()
	}
	//阻塞等待
	md.waitGroup.Wait()
	//检查异常
	if md.exception != NoException {
		return errorMap[md.exception]
	}
	return nil
}

// publisher 任务发布
func (md *m3u8downloader) publisher() {
	for i:=0; i<md.config.TotalNum; i++{
		md.taskChannel <- i
	}
	close(md.taskChannel)
}

func (md *m3u8downloader) download(threadId int,downloadModel func(int,int,[]byte)) {
	var index int
	var ok bool
	var body []byte
	for {
		//从管道接收信息
		index,ok = <- md.taskChannel
		if !ok {
			break
		} //拼接路径
		md.buffer[threadId].Write(md.config.noSuffixUrl)
		md.buffer[threadId].Write(md.suffixList[index])
		//尝试下载，错误达到一定次数停止下载
		for{
			body = md.simplyHttpGet(md.buffer[threadId].String())
			if md.config.errCount < md.config.NumOfThreads {
				if md.exception != NoException {
					//如果出现范围允许的错误，则重试
					md.config.errCount++
					md.exception = NoException
				}//若没有出现错误则执行回调函数
				break
			}else{
				//若出现严重错误，则通知其他线程，停止工作
				md.waitGroup.Done()
				return
			}
		}
		//执行下载回调函数
		downloadModel(index,threadId,body)
	}
	md.waitGroup.Done()
}

// WriteIntoCacheAndSaveProcessor 使用缓存下载模式的处理器函数
func (md *m3u8downloader) WriteIntoCacheAndSaveProcessor(){
	var buffer StringBuilder
	buffer.Write(md.config.SaveDirectory)
	buffer.Write(md.config.VideoName)
	movie, err := os.OpenFile(buffer.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		md.config.errCount+=md.config.NumOfThreads
		md.exception = IOException
		return
	}
	defer movie.Close()
	var body []byte
	var i int
	var ok bool
	//自旋
	for i < md.config.TotalNum {
		//如果字典中有了对应索引的id相关的内容，则追加写入然后继续新的尝试，
		//否则睡眠等待一定时间后再次尝试，以减少资源的消耗，直到成功为止
		if body,ok = md.cacheMap[i]; ok {
			movie.Write(body)
			md.cacheMap[i] = nil
			delete(md.cacheMap,i)
			i++
		}else{
			time.Sleep(250*time.Millisecond)
		}
	}
	//ReSet The Map And Help GC
	md.cacheMap = nil
	md.cacheMap = map[int][]byte{}
}


// WriteIntoCacheAndSave 写入缓存，最后保存
func (md *m3u8downloader) WriteIntoCacheAndSave(index,threadId int,body []byte){
	md.cacheMap[index] = body
	md.config.completeCount++
	md.buffer[threadId].Reset()
}


// SaveAsTsFileAndMerge 下载为ts文件，最后手动合并
func (md *m3u8downloader) SaveAsTsFileAndMerge(index,threadId int,body []byte){
	md.buffer[threadId].Reset()
	md.buffer[threadId].Write(defaultSaveDirectory)
	md.buffer[threadId].Write(md.suffixList[index])

	md.suffixList[index] = md.buffer[threadId].GetBuffer()

	movie, err := os.OpenFile(md.buffer[threadId].String(), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		md.config.errCount+=md.config.NumOfThreads
		md.exception = IOException
		return
	}
	movie.Write(body)
	movie.Close()
	md.buffer[threadId].Reset()
	md.config.completeCount++
	body = nil
}


// MergeFile 合并文件
func (md *m3u8downloader) MergeFile() error {
	return md.mergeFile()
}

// MergeFileToDir 合并文件
func (md *m3u8downloader) MergeFileToDir(dir string) error {
	md.config.SaveDirectory = []byte(dir)
	return md.mergeFile()
}

// MergeFileToCompletelyDir 合并文件
func (md *m3u8downloader) MergeFileToCompletelyDir(dir string,movieName string) error {
	md.config.SaveDirectory = []byte(dir)
	md.config.VideoName = []byte(movieName)
	return  md.mergeFile()
}


// mergeFile
func (md *m3u8downloader) mergeFile() error{
	var buffer StringBuilder
	var err error
	var movie,tsFile *os.File
	buffer.Write(md.config.SaveDirectory)
	buffer.Write(md.config.VideoName)
	movie,err = os.OpenFile(buffer.String(),os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err!= nil{
		return err
	}
	defer movie.Close()
	var i int
	var body []byte
	var tempStr string
	for i=0;i<len(md.suffixList);i++ {
		tempStr = *(*string)(unsafe.Pointer(&md.suffixList[i]))
		tsFile,err= os.OpenFile(tempStr,os.O_CREATE|os.O_RDONLY, os.ModePerm)
		if err!= nil{
			return err
		}
		body,err = ioutil.ReadAll(tsFile)
		if err!= nil{
			return err
		}
		_,err = movie.Write(body)
		if err!= nil{
			return err
		}
		tsFile.Close()
		err = os.Remove(tempStr)
		if err != nil {
			return err
		}
		md.suffixList[i] = nil
	}
	return nil
}


//func httpDo() {
//	client := &http.Client{}
//
//	req, err := http.NewRequest("POST", "http://www.01happy.com/demo/accept.php", strings.NewReader("name=cjb"))
//	if err != nil {
//		// handle error
//	}
//
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Set("Cookie", "name=anny")
//
//	resp, err := client.Do(req)
//
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//
//	fmt.Println(string(body))
//}
