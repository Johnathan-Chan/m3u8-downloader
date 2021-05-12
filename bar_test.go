package M3u8Downloader

import (
	"fmt"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

// encryption
func TestBar(t *testing.T) {
	var bar = NewBar(36)
	bar.Setting().SetShowModel(LinuxTerminal)
	for i := 0; i <= 36; i++ {
		time.Sleep(100 * time.Millisecond)
		bar.Update(int64(i))
	}
	bar.Finish()
}

//█████
func TestRune(t *testing.T) {
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

func TestByte(t *testing.T) {
	b := make([]byte, 10)
	fmt.Printf("First Address:%p\n", &b[0])
	fmt.Printf("  End Address:%p\n", &b[9])
	fmt.Printf("Pointer[6]:%p\n", &b[5])
	fuck(b[5:8])
}

func fuck(b []byte) {
	fmt.Printf("Pointer[0]:%p\n", &b[0])
	fmt.Printf("Pointer[2]:%p\n", &b[2])

}

func TestPrint(t *testing.T) {
	var symbol1 rune = '█'
	var symbol2 byte = '-'
	var l, r, i int
	body := make([]byte, 200)
	for i = 0; i < 10; i++ {
		body[i] = symbol2
	}
	fmt.Println(numOfByte(body, '-'))
	l = 0
	r = 10
	n := utf8.EncodeRune(body[l:l+utf8.UTFMax], symbol1) - 1
	fmt.Println(numOfByte(body, '-'))
	//b.buf = b.buf[:l+n]
	excursion(body, '-', r, r+n)
	r += n
	fmt.Println(numOfByte(body, '-'))
	fmt.Println(*(*string)(unsafe.Pointer(&body)), "Hello")
	fmt.Println(int(symbol1))
	b1 := body[:r]
	//fmt.Printf("Pointer1:%p\nPointer2%p\n",&body[0],&b1[0])
	//fmt.Printf("Pointer1:%p\nPointer2%p\n",&body,&b1)
	//var ia *int
	//fmt.Println("body:",unsafe.Sizeof(ia)," b1:",unsafe.Sizeof(b1))
	fmt.Println(*(*string)(unsafe.Pointer(&b1)), "Hello")
}

func TestNewBar(t *testing.T) {
	var symbol1 rune = '█'
	var symbol2 byte = '-'
	var l, r, i int
	body := make([]byte, 151)
	for i = 0; i < 50; i++ {
		body[i] = symbol2
	}
	l = 0
	r = 50
	var b1 []byte

	for i = 0; i < 50; i++ {
		n := utf8.EncodeRune(body[l:l+utf8.UTFMax], symbol1)
		//fmt.Println(l,r,r+n-1)
		excursion(body, '-', r, r+n-1)
		l += n
		r += n - 1
		b1 = body[:r]
		fmt.Printf("\r[%s] %3.2f%%", *(*string)(unsafe.Pointer(&b1)), float32((float32(i+1)/50.0)*100.0))
		time.Sleep(200 * time.Millisecond)
	}
	fmt.Println()
}

func excursion(buffer []byte, symbol byte, start, end int) {
	for i := start; i < end; i++ {
		if int(buffer[i]) < 100 {
			buffer[i] = symbol
		}
	}
}

func numOfByte(buffer []byte, sym byte) int {
	var count int
	for i := 0; i < len(buffer); i++ {
		if buffer[i] == sym {
			count++
		}
	}
	return count
}
