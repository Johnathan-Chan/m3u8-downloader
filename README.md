# M3u8Downloader

### 介绍
基于golang对.m3u8链接的多线程下载模块,如果配合命令行解析或图形界面，可以作为一个运行在命令行或桌面上的下载器
，也可以作为一个模块，配合爬虫一起使用，当前还有一些不完善的地方，但使用是完全ok,有时间再更新,如果有更好的想法
或方案，请大神不吝赐教


### 安装教程

只需要了解接口就可以做到开箱即用

### 使用说明

1、 你可以设置进度条，但默认模式最好在WindowsTerminal中或Linux环境使用，在Windows CMD中会出现显示异常，需要手动更改\
2、以下是一个简单示例：
~~~go
func main() {
	//首先创建下载器对象
	m3u8 := NewDownloader()
	//进行基本设置
	m3u8.SetUrl(TestDownloadUrl)
	m3u8.SetMovieName("小猪佩奇")
	//开始下载，使用默认下载方式，其他见文档
	if m3u8.DefaultDownload() {
		fmt.Println("下载成功")
	}
}
~~~
### M3u8Download 接口文档

##### func NewDownloader 
~~~go
func NewDownloader() M3u8Downloader
~~~
- 创建一个新下载器对象

##### func NewDownloaderWithConfig
~~~go
func NewDownloaderWithConfig(config *DownloadConfig) M3u8Downloader
~~~
- 使用自定义配置创建下载器对象\
  自定义配置为一个结构体：DownloadConfig 公有字段为配置信息，但不可滥用
  

##### func DefaultDownload

~~~go
func DefaultDownload() bool
~~~
- 默认下载方式，建议使用，方便快捷,下载为ts文件，然后合并，将会返回一个下载状态\
true表示下载并合并文件成功，反之则为失败


##### func Download

~~~go
func Download() error
~~~
- 根据之前的配置开始执行下载

##### func MergeFile

~~~go
func MergeFile() error
~~~
- 默认合并文件方法

##### func MergeFileInDir
~~~go
func MergeFileInDir(path string, saveName string) error
~~~
- 将合并后的视频文件保存到目录dir中

##### func ParseM3u8FileEncrypted
~~~go
func ParseM3u8FileEncrypted(link string) (*Result, error)
~~~
- 解析加密m3u8文件

##### func SetUrl
~~~go
func SetUrl(url string)
~~~
- 设置url

##### func SetIfShowTheBar
~~~go
func SetIfShowTheBar(ifShow bool)
~~~
- 是否显示进度条

##### func SetNumOfThread
~~~go
func SetNumOfThread(num int)
~~~
- 设置下载线程的数量

##### func SetMovieName
~~~go
func SetMovieName(videoName string)
~~~
- 设置视频的文件名

##### func SetSaveDirectory
~~~go
func SetSaveDirectory(targetDir string)
~~~
- 设置保存目录

##### func SetDownloadModel
~~~go
func SetDownloadModel(model DownloadModelType) 
~~~
- 设置下载模式:\
  WriteIntoCacheAndSaveModel 使用缓存下载模式的处理器函数,暂时废弃，因为会占用大量内存，且视频质量不高,以后更新会解决此问题\
  \
  SaveAsTsFileAndMergeModel 首先下载ts文件，并将文件合并成mp4,目前使用这种方式，下载会很稳定且下载的质量可以保证

### Bar 接口文档

##### func NewBar
~~~go
func NewBar(total int64) Bar
~~~
- 创建进度条对象，显示模式默认为在Linux Terminal下\
  默认范围:[0，total),默认显示图案：=，默认显示颜色：绿色（32）

##### func NewOptionWithGraphAndModel
~~~go
func NewOptionWithGraphAndModel(start, total int64, completedIcon rune, model ModelType) Bar
~~~
- 创建自定义范围、图案的进度条对象，并自定义显示模式\
  start, total表示起止范围，completedIcon表示已完成部分的图标、
  model 表示显示模式,共两种：\
  LinuxTerminal     适用于Windows和Linux中
  WindowsCmd 只适用于WindowsCMD中

##### func Play
~~~go
func Update(cur int64)
~~~
- 执行一次记录

##### func Finish 
~~~go
func Finish()
~~~
- 完成处理方法，在完成后会将相关数据归零，以便重新使用

#####func Setting
~~~go
func Setting() *BarConfig
~~~
- 设置信息方法，通过此方法调用其他设置方法

##### func UpdateConfig
~~~go
func UpdateConfig(newConfig *BarConfig)
~~~
- 更新配置信息对象

#####func ReSetRange
~~~go
func ReSetRange(start, total int64)
~~~ 
- 重新设置进度条范围
