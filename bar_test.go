package main

import (
	"fmt"
	"testing"
	"time"
)

func TestBar(t *testing.T) {
	var bar = NewBar(36)
	bar.Setting().SetShowModel(LinuxTerminal)
	for i:= 0; i<=36; i++{
		time.Sleep(1000*time.Millisecond)
		bar.Play(int64(i))
	}
	bar.Finish()
	//fmt.Printf("[\033[1;40;32m%s\033[0m]\n",  "testPrintColor",)
	//fmt.Printf("\033[4;31;40m%s\033[0m\n","下划线 -  红色文字，黑色底哒")
}

//█████
func TestRune(t *testing.T){
	var a StringBuilder
	for i := 0; i < 6; i++ {
		a.WriteRune('█')
	}
	for i := 6; i < 10; i++ {
		a.WriteRune('-')
	}
	fmt.Println(len(a.GetBuffer()))
	fmt.Println(a.String())
}


func TestByte(t *testing.T){
	b:=make([]byte,10)
	fmt.Printf("First Address:%p\n",&b[0])
	fmt.Printf("  End Address:%p\n",&b[9])
	fmt.Printf("Pointer[6]:%p\n",&b[5])
	fuck(b[5:8])
}


func fuck(b []byte){
	fmt.Printf("Pointer[0]:%p\n",&b[0])
	fmt.Printf("Pointer[2]:%p\n",&b[2])

}


