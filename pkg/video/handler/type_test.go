package handler

import "testing"

func Test_genSecFromDuration(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			"should 90 sec",
			args{"90"},
			90,
			false,
		},
		{
			"should 90 sec",
			args{"90s"},
			90,
			false,
		},
		{
			"should 90 sec",
			args{"1m30s"},
			90,
			false,
		}, {
			"should 90 sec",
			args{"0h1m30s"},
			90,
			false,
		}, {
			"should 90 sec",
			args{" 1m30s "},
			90,
			false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := genSecFromDuration(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("genSecFromDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("genSecFromDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
