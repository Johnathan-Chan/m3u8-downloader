package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
	"unsafe"
)

func processNum(n int)[]byte{
	if n < 10 {
		return []byte{48,48,48, byte(48 + n),'.','t','s'}
	}else if n <100 {
		return []byte{48,48, byte(48+(n / 10)), byte(48+(n % 10)),'.','t','s'}
	}else if n < 1000 {
		return []byte{48,byte(48+(n/100)),byte(48+((n/10)%10)), byte(48+(n%10)),'.','t','s'}
	}else {
		return []byte{byte(48+(n/1000)),byte(48+((n/100)%10)),byte(48+((n/10)%10)), byte(48+(n%10)),'.','t','s'}
	}
}

// getAllNonDirectoryFile 获取所有非目录文件
func getAllNonDirectoryFile(pathName string)([]string,error){
	rd, err := ioutil.ReadDir(pathName)
	if err != nil {
		return nil, errorMap[ReadDirectoryException]
	}
	Files:=make([]string,0)
	for i:=0;i<len(rd);i++ {
		if !rd[i].IsDir() {
			fullName := pathName + "/" + rd[i].Name()
			Files = append(Files, fullName)
		}
	}
	rd = nil
	return Files, nil
}


func httpGet(url string) (io.ReadCloser, DownloadExceptionType) {
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		if err.Error()[len(err.Error())-12:] == "no such host" {
			return nil,NetworkException
		}else{
			errorMap[UnexpectedException] = err
		}
		return nil, UnexpectedException
	}
	if resp.StatusCode != 200 {
		errorMap[HttpException] = fmt.Errorf("[HttpError]:status code %d", resp.StatusCode)
		return nil, HttpException
	}
	return resp.Body, NoException
}



// parseLines 按行来解析m3u8文件
func parseLines(lines []string) (*M3u8, error) {
	var (
		i       = 0
		lineLen = len(lines)
		m3u8    = &M3u8{}
		key *Key
		seg *Segment
	)
	for ; i < lineLen; i++ {
		//TrimSpace返回字符串s的一个片段，去掉Unicode定义的所有前导和尾随空格。
		line := strings.TrimSpace(lines[i])
		if i == 0 {
			if "#EXTM3U" != line {
				return nil, errorMap[InvalidM3u8Exception]
			}
			continue
		}
		switch {
		case line == "":
			continue
		case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
			i++
			m3u8.MasterPlaylistURIs = append(m3u8.MasterPlaylistURIs, lines[i])
			continue
		case !strings.HasPrefix(line, "#"):
			seg = new(Segment)
			seg.URI = line
			m3u8.Segments = append(m3u8.Segments, seg)
			seg.Key = key
			continue
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			params := parseLineParameters(line)
			if len(params) == 0 {
				return nil, errorMap[InvalidEXT_X_KEY]
			}
			key = new(Key)
			method := CryptMethod(params["METHOD"])
			if method != "" && method != CryptMethodAES && method != CryptMethodNONE {
				return nil, errorMap[InvalidEXT_X_KEYMethod]
			}
			key.Method = method
			key.URI = params["URI"]
			key.IV = params["IV"]
		default:
			continue
		}
	}
	return m3u8, nil
}

//合并文件主函数
func mergeFile(path string, fileList []string,saveName string)error{
	var (
		buffer StringBuilder
		err error
		movie *os.File
	)
	buffer.WriteString(path)
	buffer.WriteString(saveName)
	movie,err = os.OpenFile(buffer.String(),os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err!= nil{
		return err
	}
	defer movie.Close()
	var (
		tsFile *os.File
		body []byte
	)
	for i:=0;i<len(fileList);i++ {
		tsFile,err= os.OpenFile(fileList[i],os.O_CREATE|os.O_RDONLY, os.ModePerm)
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
		err = os.Remove(fileList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// getUnixTimeAndToByte 根据加当前时间戳设置为默认名称
func getUnixTimeAndToByte() []byte {
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


func parseLineParameters(line string) map[string]string {
	r := lineParameterPattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}

// ResolveURL 处理Url
func ResolveURL(u *url.URL, p string) string {
	if strings.HasPrefix(p, "https://") || strings.HasPrefix(p, "http://") {
		return p
	}
	var baseURL string
	if strings.Index(p, "/") == 0 {
		baseURL = u.Scheme + "://" + u.Host
	} else {
		tU := u.String()
		baseURL = tU[0:strings.LastIndex(tU, "/")]
	}
	return baseURL + path.Join("/", p)
}

func AES128Encrypt(origData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) == 0 {
		iv = key
	}
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AES128Decrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) == 0 {
		iv = key
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
