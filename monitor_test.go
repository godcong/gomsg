package gomsg

import (
	"fmt"
	"testing"
)

type monitorS struct {
	s string
	i int
}

func TestNew(t *testing.T) {
	type args struct {
		s   *monitorS
		key string
		fn  []MsgFn[*monitorS]
	}

	tests := []struct {
		name string
		args args
		want *monitorS
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				s: &monitorS{
					s: "111",
					i: 1,
				},
				key: "str111",
				fn: []MsgFn[*monitorS]{
					func(key string, val *monitorS) {
						fmt.Println("called", key, val.s)
					},
				},
			},
			want: &monitorS{
				s: "111",
				i: 1,
			},
		},
		{
			name: "",
			args: args{
				s: &monitorS{
					s: "222",
					i: 2,
				},
				key: "str222",
				fn: []MsgFn[*monitorS]{
					func(key string, val *monitorS) {
						fmt.Println("called", key, val.s)
					},
				},
			},
			want: &monitorS{
				s: "222",
				i: 2,
			},
		},
	}
	m := New[*monitorS](1024)

	for _, tt := range tests {
		for i := range tt.args.fn {
			m.WaitFor(tt.args.key, tt.args.fn[i])
		}
		t.Run(tt.name, func(t *testing.T) {
			m.Send(tt.args.key, tt.args.s)
		})
	}
}
