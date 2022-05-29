package gochacha

import (
	"reflect"
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
			name: "rfc8439 test vector",
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
			name: "rfc8439 test vector",
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

func Test_state_quarterRound(t *testing.T) {
	type args struct {
		x uint
		y uint
		z uint
		w uint
	}
	tests := []struct {
		name string
		s    state
		args args
		want state
	}{
		{
			name: "rfc8439 test vector",
			s: []uint32{
				0x879531e0, 0xc5ecf37d, 0x516461b1, 0xc9a62f8a,
				0x44c20ef3, 0x3390af7f, 0xd9fc690b, 0x2a5f714c,
				0x53372767, 0xb00a5631, 0x974c541a, 0x359e9963,
				0x5c971061, 0x3d631689, 0x2098d9d6, 0x91dbd320,
			},
			args: args{2, 7, 8, 13},
			want: []uint32{
				0x879531e0, 0xc5ecf37d, 0xbdb886dc, 0xc9a62f8a,
				0x44c20ef3, 0x3390af7f, 0xd9fc690b, 0xcfacafd2,
				0xe46bea80, 0xb00a5631, 0x974c541a, 0x359e9963,
				0x5c971061, 0xccc07c79, 0x2098d9d6, 0x91dbd320,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.quarterRound(tt.args.x, tt.args.y, tt.args.z, tt.args.w)
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("quarterRound() got = %x, want %x", tt.s, tt.want)
			}
		})
	}
}

func TestNewState(t *testing.T) {
	type args struct {
		key     []byte
		nonce   []byte
		counter uint32
	}
	tests := []struct {
		name    string
		args    args
		want    state
		wantErr bool
	}{
		{
			name: "new",
			args: args{
				key: []byte{
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b,
					0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
					0x18, 0x19, 0x1a, 0x1b,
					0x1c, 0x1d, 0x1e, 0x1f,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x09,
					0x00, 0x00, 0x00, 0x4a,
					0x00, 0x00, 0x00, 0x00,
				},
				counter: 1,
			},
			want: []uint32{
				0x61707865, 0x3320646e, 0x79622d32, 0x6b206574,
				0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c,
				0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c,
				0x00000001, 0x09000000, 0x4a000000, 0x00000000,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newState(tt.args.key, tt.args.nonce, tt.args.counter)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewState() = \n%x, want \n%x", got, tt.want)
			}
		})
	}
}

func Test_state_innerBlock(t *testing.T) {
	tests := []struct {
		name string
		s    state
		want state
	}{
		{
			name: "rfc8439 test vector",
			s: []uint32{
				0x61707865, 0x3320646e, 0x79622d32, 0x6b206574,
				0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c,
				0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c,
				0x00000001, 0x09000000, 0x4a000000, 0x00000000,
			},
			want: []uint32{
				0x837778ab, 0xe238d763, 0xa67ae21e, 0x5950bb2f,
				0xc4f2d0c7, 0xfc62bb2f, 0x8fa018fc, 0x3f5ec7b7,
				0x335271c2, 0xf29489f3, 0xeabda8fc, 0x82e46ebd,
				0xd19c12b4, 0xb04e16de, 0x9e83d0cb, 0x4e3c50a2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 10; i++ {
				tt.s.innerBlock()
			}
			if !reflect.DeepEqual(tt.s, tt.want) {
				t.Errorf("innerBlock=\n%x, want=\n%x", tt.s, tt.want)
			}
		})
	}
}

func Test_state_clone(t *testing.T) {
	s := state{
		0x61707865, 0x3320646e, 0x79622d32, 0x6b206574,
		0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c,
		0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c,
		0x00000001, 0x09000000, 0x4a000000, 0x00000000,
	}
	want := state{
		0x61707865, 0x3320646e, 0x79622d32, 0x6b206574,
		0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c,
		0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c,
		0x00000001, 0x09000000, 0x4a000000, 0x00000000,
	}
	newS := s.clone()
	if !reflect.DeepEqual(newS, want) {
		t.Errorf("clone=\n%x, want=\n%x", newS, want)
	}
}

func Test_state_add(t *testing.T) {
	type args struct {
		other state
	}
	tests := []struct {
		name string
		s    state
		args args
		want state
	}{
		{
			name: "rfc8439 test vector",
			s: state{
				0x837778ab, 0xe238d763, 0xa67ae21e, 0x5950bb2f,
				0xc4f2d0c7, 0xfc62bb2f, 0x8fa018fc, 0x3f5ec7b7,
				0x335271c2, 0xf29489f3, 0xeabda8fc, 0x82e46ebd,
				0xd19c12b4, 0xb04e16de, 0x9e83d0cb, 0x4e3c50a2,
			},
			args: args{
				state{
					0x61707865, 0x3320646e, 0x79622d32, 0x6b206574,
					0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c,
					0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c,
					0x00000001, 0x09000000, 0x4a000000, 0x00000000,
				},
			},
			want: state{
				0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3,
				0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3,
				0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9,
				0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.add(tt.args.other)
		})
	}
}

func Test_state_serialize(t *testing.T) {
	tests := []struct {
		name string
		s    state
		want []byte
	}{
		{
			name: "rfc8439 test vector",
			s: state{
				0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3,
				0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3,
				0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9,
				0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2,
			},
			want: []byte{
				0x10, 0xf1, 0xe7, 0xe4, 0xd1, 0x3b, 0x59, 0x15, 0x50, 0x0f, 0xdd, 0x1f, 0xa3, 0x20, 0x71, 0xc4, //.....;Y.P.... q.
				0xc7, 0xd1, 0xf4, 0xc7, 0x33, 0xc0, 0x68, 0x03, 0x04, 0x22, 0xaa, 0x9a, 0xc3, 0xd4, 0x6c, 0x4e, //....3.h.."....lN
				0xd2, 0x82, 0x64, 0x46, 0x07, 0x9f, 0xaa, 0x09, 0x14, 0xc2, 0xd7, 0x05, 0xd9, 0x8b, 0x02, 0xa2, //..dF............
				0xb5, 0x12, 0x9c, 0xd1, 0xde, 0x16, 0x4e, 0xb9, 0xcb, 0xd0, 0x83, 0xe8, 0xa2, 0x50, 0x3c, 0x4e, //......N......P<N
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("state.serialize() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_block(t *testing.T) {
	type args struct {
		key     []byte
		nonce   []byte
		counter uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "rfc8439 test vector",
			args: args{
				key: []byte{
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b,
					0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
					0x18, 0x19, 0x1a, 0x1b,
					0x1c, 0x1d, 0x1e, 0x1f,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x09,
					0x00, 0x00, 0x00, 0x4a,
					0x00, 0x00, 0x00, 0x00,
				},
				counter: 1,
			},
			want: []byte{
				0x10, 0xf1, 0xe7, 0xe4, 0xd1, 0x3b, 0x59, 0x15, 0x50, 0x0f, 0xdd, 0x1f, 0xa3, 0x20, 0x71, 0xc4, //.....;Y.P.... q.
				0xc7, 0xd1, 0xf4, 0xc7, 0x33, 0xc0, 0x68, 0x03, 0x04, 0x22, 0xaa, 0x9a, 0xc3, 0xd4, 0x6c, 0x4e, //....3.h.."....lN
				0xd2, 0x82, 0x64, 0x46, 0x07, 0x9f, 0xaa, 0x09, 0x14, 0xc2, 0xd7, 0x05, 0xd9, 0x8b, 0x02, 0xa2, //..dF............
				0xb5, 0x12, 0x9c, 0xd1, 0xde, 0x16, 0x4e, 0xb9, 0xcb, 0xd0, 0x83, 0xe8, 0xa2, 0x50, 0x3c, 0x4e, //......N......P<N
			},
		},
		{
			name: "rfc8439 test vector#1",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				counter: 0,
			},
			want: []byte{
				0x76, 0xb8, 0xe0, 0xad, 0xa0, 0xf1, 0x3d, 0x90,
				0x40, 0x5d, 0x6a, 0xe5, 0x53, 0x86, 0xbd, 0x28,
				0xbd, 0xd2, 0x19, 0xb8, 0xa0, 0x8d, 0xed, 0x1a,
				0xa8, 0x36, 0xef, 0xcc, 0x8b, 0x77, 0x0d, 0xc7,
				0xda, 0x41, 0x59, 0x7c, 0x51, 0x57, 0x48, 0x8d,
				0x77, 0x24, 0xe0, 0x3f, 0xb8, 0xd8, 0x4a, 0x37,
				0x6a, 0x43, 0xb8, 0xf4, 0x15, 0x18, 0xa1, 0x1c,
				0xc3, 0x87, 0xb6, 0x69, 0xb2, 0xee, 0x65, 0x86,
			},
		},
		{
			name: "rfc8439 test vector#2",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				counter: 1,
			},
			want: []byte{
				0x9f, 0x07, 0xe7, 0xbe, 0x55, 0x51, 0x38, 0x7a,
				0x98, 0xba, 0x97, 0x7c, 0x73, 0x2d, 0x08, 0x0d,
				0xcb, 0x0f, 0x29, 0xa0, 0x48, 0xe3, 0x65, 0x69,
				0x12, 0xc6, 0x53, 0x3e, 0x32, 0xee, 0x7a, 0xed,
				0x29, 0xb7, 0x21, 0x76, 0x9c, 0xe6, 0x4e, 0x43,
				0xd5, 0x71, 0x33, 0xb0, 0x74, 0xd8, 0x39, 0xd5,
				0x31, 0xed, 0x1f, 0x28, 0x51, 0x0a, 0xfb, 0x45,
				0xac, 0xe1, 0x0a, 0x1f, 0x4b, 0x79, 0x4d, 0x6f,
			},
		},
		{
			name: "rfc8439 test vector#3",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				counter: 1,
			},
			want: []byte{
				0x3a, 0xeb, 0x52, 0x24, 0xec, 0xf8, 0x49, 0x92,
				0x9b, 0x9d, 0x82, 0x8d, 0xb1, 0xce, 0xd4, 0xdd,
				0x83, 0x20, 0x25, 0xe8, 0x01, 0x8b, 0x81, 0x60,
				0xb8, 0x22, 0x84, 0xf3, 0xc9, 0x49, 0xaa, 0x5a,
				0x8e, 0xca, 0x00, 0xbb, 0xb4, 0xa7, 0x3b, 0xda,
				0xd1, 0x92, 0xb5, 0xc4, 0x2f, 0x73, 0xf2, 0xfd,
				0x4e, 0x27, 0x36, 0x44, 0xc8, 0xb3, 0x61, 0x25,
				0xa6, 0x4a, 0xdd, 0xeb, 0x00, 0x6c, 0x13, 0xa0,
			},
		},
		{
			name: "rfc8439 test vector#4",
			args: args{
				key: []byte{
					0x00, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				counter: 2,
			},
			want: []byte{
				0x72, 0xd5, 0x4d, 0xfb, 0xf1, 0x2e, 0xc4, 0x4b,
				0x36, 0x26, 0x92, 0xdf, 0x94, 0x13, 0x7f, 0x32,
				0x8f, 0xea, 0x8d, 0xa7, 0x39, 0x90, 0x26, 0x5e,
				0xc1, 0xbb, 0xbe, 0xa1, 0xae, 0x9a, 0xf0, 0xca,
				0x13, 0xb2, 0x5a, 0xa2, 0x6c, 0xb4, 0xa6, 0x48,
				0xcb, 0x9b, 0x9d, 0x1b, 0xe6, 0x5b, 0x2c, 0x09,
				0x24, 0xa6, 0x6c, 0x54, 0xd5, 0x45, 0xec, 0x1b,
				0x73, 0x74, 0xf4, 0x87, 0x2e, 0x99, 0xf0, 0x96,
			},
		},
		{
			name: "rfc8439 test vector#5",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				},
				counter: 0,
			},
			want: []byte{
				0xc2, 0xc6, 0x4d, 0x37, 0x8c, 0xd5, 0x36, 0x37,
				0x4a, 0xe2, 0x04, 0xb9, 0xef, 0x93, 0x3f, 0xcd,
				0x1a, 0x8b, 0x22, 0x88, 0xb3, 0xdf, 0xa4, 0x96,
				0x72, 0xab, 0x76, 0x5b, 0x54, 0xee, 0x27, 0xc7,
				0x8a, 0x97, 0x0e, 0x0e, 0x95, 0x5c, 0x14, 0xf3,
				0xa8, 0x8e, 0x74, 0x1b, 0x97, 0xc2, 0x86, 0xf7,
				0x5f, 0x8f, 0xc2, 0x99, 0xe8, 0x14, 0x83, 0x62,
				0xfa, 0x19, 0x8a, 0x39, 0x53, 0x1b, 0xed, 0x6d,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := block(tt.args.key, tt.args.nonce, tt.args.counter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("block() = \n%x, want \n%x", got, tt.want)
			}
		})
	}
}

func Test_encrypt(t *testing.T) {
	type args struct {
		key       []byte
		nonce     []byte
		plaintext []byte
		counter   uint32
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "rfc8439 test vector",
			args: args{
				key: []byte{
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
					0x08, 0x09, 0x0a, 0x0b,
					0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17,
					0x18, 0x19, 0x1a, 0x1b,
					0x1c, 0x1d, 0x1e, 0x1f,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x4a,
					0x00, 0x00, 0x00, 0x00,
				},
				plaintext: []byte{
					0x4c, 0x61, 0x64, 0x69, 0x65, 0x73, 0x20, 0x61,
					0x6e, 0x64, 0x20, 0x47, 0x65, 0x6e, 0x74, 0x6c,
					0x65, 0x6d, 0x65, 0x6e, 0x20, 0x6f, 0x66, 0x20,
					0x74, 0x68, 0x65, 0x20, 0x63, 0x6c, 0x61, 0x73,
					0x73, 0x20, 0x6f, 0x66, 0x20, 0x27, 0x39, 0x39,
					0x3a, 0x20, 0x49, 0x66, 0x20, 0x49, 0x20, 0x63,
					0x6f, 0x75, 0x6c, 0x64, 0x20, 0x6f, 0x66, 0x66,
					0x65, 0x72, 0x20, 0x79, 0x6f, 0x75, 0x20, 0x6f,
					0x6e, 0x6c, 0x79, 0x20, 0x6f, 0x6e, 0x65, 0x20,
					0x74, 0x69, 0x70, 0x20, 0x66, 0x6f, 0x72, 0x20,
					0x74, 0x68, 0x65, 0x20, 0x66, 0x75, 0x74, 0x75,
					0x72, 0x65, 0x2c, 0x20, 0x73, 0x75, 0x6e, 0x73,
					0x63, 0x72, 0x65, 0x65, 0x6e, 0x20, 0x77, 0x6f,
					0x75, 0x6c, 0x64, 0x20, 0x62, 0x65, 0x20, 0x69,
					0x74, 0x2e,
				},
				counter: 1,
			},
			want: []byte{
				0x6e, 0x2e, 0x35, 0x9a, 0x25, 0x68, 0xf9, 0x80,
				0x41, 0xba, 0x07, 0x28, 0xdd, 0x0d, 0x69, 0x81,
				0xe9, 0x7e, 0x7a, 0xec, 0x1d, 0x43, 0x60, 0xc2,
				0x0a, 0x27, 0xaf, 0xcc, 0xfd, 0x9f, 0xae, 0x0b,
				0xf9, 0x1b, 0x65, 0xc5, 0x52, 0x47, 0x33, 0xab,
				0x8f, 0x59, 0x3d, 0xab, 0xcd, 0x62, 0xb3, 0x57,
				0x16, 0x39, 0xd6, 0x24, 0xe6, 0x51, 0x52, 0xab,
				0x8f, 0x53, 0x0c, 0x35, 0x9f, 0x08, 0x61, 0xd8,
				0x07, 0xca, 0x0d, 0xbf, 0x50, 0x0d, 0x6a, 0x61,
				0x56, 0xa3, 0x8e, 0x08, 0x8a, 0x22, 0xb6, 0x5e,
				0x52, 0xbc, 0x51, 0x4d, 0x16, 0xcc, 0xf8, 0x06,
				0x81, 0x8c, 0xe9, 0x1a, 0xb7, 0x79, 0x37, 0x36,
				0x5a, 0xf9, 0x0b, 0xbf, 0x74, 0xa3, 0x5b, 0xe6,
				0xb4, 0x0b, 0x8e, 0xed, 0xf2, 0x78, 0x5e, 0x42,
				0x87, 0x4d,
			},
		},
		{
			name: "rfc8439 test vector#1",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
				},
				plaintext: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				counter: 0,
			},
			want: []byte{
				0x76, 0xb8, 0xe0, 0xad, 0xa0, 0xf1, 0x3d, 0x90, 0x40, 0x5d, 0x6a, 0xe5, 0x53, 0x86, 0xbd, 0x28,
				0xbd, 0xd2, 0x19, 0xb8, 0xa0, 0x8d, 0xed, 0x1a, 0xa8, 0x36, 0xef, 0xcc, 0x8b, 0x77, 0x0d, 0xc7,
				0xda, 0x41, 0x59, 0x7c, 0x51, 0x57, 0x48, 0x8d, 0x77, 0x24, 0xe0, 0x3f, 0xb8, 0xd8, 0x4a, 0x37,
				0x6a, 0x43, 0xb8, 0xf4, 0x15, 0x18, 0xa1, 0x1c, 0xc3, 0x87, 0xb6, 0x69, 0xb2, 0xee, 0x65, 0x86,
			},
		},
		{
			name: "rfc8439 test vector#2",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x02,
				},
				plaintext: []byte{
					0x41, 0x6e, 0x79, 0x20, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x20, 0x74,
					0x6f, 0x20, 0x74, 0x68, 0x65, 0x20, 0x49, 0x45, 0x54, 0x46, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x6e,
					0x64, 0x65, 0x64, 0x20, 0x62, 0x79, 0x20, 0x74, 0x68, 0x65, 0x20, 0x43, 0x6f, 0x6e, 0x74, 0x72,
					0x69, 0x62, 0x75, 0x74, 0x6f, 0x72, 0x20, 0x66, 0x6f, 0x72, 0x20, 0x70, 0x75, 0x62, 0x6c, 0x69,
					0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x20, 0x61, 0x73, 0x20, 0x61, 0x6c, 0x6c, 0x20, 0x6f, 0x72,
					0x20, 0x70, 0x61, 0x72, 0x74, 0x20, 0x6f, 0x66, 0x20, 0x61, 0x6e, 0x20, 0x49, 0x45, 0x54, 0x46,
					0x20, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x2d, 0x44, 0x72, 0x61, 0x66, 0x74, 0x20,
					0x6f, 0x72, 0x20, 0x52, 0x46, 0x43, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x61, 0x6e, 0x79, 0x20, 0x73,
					0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x20, 0x6d, 0x61, 0x64, 0x65, 0x20, 0x77, 0x69,
					0x74, 0x68, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74,
					0x20, 0x6f, 0x66, 0x20, 0x61, 0x6e, 0x20, 0x49, 0x45, 0x54, 0x46, 0x20, 0x61, 0x63, 0x74, 0x69,
					0x76, 0x69, 0x74, 0x79, 0x20, 0x69, 0x73, 0x20, 0x63, 0x6f, 0x6e, 0x73, 0x69, 0x64, 0x65, 0x72,
					0x65, 0x64, 0x20, 0x61, 0x6e, 0x20, 0x22, 0x49, 0x45, 0x54, 0x46, 0x20, 0x43, 0x6f, 0x6e, 0x74,
					0x72, 0x69, 0x62, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x2e, 0x20, 0x53, 0x75, 0x63, 0x68, 0x20,
					0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x20, 0x69, 0x6e, 0x63, 0x6c, 0x75,
					0x64, 0x65, 0x20, 0x6f, 0x72, 0x61, 0x6c, 0x20, 0x73, 0x74, 0x61, 0x74, 0x65, 0x6d, 0x65, 0x6e,
					0x74, 0x73, 0x20, 0x69, 0x6e, 0x20, 0x49, 0x45, 0x54, 0x46, 0x20, 0x73, 0x65, 0x73, 0x73, 0x69,
					0x6f, 0x6e, 0x73, 0x2c, 0x20, 0x61, 0x73, 0x20, 0x77, 0x65, 0x6c, 0x6c, 0x20, 0x61, 0x73, 0x20,
					0x77, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x20, 0x61, 0x6e, 0x64, 0x20, 0x65, 0x6c, 0x65, 0x63,
					0x74, 0x72, 0x6f, 0x6e, 0x69, 0x63, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61,
					0x74, 0x69, 0x6f, 0x6e, 0x73, 0x20, 0x6d, 0x61, 0x64, 0x65, 0x20, 0x61, 0x74, 0x20, 0x61, 0x6e,
					0x79, 0x20, 0x74, 0x69, 0x6d, 0x65, 0x20, 0x6f, 0x72, 0x20, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x2c,
					0x20, 0x77, 0x68, 0x69, 0x63, 0x68, 0x20, 0x61, 0x72, 0x65, 0x20, 0x61, 0x64, 0x64, 0x72, 0x65,
					0x73, 0x73, 0x65, 0x64, 0x20, 0x74, 0x6f,
				},
				counter: 1,
			},
			want: []byte{
				0xa3, 0xfb, 0xf0, 0x7d, 0xf3, 0xfa, 0x2f, 0xde, 0x4f, 0x37, 0x6c, 0xa2, 0x3e, 0x82, 0x73, 0x70,
				0x41, 0x60, 0x5d, 0x9f, 0x4f, 0x4f, 0x57, 0xbd, 0x8c, 0xff, 0x2c, 0x1d, 0x4b, 0x79, 0x55, 0xec,
				0x2a, 0x97, 0x94, 0x8b, 0xd3, 0x72, 0x29, 0x15, 0xc8, 0xf3, 0xd3, 0x37, 0xf7, 0xd3, 0x70, 0x05,
				0x0e, 0x9e, 0x96, 0xd6, 0x47, 0xb7, 0xc3, 0x9f, 0x56, 0xe0, 0x31, 0xca, 0x5e, 0xb6, 0x25, 0x0d,
				0x40, 0x42, 0xe0, 0x27, 0x85, 0xec, 0xec, 0xfa, 0x4b, 0x4b, 0xb5, 0xe8, 0xea, 0xd0, 0x44, 0x0e,
				0x20, 0xb6, 0xe8, 0xdb, 0x09, 0xd8, 0x81, 0xa7, 0xc6, 0x13, 0x2f, 0x42, 0x0e, 0x52, 0x79, 0x50,
				0x42, 0xbd, 0xfa, 0x77, 0x73, 0xd8, 0xa9, 0x05, 0x14, 0x47, 0xb3, 0x29, 0x1c, 0xe1, 0x41, 0x1c,
				0x68, 0x04, 0x65, 0x55, 0x2a, 0xa6, 0xc4, 0x05, 0xb7, 0x76, 0x4d, 0x5e, 0x87, 0xbe, 0xa8, 0x5a,
				0xd0, 0x0f, 0x84, 0x49, 0xed, 0x8f, 0x72, 0xd0, 0xd6, 0x62, 0xab, 0x05, 0x26, 0x91, 0xca, 0x66,
				0x42, 0x4b, 0xc8, 0x6d, 0x2d, 0xf8, 0x0e, 0xa4, 0x1f, 0x43, 0xab, 0xf9, 0x37, 0xd3, 0x25, 0x9d,
				0xc4, 0xb2, 0xd0, 0xdf, 0xb4, 0x8a, 0x6c, 0x91, 0x39, 0xdd, 0xd7, 0xf7, 0x69, 0x66, 0xe9, 0x28,
				0xe6, 0x35, 0x55, 0x3b, 0xa7, 0x6c, 0x5c, 0x87, 0x9d, 0x7b, 0x35, 0xd4, 0x9e, 0xb2, 0xe6, 0x2b,
				0x08, 0x71, 0xcd, 0xac, 0x63, 0x89, 0x39, 0xe2, 0x5e, 0x8a, 0x1e, 0x0e, 0xf9, 0xd5, 0x28, 0x0f,
				0xa8, 0xca, 0x32, 0x8b, 0x35, 0x1c, 0x3c, 0x76, 0x59, 0x89, 0xcb, 0xcf, 0x3d, 0xaa, 0x8b, 0x6c,
				0xcc, 0x3a, 0xaf, 0x9f, 0x39, 0x79, 0xc9, 0x2b, 0x37, 0x20, 0xfc, 0x88, 0xdc, 0x95, 0xed, 0x84,
				0xa1, 0xbe, 0x05, 0x9c, 0x64, 0x99, 0xb9, 0xfd, 0xa2, 0x36, 0xe7, 0xe8, 0x18, 0xb0, 0x4b, 0x0b,
				0xc3, 0x9c, 0x1e, 0x87, 0x6b, 0x19, 0x3b, 0xfe, 0x55, 0x69, 0x75, 0x3f, 0x88, 0x12, 0x8c, 0xc0,
				0x8a, 0xaa, 0x9b, 0x63, 0xd1, 0xa1, 0x6f, 0x80, 0xef, 0x25, 0x54, 0xd7, 0x18, 0x9c, 0x41, 0x1f,
				0x58, 0x69, 0xca, 0x52, 0xc5, 0xb8, 0x3f, 0xa3, 0x6f, 0xf2, 0x16, 0xb9, 0xc1, 0xd3, 0x00, 0x62,
				0xbe, 0xbc, 0xfd, 0x2d, 0xc5, 0xbc, 0xe0, 0x91, 0x19, 0x34, 0xfd, 0xa7, 0x9a, 0x86, 0xf6, 0xe6,
				0x98, 0xce, 0xd7, 0x59, 0xc3, 0xff, 0x9b, 0x64, 0x77, 0x33, 0x8f, 0x3d, 0xa4, 0xf9, 0xcd, 0x85,
				0x14, 0xea, 0x99, 0x82, 0xcc, 0xaf, 0xb3, 0x41, 0xb2, 0x38, 0x4d, 0xd9, 0x02, 0xf3, 0xd1, 0xab,
				0x7a, 0xc6, 0x1d, 0xd2, 0x9c, 0x6f, 0x21, 0xba, 0x5b, 0x86, 0x2f, 0x37, 0x30, 0xe3, 0x7c, 0xfd,
				0xc4, 0xfd, 0x80, 0x6c, 0x22, 0xf2, 0x21,
			},
		},
		{
			name: "rfc8439 test vector#3",
			args: args{
				key: []byte{
					0x1c, 0x92, 0x40, 0xa5, 0xeb, 0x55, 0xd3, 0x8a, 0xf3, 0x33, 0x88, 0x86, 0x04, 0xf6, 0xb5, 0xf0,
					0x47, 0x39, 0x17, 0xc1, 0x40, 0x2b, 0x80, 0x09, 0x9d, 0xca, 0x5c, 0xbc, 0x20, 0x70, 0x75, 0xc0,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x02,
				},
				plaintext: []byte{
					0x27, 0x54, 0x77, 0x61, 0x73, 0x20, 0x62, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x67, 0x2c, 0x20, 0x61,
					0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x6c, 0x69, 0x74, 0x68, 0x79, 0x20, 0x74, 0x6f,
					0x76, 0x65, 0x73, 0x0a, 0x44, 0x69, 0x64, 0x20, 0x67, 0x79, 0x72, 0x65, 0x20, 0x61, 0x6e, 0x64,
					0x20, 0x67, 0x69, 0x6d, 0x62, 0x6c, 0x65, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x77,
					0x61, 0x62, 0x65, 0x3a, 0x0a, 0x41, 0x6c, 0x6c, 0x20, 0x6d, 0x69, 0x6d, 0x73, 0x79, 0x20, 0x77,
					0x65, 0x72, 0x65, 0x20, 0x74, 0x68, 0x65, 0x20, 0x62, 0x6f, 0x72, 0x6f, 0x67, 0x6f, 0x76, 0x65,
					0x73, 0x2c, 0x0a, 0x41, 0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x6f, 0x6d, 0x65, 0x20,
					0x72, 0x61, 0x74, 0x68, 0x73, 0x20, 0x6f, 0x75, 0x74, 0x67, 0x72, 0x61, 0x62, 0x65, 0x2e,
				},
				counter: 42,
			},
			want: []byte{
				0x62, 0xe6, 0x34, 0x7f, 0x95, 0xed, 0x87, 0xa4, 0x5f, 0xfa, 0xe7, 0x42, 0x6f, 0x27, 0xa1, 0xdf,
				0x5f, 0xb6, 0x91, 0x10, 0x04, 0x4c, 0x0d, 0x73, 0x11, 0x8e, 0xff, 0xa9, 0x5b, 0x01, 0xe5, 0xcf,
				0x16, 0x6d, 0x3d, 0xf2, 0xd7, 0x21, 0xca, 0xf9, 0xb2, 0x1e, 0x5f, 0xb1, 0x4c, 0x61, 0x68, 0x71,
				0xfd, 0x84, 0xc5, 0x4f, 0x9d, 0x65, 0xb2, 0x83, 0x19, 0x6c, 0x7f, 0xe4, 0xf6, 0x05, 0x53, 0xeb,
				0xf3, 0x9c, 0x64, 0x02, 0xc4, 0x22, 0x34, 0xe3, 0x2a, 0x35, 0x6b, 0x3e, 0x76, 0x43, 0x12, 0xa6,
				0x1a, 0x55, 0x32, 0x05, 0x57, 0x16, 0xea, 0xd6, 0x96, 0x25, 0x68, 0xf8, 0x7d, 0x3f, 0x3f, 0x77,
				0x04, 0xc6, 0xa8, 0xd1, 0xbc, 0xd1, 0xbf, 0x4d, 0x50, 0xd6, 0x15, 0x4b, 0x6d, 0xa7, 0x31, 0xb1,
				0x87, 0xb5, 0x8d, 0xfd, 0x72, 0x8a, 0xfa, 0x36, 0x75, 0x7a, 0x79, 0x7a, 0xc1, 0x88, 0xd1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encrypt(tt.args.key, tt.args.nonce, tt.args.plaintext, tt.args.counter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
