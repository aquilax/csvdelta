package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_process(t *testing.T) {
	type args struct {
		o options
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			"Empty input produces empty output",
			args{
				options{
					ignoreFirst: false,
					columns:     []int{1},
				},
				strings.NewReader(``),
			},
			``,
			false,
		},
		{
			"Works as expected with correct input",
			args{
				options{
					ignoreFirst: false,
					columns:     []int{1},
				},
				strings.NewReader(`test1,1
test2,10
test3,-10
`),
			},
			`test1,1
test2,9
test3,-20
`,
			false,
		},
		{
			"Ignores the header if requested",
			args{
				options{
					ignoreFirst: true,
					columns:     []int{1},
				},
				strings.NewReader(`name,value
test1,1
test2,10
test3,-10
`),
			},
			`test1,1
test2,9
test3,-20
`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := process(tt.args.o, tt.args.r, w); (err != nil) != tt.wantErr {
				t.Errorf("process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("process() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
