package tdtidy

import (
	"reflect"
	"testing"
)

func Test_chunk(t *testing.T) {
	type args struct {
		items     []string
		chunkSize int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Regular case #1",
			args: args{
				items:     []string{"a", "b", "c", "d", "e", "f"},
				chunkSize: 2,
			},
			want: [][]string{{"a", "b"}, {"c", "d"}, {"e", "f"}},
		},
		{
			name: "Regular case #2",
			args: args{
				items:     []string{"a", "b", "c"},
				chunkSize: 1,
			},
			want: [][]string{{"a"}, {"b"}, {"c"}},
		},
		{
			name: "Regular case #3",
			args: args{
				items:     []string{"a", "b", "c"},
				chunkSize: 5,
			},
			want: [][]string{{"a", "b", "c"}},
		},
		{
			name: "Empty case",
			args: args{
				items:     []string{},
				chunkSize: 2,
			},
			want: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chunk(tt.args.items, tt.args.chunkSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("chunk() = %v, want %v", got, tt.want)
			}
		})
	}
}
