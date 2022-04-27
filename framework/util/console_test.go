package util

import "testing"

func TestFmtPrint(t *testing.T) {
	type args struct {
		arr [][]string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				arr: [][]string{
					{"te", "test", "sdf"},
					{"test111", "test123", "78907"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FmtPrint(tt.args.arr)
		})
	}
}
