package stringsz

import (
	"context"
	"regexp"
	"testing"
)

func TestGetVersionAsInteger(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1 digit",
			args: args{
				version: "1.2.3",
			},
			want: "1002003",
		},
		{
			name: "2 digit",
			args: args{
				version: "11.22.33",
			},
			want: "11022033",
		},
		{
			name: "3 digit",
			args: args{
				version: "111.222.333",
			},
			want: "111222333",
		},
		{
			name: "4 digit",
			args: args{
				version: "1111.2222.3333",
			},
			want: "1111999999",
		},
		{
			name: "2 digit mix string",
			args: args{
				version: "ver11.22.33beta",
			},
			want: "11022033",
		},
		{
			name: "4 digit mix string",
			args: args{
				version: "ver#1111.2222.3333alpha",
			},
			want: "1111999999",
		},
		{
			name: "No number",
			args: args{
				version: "x",
			},
			want: "0",
		},
		{
			name: "0.0.0",
			args: args{
				version: "0.0.0",
			},
			want: "0",
		},
		{
			name: "0.0.x",
			args: args{
				version: "0.0.2",
			},
			want: "2",
		},
		{
			name: "0.y.x",
			args: args{
				version: "0.3.2",
			},
			want: "3002",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVersionAsInteger(tt.args.version); got != tt.want {
				t.Errorf("GetVersionAsInteger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOfNth(t *testing.T) {
	type args struct {
		s      string
		substr string
		nth    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "nth -1=1",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    -1,
			},
			want: 3,
		},
		{
			name: "nth 0=1",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    0,
			},
			want: 3,
		},
		{
			name: "nth 1",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    1,
			},
			want: 3,
		},
		{
			name: "nth 2",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    2,
			},
			want: 9,
		},
		{
			name: "nth 3",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    3,
			},
			want: 15,
		},
		{
			name: "nth 4",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    4,
			},
			want: -1,
		},
		{
			name: "nth 777",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "_x_",
				nth:    777,
			},
			want: -1,
		},
		{
			name: "substr empty return 0",
			args: args{
				s:      "abc_x_def_x_ghi_x_jkl",
				substr: "",
				nth:    1,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOfNth(tt.args.s, tt.args.substr, tt.args.nth); got != tt.want {
				t.Errorf("IndexOfNth() = %v, want %v", got, tt.want)
			}
		})
	}
}

type contextKey string

func (c contextKey) String() string {
	return "mypackage context key " + string(c)
}

func TestPrintContextInternals(t *testing.T) {
	type args struct {
		ctx   interface{}
		inner bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "for debug: anything is fine, inner: false",
			args: args{
				ctx:   context.WithValue(context.Background(), contextKey("key1"), "value1"),
				inner: true,
			},
		},
		{
			name: "for debug: anything is fine, inner: true",
			args: args{
				ctx:   context.WithValue(context.Background(), contextKey("key1"), "value1"),
				inner: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrintContextInternals(tt.args.ctx, tt.args.inner)
		})
	}
}

func TestReplaceAllStringSubmatchFunc(t *testing.T) {
	type args struct {
		re   *regexp.Regexp
		str  string
		repl func([]string) string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				re:  regexp.MustCompile(`(\w+):(.*?)(\d+?)`),
				str: "$email:someone789@gmail.com",
				repl: func(ss []string) string {
					return ss[3] + "_" + ss[2] + "_" + ss[1]
				},
			},
			want: "$7_someone_email89@gmail.com",
		},
		{
			name: "case 2 with empty group, not -1 case",
			args: args{
				re:  regexp.MustCompile(`(\w+)(x?):(.*?)(\d+?)`),
				str: "$email:someone789@gmail.com",
				repl: func(ss []string) string {
					return ss[4] + "_" + ss[3] + "_" + ss[2] + "_" + ss[1]
				},
			},
			want: "$7_someone__email89@gmail.com",
		},
		{
			name: "case 3 with no match",
			args: args{
				re:  regexp.MustCompile(`(xyz)`),
				str: "$email:someone789@gmail.com",
				repl: func(ss []string) string {
					return "aa" + ss[1] + "aa"
				},
			},
			want: "$email:someone789@gmail.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReplaceAllStringSubmatchFunc(tt.args.re, tt.args.str, tt.args.repl); got != tt.want {
				t.Errorf("ReplaceAllStringSubmatchFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSnake2Title(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "general",
			args: args{
				s: "i_love_golang",
			},
			want: "I Love Golang",
		},
		{
			name: "single",
			args: args{
				s: "i",
			},
			want: "I",
		},
		{
			name: "empty",
			args: args{
				s: "",
			},
			want: "",
		},
		{
			name: "prefix suffix or multiple sep",
			args: args{
				s: "___i__love__golang___",
			},
			want: "I Love Golang",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Snake2Title(tt.args.s); got != tt.want {
				t.Errorf("Snake2Title() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToCrc32(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "case 1",
			args: args{
				v: "test2crc",
			},
			want: "CC5374E4",
		},
		{
			name: "case 2",
			args: args{
				v: "777",
			},
			want: "F6DF2F3C",
		},
		{
			name: "case 3",
			args: args{
				v: 777,
			},
			want: "F6DF2F3C",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToCrc32(tt.args.v); got != tt.want {
				t.Errorf("ToCrc32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJson(t *testing.T) {
	type args struct {
		obj    interface{}
		indent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "string",
			args: args{
				obj:    "string1",
				indent: "\t",
			},
			want: `"string1"`,
		},
		{
			name: "array int",
			args: args{
				obj:    []int{7, 8, 9},
				indent: "\t",
			},
			want: "[\n\t7,\n\t8,\n\t9\n]",
		},
		{
			name: "array mix",
			args: args{
				obj:    []interface{}{7, "8", true, nil},
				indent: "\t",
			},
			want: "[\n\t7,\n\t\"8\",\n\ttrue,\n\tnull\n]",
		},
		{
			name: "map[string]interface{}",
			args: args{
				obj:    map[string]interface{}{"a": 7, "b": "8", "c": true, "d": nil},
				indent: "\t",
			},
			want: "{\n\t\"a\": 7,\n\t\"b\": \"8\",\n\t\"c\": true,\n\t\"d\": null\n}",
		},
		{
			name: "map[string]interface{} with zz indent",
			args: args{
				obj:    map[string]interface{}{"a": 7, "b": "8", "c": true, "d": nil},
				indent: "zz",
			},
			want: "{\nzz\"a\": 7,\nzz\"b\": \"8\",\nzz\"c\": true,\nzz\"d\": null\n}",
		},
		{
			name: "map[string]interface{} without indent",
			args: args{
				obj:    map[string]interface{}{"a": 7, "b": "8", "c": true, "d": nil},
				indent: "",
			},
			want: "{\"a\":7,\"b\":\"8\",\"c\":true,\"d\":null}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToJson(tt.args.obj, tt.args.indent); got != tt.want {
				t.Errorf("ToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
