package toychacha

import (
	"math/rand/v2"
	"testing"
)

func Test_blockChacha8rand(t *testing.T) {
	{
		key := [32]byte{}
		key[0] = 0x01

		cc := rand.New(NewChaCha8(key))
		rr := rand.New(rand.NewChaCha8(key))

		for i := 0; i < 32; i++ {
			got := cc.Uint64()
			want := rr.Uint64()
			if got != want {
				t.Fatalf("invalid uint64, got=%X, want=%X", got, want)
			}
		}
	}
	{
		key := [32]byte{
			0x01, 0x02, 0x03, 0x04,
			0x05, 0x06, 0x07, 0x08,
			0x09, 0x0A, 0x0B, 0x0C,
			0x0D, 0x0E, 0x0F, 0x10,
			0x11, 0x12, 0x13, 0x14,
			0x15, 0x16, 0x17, 0x18,
			0x19, 0x1A, 0x1B, 0x1C,
		}
		cc := rand.New(NewChaCha8(key))
		rr := rand.New(rand.NewChaCha8(key))

		for i := 0; i < 32; i++ {
			got := cc.Uint64()
			want := rr.Uint64()
			if got != want {
				t.Fatalf("invalid uint64, got=%X, want=%X", got, want)
			}
		}
	}
}
