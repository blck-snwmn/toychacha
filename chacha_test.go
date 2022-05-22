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
			name: "rfc7539 test vector",
			s: [][]uint32{
				{0x879531e0, 0xc5ecf37d, 0x516461b1, 0xc9a62f8a},
				{0x44c20ef3, 0x3390af7f, 0xd9fc690b, 0x2a5f714c},
				{0x53372767, 0xb00a5631, 0x974c541a, 0x359e9963},
				{0x5c971061, 0x3d631689, 0x2098d9d6, 0x91dbd320},
			},
			args: args{2, 7, 8, 13},
			want: [][]uint32{
				{0x879531e0, 0xc5ecf37d, 0xbdb886dc, 0xc9a62f8a},
				{0x44c20ef3, 0x3390af7f, 0xd9fc690b, 0xcfacafd2},
				{0xe46bea80, 0xb00a5631, 0x974c541a, 0x359e9963},
				{0x5c971061, 0xccc07c79, 0x2098d9d6, 0x91dbd320},
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
			want: [][]uint32{
				{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574},
				{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c},
				{0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c},
				{0x00000001, 0x09000000, 0x4a000000, 0x00000000},
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
			name: "rfc7539 test vector",
			s: [][]uint32{
				{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574},
				{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c},
				{0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c},
				{0x00000001, 0x09000000, 0x4a000000, 0x00000000},
			},
			want: [][]uint32{
				{0x837778ab, 0xe238d763, 0xa67ae21e, 0x5950bb2f},
				{0xc4f2d0c7, 0xfc62bb2f, 0x8fa018fc, 0x3f5ec7b7},
				{0x335271c2, 0xf29489f3, 0xeabda8fc, 0x82e46ebd},
				{0xd19c12b4, 0xb04e16de, 0x9e83d0cb, 0x4e3c50a2},
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
		{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574},
		{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c},
		{0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c},
		{0x00000001, 0x09000000, 0x4a000000, 0x00000000},
	}
	want := state{
		{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574},
		{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c},
		{0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c},
		{0x00000001, 0x09000000, 0x4a000000, 0x00000000},
	}
	newS := s.clone()
	if !reflect.DeepEqual(newS, want) {
		t.Errorf("clone=\n%x, want=\n%x", newS, want)
	}

	{
		newS[0][0] = 0x0
		if s[0][0] != 0x61707865 {
			t.Errorf("failed to deep copy. got=%d, want=%d", s[0][0], 0x61707865)
		}
	}
	{
		newS[1][3] = 0x0
		if s[1][3] != 0x0f0e0d0c {
			t.Errorf("failed to deep copy. got=%d, want=%d", s[1][3], 0x0f0e0d0c)
		}
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
			name: "rfc7539 test vector",
			s: state{
				{0x837778ab, 0xe238d763, 0xa67ae21e, 0x5950bb2f},
				{0xc4f2d0c7, 0xfc62bb2f, 0x8fa018fc, 0x3f5ec7b7},
				{0x335271c2, 0xf29489f3, 0xeabda8fc, 0x82e46ebd},
				{0xd19c12b4, 0xb04e16de, 0x9e83d0cb, 0x4e3c50a2},
			},
			args: args{
				state{
					{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574},
					{0x03020100, 0x07060504, 0x0b0a0908, 0x0f0e0d0c},
					{0x13121110, 0x17161514, 0x1b1a1918, 0x1f1e1d1c},
					{0x00000001, 0x09000000, 0x4a000000, 0x00000000},
				},
			},
			want: state{
				{0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3},
				{0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3},
				{0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9},
				{0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2},
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
			name: "rfc7539 test vector",
			s: state{
				{0xe4e7f110, 0x15593bd1, 0x1fdd0f50, 0xc47120a3},
				{0xc7f4d1c7, 0x0368c033, 0x9aaa2204, 0x4e6cd4c3},
				{0x466482d2, 0x09aa9f07, 0x05d7c214, 0xa2028bd9},
				{0xd19c12b5, 0xb94e16de, 0xe883d0cb, 0x4e3c50a2},
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
			name: "rfc7539 test vector",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := block(tt.args.key, tt.args.nonce, tt.args.counter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("block() = \n%x, want \n%x", got, tt.want)
			}
		})
	}
}
