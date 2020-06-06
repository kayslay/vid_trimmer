package service

import (
	"os"
	"testing"
)

func TestGeneratePath(t *testing.T) {
	type args struct {
		key   string
		dType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"file",
			args{"hello", "mp4"},
			os.TempDir() + "/hello.mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GeneratePath(tt.args.key, tt.args.dType); got != tt.want {
				t.Errorf("GeneratePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
