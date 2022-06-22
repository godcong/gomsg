package gomsg

import (
	"fmt"
	"testing"
)

type data struct {
	s string
	i int
}

func TestNew(t *testing.T) {
	type args struct {
		s   *data
		key string
		fn  []MsgFn[*data]
	}

	tests := []struct {
		name string
		args args
		want *data
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				s: &data{
					s: "111",
					i: 1,
				},
				key: "str111",
				fn: []MsgFn[*data]{
					func(key string, val *data) {
						fmt.Println("called", key, val.s)
					},
				},
			},
			want: &data{
				s: "111",
				i: 1,
			},
		},
		{
			name: "",
			args: args{
				s: &data{
					s: "222",
					i: 2,
				},
				key: "str222",
				fn: []MsgFn[*data]{
					func(key string, val *data) {
						fmt.Println("called", key, val.s)
					},
				},
			},
			want: &data{
				s: "222",
				i: 2,
			},
		},
	}
	m := New[*data](1024)

	for _, tt := range tests {
		for i := range tt.args.fn {
			m.WaitFor(tt.args.key, tt.args.fn[i])
		}
		t.Run(tt.name, func(t *testing.T) {
			m.Send(tt.args.key, tt.args.s)
		})
	}
}
