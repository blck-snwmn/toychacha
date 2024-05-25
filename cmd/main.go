package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"

	"github.com/blck-snwmn/toychacha"
	"github.com/pkg/profile"
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("time=%v\n", time.Since(now))
	}()

	b := flag.Bool("b", false, "")
	flag.Parse()
	if *b {
		defer profile.Start().Stop()
	}
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

	plaintext := []byte(`
	this is a long text. this is a pen. my name is alice. my job is a teacher.
	tell me your name. tell me your job. tell me your hobby. tell me your favorite.
	my hobby is reading. my favorite book is alice in wonderland.
	my favorite food is sushi. 

	this is a plaintext. 
	this is a plaintext.
	this is a plaintext.
	this is a plaintext.
	this is a plaintext.
	`)

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
	// {
	// 	f, err := os.Create("heap.prof")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	err = pprof.WriteHeapProfile(f)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
}
