package M3u8Downloader

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"regexp"
	"strings"
)

var(
	lineParameterPattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)
)

type(
	CryptMethod string
)

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


func parseLineParameters(line string) map[string]string {
	r := lineParameterPattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
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

//func download() {
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("Panic:", r)
//		}
//	}()
//	m3u8URL := "https://v.bdcache.com/vh3/640/3c24f5d470b860351bc5e670aef38446c0a125dc/master.m3u8"
//	m3u8ParseResult, err := ParseM3u8FileEncrypted(m3u8URL)
//	if err != nil {
//		panic(err)
//	}
//	storeFolder := "/Users/oopsguy/m3u8_down/s"
//	if err := os.MkdirAll(storeFolder, 0777); err != nil {
//		panic(err)
//	}
//
//	var wg sync.WaitGroup
//	// 防止协程启动过多，限制频率
//	limitChan := make(chan byte, 20)
//	// 开启协程请求
//	for idx, seg := range m3u8ParseResult.M3u8.Segments {
//		wg.Add(1)
//		go func(i int, s *Segment) {
//			defer func() {
//				wg.Done()
//				<-limitChan
//			}()
//			// 以需要命名文件
//			fullURL := ResolveURL(m3u8ParseResult.URL, s.URI)
//			body, err := httpGet(fullURL)
//			if err != nil {
//				fmt.Printf("Download failed [%s] %s\n", err.Error(), fullURL)
//				return
//			}
//			defer body.Close()
//			// 创建存在 TS 数据的文件
//			tsFile := filepath.Join(storeFolder, strconv.Itoa(i)+".ts")
//			tsFileTmpPath := tsFile + "_tmp"
//			tsFileTmp, err := os.Create(tsFileTmpPath)
//			if err != nil {
//				fmt.Printf("Create TS file failed: %s\n", err.Error())
//				return
//			}
//			//noinspection GoUnhandledErrorResult
//			defer tsFileTmp.Close()
//			bytes, err := ioutil.ReadAll(body)
//			if err != nil {
//				fmt.Printf("Read TS file failed: %s\n", err.Error())
//				return
//			}
//			// 解密 TS 数据
//			if s.Key != nil {
//				key := m3u8ParseResult.Keys[s.Key]
//				if key != "" {
//					bytes, err = AES128Decrypt(bytes, []byte(key), []byte(s.Key.IV))
//					if err != nil {
//						fmt.Printf("decryt TS failed: %s\n", err.Error())
//					}
//				}
//			}
//			syncByte := uint8(71) //0x47
//			bLen := len(bytes)
//			for j := 0; j < bLen; j++ {
//				if bytes[j] == syncByte {
//					bytes = bytes[j:]
//					break
//				}
//			}
//			if _, err := tsFileTmp.Write(bytes); err != nil {
//				fmt.Printf("Save TS file failed:%s\n", err.Error())
//				return
//			}
//			_ = tsFileTmp.Close()
//			// 重命名为正式文件
//			if err = os.Rename(tsFileTmpPath, tsFile); err != nil {
//				fmt.Printf("Rename TS file failed: %s\n", err.Error())
//				return
//			}
//			fmt.Printf("下载成功：%s\n", fullURL)
//		}(idx, seg)
//		limitChan <- 1
//	}
//	wg.Wait()
//
//	// 按 ts 文件名顺序合并文件
//	// 由于是从 0 开始计算，只需要递增到 len(m3u8ParseResult.M3u8.Segments)-1 即可
//	mainFile, err := os.Create(filepath.Join(storeFolder, "main.ts"))
//	if err != nil {
//		panic(err)
//	}
//	//noinspection GoUnhandledErrorResult
//	defer mainFile.Close()
//	for i := 0; i < len(m3u8ParseResult.M3u8.Segments); i++ {
//		bytes, err := ioutil.ReadFile(filepath.Join(storeFolder, strconv.Itoa(i)+".ts"))
//		if err != nil {
//			fmt.Println(err.Error())
//			continue
//		}
//		if _, err := mainFile.Write(bytes); err != nil {
//			fmt.Println(err.Error())
//			continue
//		}
//	}
//	_ = mainFile.Sync()
//	fmt.Println("下载完成")
//}