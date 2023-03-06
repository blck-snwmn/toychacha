package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime/pprof"

	"github.com/blck-snwmn/toychacha"
)

func main() {
	// {
	// 	f, err := os.Create("cpu.prof")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	err = pprof.StartCPUProfile(f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer pprof.StopCPUProfile()
	// }

	plaintext := make([]byte, 100000)
	_, _ = rand.Read(plaintext)

	aad := []byte{
		0x50, 0x51, 0x52, 0x53,
		0xc0, 0xc1, 0xc2, 0xc3,
		0xc4, 0xc5, 0xc6, 0xc7,
	}
	key := []byte{
		0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87,
		0x88, 0x89, 0x8a, 0x8b, 0x8c, 0x8d, 0x8e, 0x8f,
		0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97,
		0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f,
	}

	nonce := []byte{
		0x07, 0x00, 0x00, 0x00, 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47,
	}

	tcp, _ := toychacha.New(key)
	aead := make([]byte, len(plaintext)+tcp.Overhead())
	aead = tcp.Seal(aead, nonce, []byte(plaintext), aad)
	// fmt.Printf("AEAD=%X\n", aead)

	p := make([]byte, len(plaintext))
	p, _ = tcp.Open(p, nonce, aead, aad)
	fmt.Printf("seal -> open =%v\n", reflect.DeepEqual(plaintext, p))
	{
		f, err := os.Create("heap.prof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}
