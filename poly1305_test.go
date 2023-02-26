package toychacha

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
			name: "rfc8439 test vector",
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
		{
			name: "rfc8439 test vector#1",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#2",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x36, 0xe5, 0xf6, 0xb5, 0xc5, 0xe0, 0x60, 0x70,
					0xf0, 0xef, 0xca, 0x96, 0x22, 0x7a, 0x86, 0x3e,
				},
				msg: []byte{
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
			},
			want: []byte{
				0x36, 0xe5, 0xf6, 0xb5, 0xc5, 0xe0, 0x60, 0x70,
				0xf0, 0xef, 0xca, 0x96, 0x22, 0x7a, 0x86, 0x3e,
			},
		},
		{
			name: "rfc8439 test vector#3",
			args: args{
				key: []byte{
					0x36, 0xe5, 0xf6, 0xb5, 0xc5, 0xe0, 0x60, 0x70,
					0xf0, 0xef, 0xca, 0x96, 0x22, 0x7a, 0x86, 0x3e,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
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
			},
			want: []byte{
				0xf3, 0x47, 0x7e, 0x7c, 0xd9, 0x54, 0x17, 0xaf,
				0x89, 0xa6, 0xb8, 0x79, 0x4c, 0x31, 0x0c, 0xf0,
			},
		},
		{
			name: "rfc8439 test vector#4",
			args: args{
				key: []byte{
					0x1c, 0x92, 0x40, 0xa5, 0xeb, 0x55, 0xd3, 0x8a, 0xf3, 0x33, 0x88, 0x86, 0x04, 0xf6, 0xb5, 0xf0,
					0x47, 0x39, 0x17, 0xc1, 0x40, 0x2b, 0x80, 0x09, 0x9d, 0xca, 0x5c, 0xbc, 0x20, 0x70, 0x75, 0xc0,
				},
				msg: []byte{
					0x27, 0x54, 0x77, 0x61, 0x73, 0x20, 0x62, 0x72, 0x69, 0x6c, 0x6c, 0x69, 0x67, 0x2c, 0x20, 0x61,
					0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x73, 0x6c, 0x69, 0x74, 0x68, 0x79, 0x20, 0x74, 0x6f,
					0x76, 0x65, 0x73, 0x0a, 0x44, 0x69, 0x64, 0x20, 0x67, 0x79, 0x72, 0x65, 0x20, 0x61, 0x6e, 0x64,
					0x20, 0x67, 0x69, 0x6d, 0x62, 0x6c, 0x65, 0x20, 0x69, 0x6e, 0x20, 0x74, 0x68, 0x65, 0x20, 0x77,
					0x61, 0x62, 0x65, 0x3a, 0x0a, 0x41, 0x6c, 0x6c, 0x20, 0x6d, 0x69, 0x6d, 0x73, 0x79, 0x20, 0x77,
					0x65, 0x72, 0x65, 0x20, 0x74, 0x68, 0x65, 0x20, 0x62, 0x6f, 0x72, 0x6f, 0x67, 0x6f, 0x76, 0x65,
					0x73, 0x2c, 0x0a, 0x41, 0x6e, 0x64, 0x20, 0x74, 0x68, 0x65, 0x20, 0x6d, 0x6f, 0x6d, 0x65, 0x20,
					0x72, 0x61, 0x74, 0x68, 0x73, 0x20, 0x6f, 0x75, 0x74, 0x67, 0x72, 0x61, 0x62, 0x65, 0x2e,
				},
			},
			want: []byte{
				0x45, 0x41, 0x66, 0x9a, 0x7e, 0xaa, 0xee, 0x61, 0xe7, 0x08, 0xdc, 0x7c, 0xbc, 0xc5, 0xeb, 0x62,
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
				msg: []byte{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#6",
			args: args{
				key: []byte{
					0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
				msg: []byte{
					0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#7",
			args: args{
				key: []byte{
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xF0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0x11, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#8",
			args: args{
				key: []byte{
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFB, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE,
					0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE, 0xFE,
					0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
					0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
				},
			},
			want: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#9",
			args: args{
				key: []byte{
					0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0xFD, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
					0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				},
			},
			want: []byte{
				0xFA, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
			},
		},
		{
			name: "rfc8439 test vector#10",
			args: args{
				key: []byte{
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0xE3, 0x35, 0x94, 0xD7, 0x50, 0x5E, 0x43, 0xB9,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x33, 0x94, 0xD7, 0x50, 0x5E, 0x43, 0x79, 0xCD,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x55, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "rfc8439 test vector#11",
			args: args{
				key: []byte{
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				msg: []byte{
					0xE3, 0x35, 0x94, 0xD7, 0x50, 0x5E, 0x43, 0xB9,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x33, 0x94, 0xD7, 0x50, 0x5E, 0x43, 0x79, 0xCD,
					0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
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
			name: "rfc8439 test vector",
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
		{
			name: "rfc8439 test vector#1",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: []byte{
				0x76, 0xb8, 0xe0, 0xad, 0xa0, 0xf1, 0x3d, 0x90, 0x40, 0x5d, 0x6a, 0xe5, 0x53, 0x86, 0xbd, 0x28,
				0xbd, 0xd2, 0x19, 0xb8, 0xa0, 0x8d, 0xed, 0x1a, 0xa8, 0x36, 0xef, 0xcc, 0x8b, 0x77, 0x0d, 0xc7,
			},
		},
		{
			name: "rfc8439 test vector#2",
			args: args{
				key: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
				},
				nonce: []byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				},
			},
			want: []byte{
				0xec, 0xfa, 0x25, 0x4f, 0x84, 0x5f, 0x64, 0x74, 0x73, 0xd3, 0xcb, 0x14, 0x0d, 0xa9, 0xe8, 0x76,
				0x06, 0xcb, 0x33, 0x06, 0x6c, 0x44, 0x7b, 0x87, 0xbc, 0x26, 0x66, 0xdd, 0xe3, 0xfb, 0xb7, 0x39,
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
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
				},
			},
			want: []byte{
				0x96, 0x5e, 0x3b, 0xc6, 0xf9, 0xec, 0x7e, 0xd9, 0x56, 0x08, 0x08, 0xf4, 0xd2, 0x29, 0xf9, 0x4b,
				0x13, 0x7f, 0xf2, 0x75, 0xca, 0x9b, 0x3f, 0xcb, 0xdd, 0x59, 0xde, 0xaa, 0xd2, 0x33, 0x10, 0xae,
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
