package helper

import (
	"testing"
)

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantOk  bool
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				path: "../httpapi",
			},
			wantOk:  true,
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				path: "./testdata",
			},
			wantOk:  false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOk != tt.wantOk {
				t.Errorf("Exists() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestEnsurePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				path: "./testdata",
			},
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				path: "./testdata",
			},
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				path: "./testdata/testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsurePath(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("EnsurePath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
