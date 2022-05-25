package gochacha

import (
	"reflect"
	"testing"
)

func Test_mac(t *testing.T) {
	type args struct {
		msg []byte
		key []byte
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
					0x85, 0xd6, 0xbe, 0x78,
					0x57, 0x55, 0x6d, 0x33,
					0x7f, 0x44, 0x52, 0xfe,
					0x42, 0xd5, 0x06, 0xa8,
					0x01, 0x03, 0x80, 0x8a,
					0xfb, 0x0d, 0xb2, 0xfd,
					0x4a, 0xbf, 0xf6, 0xaf,
					0x41, 0x49, 0xf5, 0x1b,
				},
				msg: []byte{
					0x43, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x67, 0x72,
					0x61, 0x70, 0x68, 0x69, 0x63, 0x20, 0x46, 0x6f,
					0x72, 0x75, 0x6d, 0x20, 0x52, 0x65, 0x73, 0x65,
					0x61, 0x72, 0x63, 0x68, 0x20, 0x47, 0x72, 0x6f,
					0x75, 0x70,
				},
			},
			want: []byte{
				0xa8, 0x06, 0x1d, 0xc1,
				0x30, 0x51, 0x36, 0xc6,
				0xc2, 0x2b, 0x8b, 0xaf,
				0x0c, 0x01, 0x27, 0xa9,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mac(tt.args.msg, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genMacKey(t *testing.T) {
	type args struct {
		key   []byte
		nonce []byte
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
					0x80, 0x81, 0x82, 0x83,
					0x84, 0x85, 0x86, 0x87,
					0x88, 0x89, 0x8a, 0x8b,
					0x8c, 0x8d, 0x8e, 0x8f,
					0x90, 0x91, 0x92, 0x93,
					0x94, 0x95, 0x96, 0x97,
					0x98, 0x99, 0x9a, 0x9b,
					0x9c, 0x9d, 0x9e, 0x9f,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x01, 0x02, 0x03,
					0x04, 0x05, 0x06, 0x07,
				},
			},
			want: []byte{
				0x8a, 0xd5, 0xa0, 0x8b,
				0x90, 0x5f, 0x81, 0xcc,
				0x81, 0x50, 0x40, 0x27,
				0x4a, 0xb2, 0x94, 0x71,
				0xa8, 0x33, 0xb6, 0x37,
				0xe3, 0xfd, 0x0d, 0xa5,
				0x08, 0xdb, 0xb8, 0xe2,
				0xfd, 0xd1, 0xa6, 0x46,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genMacKey(tt.args.key, tt.args.nonce); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genMacKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
