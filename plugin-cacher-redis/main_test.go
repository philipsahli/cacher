package main

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestKV_Put(t *testing.T) {
	type fields struct {
		logger hclog.Logger
	}
	type args struct {
		key   string
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "",
			fields: fields{
				logger: hclog.New(&hclog.LoggerOptions{
					Level:      hclog.Debug,
					Output:     os.Stderr,
					JSONFormat: true,
				}),
			},
			args: args{
				key:   "testkey",
				value: []byte("testvalue"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := KV{
				Logger: tt.fields.logger,
			}
			if err := k.Put(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("KV.Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
