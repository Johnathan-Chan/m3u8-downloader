package M3u8Downloader

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func TestBody(t *testing.T) {
	m3u8 := NewDownloader()
	//m3u8.SetUrl("https://v.bdcache.com/vh1/640/b243d8d9c2de0a7e10bbe1d7ecb41af73ccbb03c/master.m3u8")
	//m3u8.SetMovieName("無鬼イキトランス10パート2·樱井友树")
	m3u8.SetUrl("https://v.bdcache.com/vh1/640/e72497c9caf8146269d6b4e06e0979efc50109cb/master.m3u8")
	m3u8.SetMovieName("性感女孩镶边和更热的肛门")
	//m3u8.SetUrl(TestUrl1)
	//m3u8.SetMovieName("浴血黑帮第四季第一集")
	m3u8.SetSaveDirectory("D:/Users/oopsguy/asl")
	m3u8.SetIfShowTheBar(true)
	m3u8.SetDownloadModel(WriteIntoCacheAndSaveModel)
	//m3u8.SetDownloadModel(SaveAsTsFileAndMergeModel)
	if m3u8.DefaultDownload(){
		fmt.Println("下载成功")
	}
	fmt.Println("启动新下载：")
	m3u8.SetUrl("https://v.bdcache.com/vh1/640/09bb20a3d0f48bb872d60ad30022cd3299aea4c6/master.m3u8")
	m3u8.SetMovieName("女同性恋肛门自助餐")
	if m3u8.DefaultDownload(){
		fmt.Println("下载成功")
	}
	//m3u8.ParseM3u8File()
	//body:=httpGetBodyToByte()
	//fmt.Println(string(body))
}


func TestStr(t *testing.T){
	var id64 int64 = 99
	// method 2:
	idPointer := (*int)(unsafe.Pointer(&id64))
	idd16 := *idPointer
	fmt.Println(idd16)
	fmt.Println(reflect.TypeOf(idd16))
	//res:=strconv.FormatUint(uint64(time.Now().Unix()),10)
	//fmt.Println(res)
	//fmt.Println(getUnixTimeAndToByte())
	bytes := []byte("I am byte array !")
	str := *(*string)(unsafe.Pointer(&bytes))
	bytes[0] = 'i'
	fmt.Println(str)
}

func TestNumber(t *testing.T){
	n:=processNum(520)
	fmt.Println(*(*string)(unsafe.Pointer(&n)))
}



//基于打印结果的一些时间测试
func TestUrl(t *testing.T) {
	//m3u8URL := TestUrl1
	//t0 := time.Now()
	//m3u8ParseResult, err := parseM3u8Body(m3u8URL)
	//if err != nil {
	//	panic(err)
	//}
	//t1 := time.Now()
	//fmt.Println(t1.Sub(t0))
	////for _,v:=range m3u8ParseResult.M3u8.Segments{
	////	fmt.Println(m3u8ParseResult)
	////}
	//fmt.Println(m3u8ParseResult.M3u8.Segments[66])
}

func TestTime(t *testing.T){
	//遍历打印所有的文件名
	var s []string

	s, _ = GetAllFile("D:/Users/oopsguy/m3u8_down/s", s)

	for _, v := range s {
		fmt.Println(v)
	}

}
//获取当前目录下的文件及目录信息
func pwdTest(){
	pwd, _ := os.Getwd()
	//获取文件或目录相关信息
	fileInfoList, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(fileInfoList))
	for i := range fileInfoList {
		fmt.Println(fileInfoList[i].Name()) //打印当前文件或目录下的文件或目录名
	}
}

func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)

	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}


// 递归获取指定目录下的所有文件名
func GetAllFileTime(pathname string) ([]string, error) {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败，pathname=%v, err=%v \n",pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFileTime(fullname)
			if err != nil {
				fmt.Printf("读取文件目录失败,fullname=%v, err=%v",fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}

// 把秒级的时间戳转为time格式
func SecondToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func TestFileTime(t *testing.T) {

	// 递归获取目录下的所有文件
	var files []string
	files, _ = GetAllFileTime("D:/Users/oopsguy/m3u8_down/s")

	fmt.Println("目录下的所有文件如下")
	for i:=0;i<len(files);i++ {
		fmt.Println("文件名：",files[i])

		// 获取文件原来的访问时间，修改时间
		finfo, _ := os.Stat(files[i])

		// linux环境下代码如下
		//linuxFileAttr := finfo.Sys().(*syscall.Stat_t)
		//fmt.Println("文件创建时间", SecondToTime(linuxFileAttr.Ctim.Sec))
		//fmt.Println("最后访问时间", SecondToTime(linuxFileAttr.Atim.Sec))
		//fmt.Println("最后修改时间", SecondToTime(linuxFileAttr.Mtim.Sec))

		// windows下代码如下
		winFileAttr := finfo.Sys().(*syscall.Win32FileAttributeData)
		//fmt.Println("文件创建时间：",SecondToTime(winFileAttr.CreationTime.Nanoseconds()/1e9))
		//fmt.Println("最后访问时间：",SecondToTime(winFileAttr.LastAccessTime.Nanoseconds()/1e9))
		//fmt.Println("最后修改时间：",SecondToTime(winFileAttr.LastWriteTime.Nanoseconds()/1e9))
		fmt.Println("文件创建时间：",winFileAttr.CreationTime.Nanoseconds()/1e9)
		fmt.Println("最后访问时间：",winFileAttr.LastAccessTime.Nanoseconds()/1e9)
		fmt.Println("最后修改时间：",winFileAttr.LastWriteTime.Nanoseconds()/1e9)

	}
}


func TestDirCreat(t *testing.T){
	err := os.MkdirAll("d:/fuck/", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}