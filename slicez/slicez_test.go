package slicez

import (
	"reflect"
	"testing"
)

func TestClone(t *testing.T) {
	t.Parallel()
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case 1",
			args: args{
				a: []int{0, 1, 2, 3, 4, 5, 6, 7},
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clone(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloneCast(t *testing.T) {
	t.Parallel()
	type args struct {
		a []string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "string to int",
			args: args{
				a: []string{"0", "1", "2", "3", "4", "5", "6", "7"},
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneCast[int](tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneCast() = %v, want %v", got, tt.want)
			}
		})
	}

	type args2 struct {
		a []float64
	}
	tests2 := []struct {
		name string
		args args2
		want []string
	}{
		{
			name: "float to string",
			args: args2{
				a: []float64{3.3, 4.4, 5.5},
			},
			want: []string{"3.3", "4.4", "5.5"},
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			if got := CloneCast[string](tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneCast() = %v, want %v", got, tt.want)
			}
		})
	}
}
