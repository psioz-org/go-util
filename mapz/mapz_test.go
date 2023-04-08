package mapz

import (
	"reflect"
	"testing"
)

func TestJoin(t *testing.T) {
	t.Parallel()
	type args struct {
		out          map[string]string
		excludeEmpty bool
		in           []map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "case 1",
			args: args{
				out:          map[string]string{"a": "1", "b": "2"},
				excludeEmpty: false,
				in: []map[string]string{
					{
						"b": "3",
						"c": "",
					},
					{
						"d": "4",
						"e": "5",
					},
				},
			},
			want: map[string]string{"a": "1", "b": "2", "c": "", "d": "4", "e": "5"},
		},
		{
			name: "case 1",
			args: args{
				out:          map[string]string{"a": "1", "b": "2"},
				excludeEmpty: true,
				in: []map[string]string{
					{
						"b": "3",
						"c": "",
					},
					{
						"d": "4",
						"e": "5",
					},
				},
			},
			want: map[string]string{"a": "1", "b": "2", "d": "4", "e": "5"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Join(tt.args.out, tt.args.excludeEmpty, tt.args.in...)
			if !reflect.DeepEqual(tt.args.out, tt.want) {
				t.Errorf("CloneCast() = %v, want %v", tt.args.out, tt.want)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	t.Parallel()
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "from json string",
			args: args{
				obj: `{"x":{"y":[false,true,{"z":777}]}}`,
			},
			want: map[string]interface{}{
				"x": map[string]interface{}{
					"y": []interface{}{
						false,
						true,
						map[string]interface{}{
							"z": 777.0,
						},
					},
				},
			},
		},
		{
			name: "from byte array of json string",
			args: args{
				obj: []byte(`{"x":{"y":[false,true,{"z":777}]}}`),
			},
			want: map[string]interface{}{
				"x": map[string]interface{}{
					"y": []interface{}{
						false,
						true,
						map[string]interface{}{
							"z": 777.0,
						},
					},
				},
			},
		},
		{
			name: "from random object",
			args: args{
				obj: map[string]map[string]interface{}{
					"x": {
						"y": []interface{}{
							false,
							true,
							map[string]int{
								"z": 777,
							},
						},
					},
				},
			},
			want: map[string]interface{}{
				"x": map[string]interface{}{
					"y": []interface{}{
						false,
						true,
						map[string]interface{}{
							"z": 777.0,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToMap(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringMap(t *testing.T) {
	t.Parallel()
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "from json string",
			args: args{
				obj: `{"x":true,"y":"y","z":777,"empty":"","null":null}`,
			},
			want: map[string]string{
				"x":     "true",
				"y":     "y",
				"z":     "777",
				"empty": "",
				"null":  "<nil>",
			},
		},
		{
			name: "from byte array of json string",
			args: args{
				obj: []byte(`{"x":true,"y":"y","z":777,"empty":"","null":null}`),
			},
			want: map[string]string{
				"x":     "true",
				"y":     "y",
				"z":     "777",
				"empty": "",
				"null":  "<nil>",
			},
		},
		{
			name: "from random object",
			args: args{
				obj: map[string]interface{}{
					"x":     true,
					"y":     "y",
					"z":     777,
					"empty": "",
					"null":  nil,
				},
			},
			want: map[string]string{
				"x":     "true",
				"y":     "y",
				"z":     "777",
				"empty": "",
				"null":  "<nil>",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStringMap(tt.args.obj); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToStringMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
