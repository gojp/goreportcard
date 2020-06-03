package repository

import (
	"reflect"
	"testing"
)

func Test_badgerRepo_Get(t *testing.T) {
	br, _ := NewBadgerRepo("./.badger")
	testkey := []byte("testkey")
	testvalue := []byte("testvalue")

	if err := br.Update(testkey, testvalue); err != nil {
		t.Error(err)
		t.FailNow()
	}

	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "case 0",
			args: args{
				key: testkey,
			},
			want:    testvalue,
			wantErr: false,
		},
		{
			name: "case 0",
			args: args{
				key: []byte("not-ex"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := br.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("badgerRepo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("badgerRepo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
