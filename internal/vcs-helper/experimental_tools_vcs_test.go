package vcshelper

import (
	"testing"
)

func Test_builtinToolVCS_Download(t *testing.T) {
	cfgs := map[string]string{
		"github.com":        "git",
		"git.medlinker.com": "medgit",
	}
	btvcs := NewBuiltinToolVCS(cfgs)

	type args struct {
		remoteURI string
		localDir  string
	}
	tests := []struct {
		name         string
		args         args
		wantRepoRoot string
		wantErr      bool
	}{
		{
			name: "case 0",
			args: args{
				remoteURI: "git.medlinker.com/yeqown/micro-server-template",
				localDir:  "./testdata",
			},
			wantRepoRoot: "testdata/git.medlinker.com/yeqown/micro-server-template",
			wantErr:      false,
		},
		{
			name: "case 1",
			args: args{
				remoteURI: "github.com/yeqown/micro-server-demo",
				localDir:  "./testdata",
			},
			wantRepoRoot: "testdata/github.com/yeqown/micro-server-demo",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRepoRoot, err := btvcs.Download(tt.args.remoteURI, tt.args.localDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("builtinToolVCS.Download() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepoRoot != tt.wantRepoRoot {
				t.Errorf("builtinToolVCS.Download() = %v, want %v", gotRepoRoot, tt.wantRepoRoot)
			}
		})
	}
}
