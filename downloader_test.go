package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

func TestBody(t *testing.T) {
	m3u8 := NewDownloader()
	m3u8.Download()
	//m3u8.ParseM3u8File()
	//body:=simplyHttpGet()
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
	//fmt.Println(getUnixTimeByte())
	bytes := []byte("I am byte array !")
	str := *(*string)(unsafe.Pointer(&bytes))
	bytes[0] = 'i'
	fmt.Println(str)
}
