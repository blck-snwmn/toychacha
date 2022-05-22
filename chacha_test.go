package gochacha

import (
	"testing"
)

func Test_rotationN(t *testing.T) {
	type args struct {
		n     uint32
		shift uint
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "rfc7539 test vector",
			args: args{0x7998bfda, 7},
			want: 0xcc5fed3c,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rotationN(tt.args.n, tt.args.shift); got != tt.want {
				t.Errorf("shift() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_quarterRound(t *testing.T) {
	type args struct {
		a uint32
		b uint32
		c uint32
		d uint32
	}
	tests := []struct {
		name  string
		args  args
		want  uint32
		want1 uint32
		want2 uint32
		want3 uint32
	}{
		{
			name: "rfc7539 test vector",
			args: args{
				0x11111111,
				0x01020304,
				0x9b8d6f43,
				0x01234567,
			},
			want:  0xea2a92f4,
			want1: 0xcb1cf8ce,
			want2: 0x4581472e,
			want3: 0x5881c4bb,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := quarterRound(tt.args.a, tt.args.b, tt.args.c, tt.args.d)
			if got != tt.want {
				t.Errorf("quarterRound() got = %x, want %x", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("quarterRound() got1 = %x, want %x", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("quarterRound() got2 = %x, want %x", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("quarterRound() got3 = %x, want %x", got3, tt.want3)
			}
		})
	}
}
