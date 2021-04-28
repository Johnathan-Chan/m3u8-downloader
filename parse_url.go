package main


//func download() {
//	defer func() {
//		if r := recover(); r != nil {
//			fmt.Println("Panic:", r)
//		}
//	}()
//	m3u8URL := "https://v.bdcache.com/vh3/640/3c24f5d470b860351bc5e670aef38446c0a125dc/master.m3u8"
//	m3u8ParseResult, err := parseM3u8FileEncrypted(m3u8URL)
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