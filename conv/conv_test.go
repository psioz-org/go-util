package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"testing"
)

type other string

const (
	jsonFullCheckMapStr = `{"empty":"","false":false,"float":777.7,"int":777,"int string":"777","null":null,"string":"string","true":true}`
)

var (
	jsonFullCheckMapObj = map[string]any{"empty": "", "false": false, "float": 777.7, "int": 777, "int string": "777", "null": nil, "string": "string", "true": true}
	//When unmarshal any number is float
	jsonFullCheckMapObjOnlyFloat = map[string]any{"empty": "", "false": false, "float": 777.7, "int": 777.0, "int string": "777", "null": nil, "string": "string", "true": true}
)

type Shape interface {
	area() float64
	area2() float64
}
type Rect struct {
	X    float64
	Y    float64
	w, h float64
}

func (r Rect) area() float64 {
	return r.w * r.h
}

func (r *Rect) area2() float64 {
	return r.w * r.h
}

func (r *Rect) area3() float64 {
	return r.w * r.h
}

func func1() {
	fmt.Sprintln("func1 body")
}

func TestGetItemTestGeneric(t *testing.T) {
	t.Parallel()
	type args struct {
		obj  any
		keys string
		sep  string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "test generic success",
			args: args{
				obj: map[string]any{
					"x": map[string]any{
						"y": map[string]any{
							"z": 3,
						},
					},
				},
				keys: "x.y.z",
				sep:  ".",
			},
			want: 3,
		},
		{
			name: "test generic success array index",
			args: args{
				obj: map[string]any{
					"x": map[string]any{
						"y": []any{
							1,
							2,
						},
					},
				},
				keys: "x/y/1",
				sep:  "/",
			},
			want: 2,
		},
		{
			name: "test generic fail cast to result type",
			args: args{
				obj: map[string]any{
					"x": map[string]any{
						"y": map[string]any{
							"z": "3.3",
						},
					},
				},
				keys: "x.y.z",
				sep:  ".",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetItem[int](tt.args.obj, tt.args.keys, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetItem(t *testing.T) {
	t.Parallel()
	type args struct {
		obj  any
		keys string
		sep  string
	}
	type result struct{}
	result1 := result{}
	result2 := result{}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "test result success",
			args: args{
				obj: map[string]any{
					"x": map[string]any{
						"y": map[string]any{
							"z": result1,
						},
					},
				},
				keys: "x/y/z",
				sep:  "/",
			},
			want: result1,
		},
		{
			name: "test result success array index",
			args: args{
				obj: map[string]any{
					"x": map[string]any{
						"y": []any{
							result1,
							result2,
						},
					},
				},
				keys: "x/y/1",
				sep:  "/",
			},
			want: result2,
		},
		{
			name: "fail cast to map to get key",
			args: args{
				obj: map[string]any{
					"x": map[string]any{},
				},
				keys: "x/y/z",
				sep:  "/",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail cast to array to get key",
			args: args{
				obj: map[string]any{
					"x": map[string]any{},
				},
				keys: "x/y/0",
				sep:  "/",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get nil for non-exists",
			args: args{
				obj: map[string]any{
					"x": map[string]any{},
				},
				keys: "x/y",
				sep:  "/",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetItem[any](tt.args.obj, tt.args.keys, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernary(t *testing.T) {
	t.Parallel()
	type args[T any] struct {
		b  bool
		vt T
		vf T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[string]{
		{
			name: "case true",
			args: args[string]{
				b:  true,
				vt: "value if true",
				vf: "value if false",
			},
			want: "value if true",
		},
		{
			name: "case false",
			args: args[string]{
				b:  false,
				vt: "value if true",
				vf: "value if false",
			},
			want: "value if false",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.b, tt.args.vt, tt.args.vf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernaryInt(t *testing.T) {
	t.Parallel()
	type args[T any] struct {
		b  bool
		vt T
		vf T
	}

	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "case true",
			args: args[int]{
				b:  true,
				vt: 2,
				vf: 4,
			},
			want: 2,
		},
		{
			name: "case false",
			args: args[int]{
				b:  false,
				vt: 2,
				vf: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.args.b, tt.args.vt, tt.args.vf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTo(t *testing.T) {
	t.Parallel()
	type args struct {
		obj any
	}
	type testCase[T any] struct {
		name    string
		args    args
		wantOut T
		wantErr bool
		wantReg *regexp.Regexp
	}
	const boolTo = "bool"
	testsBool := []testCase[bool]{
		{
			name: "bool > " + boolTo,
			args: args{
				obj: true,
			},
			wantOut: true,
		},
		{
			name: "float32 0 > " + boolTo,
			args: args{
				obj: float32(0),
			},
			wantOut: false,
		},
		{
			name: "float32 1 > " + boolTo,
			args: args{
				obj: float32(1),
			},
			wantOut: true,
		},
		{
			name: "int 0 > " + boolTo,
			args: args{
				obj: 0,
			},
			wantOut: false,
		},
		{
			name: "int 1 > " + boolTo,
			args: args{
				obj: 1,
			},
			wantOut: true,
		},
		{
			name: "int 2 > " + boolTo,
			args: args{
				obj: 2,
			},
			wantOut: true,
		},
		{
			name: "string false > " + boolTo,
			args: args{
				obj: "false",
			},
			wantOut: false,
		},
		{
			name: "string true > " + boolTo,
			args: args{
				obj: "true",
			},
			wantOut: true,
		},
		{
			name: "string 0 > " + boolTo,
			args: args{
				obj: "0",
			},
			wantOut: false,
		},
		{
			name: "string 1 > " + boolTo,
			args: args{
				obj: "1",
			},
			wantOut: true,
		},
		{
			name: "string empty > " + boolTo,
			args: args{
				obj: "",
			},
			wantOut: false,
		},
		{
			name: "nil > " + boolTo,
			args: args{
				obj: nil,
			},
			wantOut: false,
		},
	}
	for _, tt := range testsBool {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[bool](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const complex64To = "complex64"
	complex64Tests := []testCase[complex64]{
		{
			name: "string > " + complex64To,
			args: args{
				obj: "(77+0i)",
			},
			wantOut: complex(float32(77), float32(0)),
		},
		{
			name: "string without parenthesis > " + complex64To,
			args: args{
				obj: "77+0i",
			},
			wantOut: complex(float32(77), float32(0)),
		},
	}
	for _, tt := range complex64Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[complex64](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const complex128To = "complex128"
	complex128Tests := []testCase[complex128]{
		{
			name: "string > " + complex128To,
			args: args{
				obj: "(77+0i)",
			},
			wantOut: complex(float64(77), float64(0)),
		},
		{
			name: "string without parenthesis > " + complex128To,
			args: args{
				obj: "77+0i",
			},
			wantOut: complex(float64(77), float64(0)),
		},
	}
	for _, tt := range complex128Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[complex128](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const float32To = "float32"
	float32Tests := []testCase[float32]{
		{
			name: "complex64 > " + float32To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + float32To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + float32To,
			args: args{
				obj: float32(77),
			},
			wantOut: float32(77),
		},
		{
			name: "float64 > " + float32To,
			args: args{
				obj: float64(77),
			},
			wantOut: float32(77),
		},
		{
			name: "float32 dot0 > " + float32To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: float32(77),
		},
		{
			name: "float64 dot0 > " + float32To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: float32(77),
		},
		{
			name: "float32 dot1 > " + float32To,
			args: args{
				obj: float32(77.1),
			},
			wantOut: float32(77.1),
		},
		{
			name: "float64 dot1 > " + float32To,
			args: args{
				obj: float64(77.1),
			},
			wantOut: float32(77.1),
		},
		{
			name: "int > " + float32To,
			args: args{
				obj: 77,
			},
			wantOut: float32(77),
		},
		{
			name: "int8 > " + float32To,
			args: args{
				obj: int8(77),
			},
			wantOut: float32(77),
		},
		{
			name: "int16 > " + float32To,
			args: args{
				obj: int16(77),
			},
			wantOut: float32(77),
		},
		{
			name: "int32 > " + float32To,
			args: args{
				obj: int32(77),
			},
			wantOut: float32(77),
		},
		{
			name: "int64 > " + float32To,
			args: args{
				obj: int64(77),
			},
			wantOut: float32(77),
		},
		{
			name: "uint > " + float32To,
			args: args{
				obj: uint(77),
			},
			wantOut: float32(77),
		},
		{
			name: "uint8 > " + float32To,
			args: args{
				obj: uint8(77),
			},
			wantOut: float32(77),
		},
		{
			name: "uint16 > " + float32To,
			args: args{
				obj: uint16(77),
			},
			wantOut: float32(77),
		},
		{
			name: "uint32 > " + float32To,
			args: args{
				obj: uint32(77),
			},
			wantOut: float32(77),
		},
		{
			name: "uint64 > " + float32To,
			args: args{
				obj: uint64(77),
			},
			wantOut: float32(77),
		},
		{
			name: "string > " + float32To,
			args: args{
				obj: "77",
			},
			wantOut: float32(77),
		},
	}
	for _, tt := range float32Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[float32](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const float64To = "float64"
	float64Tests := []testCase[float64]{
		{
			name: "complex64 > " + float64To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + float64To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + float64To,
			args: args{
				obj: float32(77),
			},
			wantOut: float64(77),
		},
		{
			name: "float64 > " + float64To,
			args: args{
				obj: float64(77),
			},
			wantOut: float64(77),
		},
		{
			name: "float32 dot0 > " + float64To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: float64(77),
		},
		{
			name: "float64 dot0 > " + float64To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: float64(77),
		},
		{
			name: "float32 dot1 > " + float64To,
			args: args{
				obj: float32(77.1),
			},
			wantOut: float64(77.1),
		},
		{
			name: "float64 dot1 > " + float64To,
			args: args{
				obj: float64(77.1),
			},
			wantOut: float64(77.1),
		},
		{
			name: "int > " + float64To,
			args: args{
				obj: 77,
			},
			wantOut: float64(77),
		},
		{
			name: "int8 > " + float64To,
			args: args{
				obj: int8(77),
			},
			wantOut: float64(77),
		},
		{
			name: "int16 > " + float64To,
			args: args{
				obj: int16(77),
			},
			wantOut: float64(77),
		},
		{
			name: "int32 > " + float64To,
			args: args{
				obj: int32(77),
			},
			wantOut: float64(77),
		},
		{
			name: "int64 > " + float64To,
			args: args{
				obj: int64(77),
			},
			wantOut: float64(77),
		},
		{
			name: "uint > " + float64To,
			args: args{
				obj: uint(77),
			},
			wantOut: float64(77),
		},
		{
			name: "uint8 > " + float64To,
			args: args{
				obj: uint8(77),
			},
			wantOut: float64(77),
		},
		{
			name: "uint16 > " + float64To,
			args: args{
				obj: uint16(77),
			},
			wantOut: float64(77),
		},
		{
			name: "uint32 > " + float64To,
			args: args{
				obj: uint32(77),
			},
			wantOut: float64(77),
		},
		{
			name: "uint64 > " + float64To,
			args: args{
				obj: uint64(77),
			},
			wantOut: float64(77),
		},
		{
			name: "string > " + float64To,
			args: args{
				obj: "77",
			},
			wantOut: float64(77),
		},
	}
	for _, tt := range float64Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[float64](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const intTo = "int"
	intTests := []testCase[int]{
		{
			name: "complex64 > " + intTo,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + intTo,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + intTo,
			args: args{
				obj: float32(77),
			},
			wantOut: int(77),
		},
		{
			name: "float64 > " + intTo,
			args: args{
				obj: float64(77),
			},
			wantOut: int(77),
		},
		{
			name: "float32 dot0 > " + intTo,
			args: args{
				obj: float32(77.0),
			},
			wantOut: int(77),
		},
		{
			name: "float64 dot0 > " + intTo,
			args: args{
				obj: float64(77.0),
			},
			wantOut: int(77),
		},
		{
			name: "float32 dot1 > " + intTo,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + intTo,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + intTo,
			args: args{
				obj: 77,
			},
			wantOut: int(77),
		},
		{
			name: "int8 > " + intTo,
			args: args{
				obj: int8(77),
			},
			wantOut: int(77),
		},
		{
			name: "int16 > " + intTo,
			args: args{
				obj: int16(77),
			},
			wantOut: int(77),
		},
		{
			name: "int32 > " + intTo,
			args: args{
				obj: int32(77),
			},
			wantOut: int(77),
		},
		{
			name: "int64 > " + intTo,
			args: args{
				obj: int64(77),
			},
			wantOut: int(77),
		},
		{
			name: "uint > " + intTo,
			args: args{
				obj: uint(77),
			},
			wantOut: int(77),
		},
		{
			name: "uint8 > " + intTo,
			args: args{
				obj: uint8(77),
			},
			wantOut: int(77),
		},
		{
			name: "uint16 > " + intTo,
			args: args{
				obj: uint16(77),
			},
			wantOut: int(77),
		},
		{
			name: "uint32 > " + intTo,
			args: args{
				obj: uint32(77),
			},
			wantOut: int(77),
		},
		{
			name: "uint64 > " + intTo,
			args: args{
				obj: uint64(77),
			},
			wantOut: int(77),
		},
		{
			name: "string > " + intTo,
			args: args{
				obj: "77",
			},
			wantOut: int(77),
		},
	}
	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[int](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const int8To = "int8"
	int8Tests := []testCase[int8]{
		{
			name: "complex64 > " + int8To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + int8To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + int8To,
			args: args{
				obj: float32(77),
			},
			wantOut: int8(77),
		},
		{
			name: "float64 > " + int8To,
			args: args{
				obj: float64(77),
			},
			wantOut: int8(77),
		},
		{
			name: "float32 dot0 > " + int8To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: int8(77),
		},
		{
			name: "float64 dot0 > " + int8To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: int8(77),
		},
		{
			name: "float32 dot1 > " + int8To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + int8To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + int8To,
			args: args{
				obj: 77,
			},
			wantOut: int8(77),
		},
		{
			name: "int8 > " + int8To,
			args: args{
				obj: int8(77),
			},
			wantOut: int8(77),
		},
		{
			name: "int16 > " + int8To,
			args: args{
				obj: int16(77),
			},
			wantOut: int8(77),
		},
		{
			name: "int32 > " + int8To,
			args: args{
				obj: int32(77),
			},
			wantOut: int8(77),
		},
		{
			name: "int64 > " + int8To,
			args: args{
				obj: int64(77),
			},
			wantOut: int8(77),
		},
		{
			name: "uint > " + int8To,
			args: args{
				obj: uint(77),
			},
			wantOut: int8(77),
		},
		{
			name: "uint8 > " + int8To,
			args: args{
				obj: uint8(77),
			},
			wantOut: int8(77),
		},
		{
			name: "uint16 > " + int8To,
			args: args{
				obj: uint16(77),
			},
			wantOut: int8(77),
		},
		{
			name: "uint32 > " + int8To,
			args: args{
				obj: uint32(77),
			},
			wantOut: int8(77),
		},
		{
			name: "uint64 > " + int8To,
			args: args{
				obj: uint64(77),
			},
			wantOut: int8(77),
		},
		{
			name: "string > " + int8To,
			args: args{
				obj: "77",
			},
			wantOut: int8(77),
		},
	}
	for _, tt := range int8Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[int8](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const int16To = "int16"
	int16Tests := []testCase[int16]{
		{
			name: "complex64 > " + int16To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + int16To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + int16To,
			args: args{
				obj: float32(77),
			},
			wantOut: int16(77),
		},
		{
			name: "float64 > " + int16To,
			args: args{
				obj: float64(77),
			},
			wantOut: int16(77),
		},
		{
			name: "float32 dot0 > " + int16To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: int16(77),
		},
		{
			name: "float64 dot0 > " + int16To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: int16(77),
		},
		{
			name: "float32 dot1 > " + int16To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + int16To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + int16To,
			args: args{
				obj: 77,
			},
			wantOut: int16(77),
		},
		{
			name: "int8 > " + int16To,
			args: args{
				obj: int8(77),
			},
			wantOut: int16(77),
		},
		{
			name: "int16 > " + int16To,
			args: args{
				obj: int16(77),
			},
			wantOut: int16(77),
		},
		{
			name: "int32 > " + int16To,
			args: args{
				obj: int32(77),
			},
			wantOut: int16(77),
		},
		{
			name: "int64 > " + int16To,
			args: args{
				obj: int64(77),
			},
			wantOut: int16(77),
		},
		{
			name: "uint > " + int16To,
			args: args{
				obj: uint(77),
			},
			wantOut: int16(77),
		},
		{
			name: "uint8 > " + int16To,
			args: args{
				obj: uint8(77),
			},
			wantOut: int16(77),
		},
		{
			name: "uint16 > " + int16To,
			args: args{
				obj: uint16(77),
			},
			wantOut: int16(77),
		},
		{
			name: "uint32 > " + int16To,
			args: args{
				obj: uint32(77),
			},
			wantOut: int16(77),
		},
		{
			name: "uint64 > " + int16To,
			args: args{
				obj: uint64(77),
			},
			wantOut: int16(77),
		},
		{
			name: "string > " + int16To,
			args: args{
				obj: "77",
			},
			wantOut: int16(77),
		},
	}
	for _, tt := range int16Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[int16](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const int32To = "int32"
	int32Tests := []testCase[int32]{
		{
			name: "complex64 > " + int32To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + int32To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + int32To,
			args: args{
				obj: float32(77),
			},
			wantOut: int32(77),
		},
		{
			name: "float64 > " + int32To,
			args: args{
				obj: float64(77),
			},
			wantOut: int32(77),
		},
		{
			name: "float32 dot0 > " + int32To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: int32(77),
		},
		{
			name: "float64 dot0 > " + int32To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: int32(77),
		},
		{
			name: "float32 dot1 > " + int32To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + int32To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + int32To,
			args: args{
				obj: 77,
			},
			wantOut: int32(77),
		},
		{
			name: "int8 > " + int32To,
			args: args{
				obj: int8(77),
			},
			wantOut: int32(77),
		},
		{
			name: "int16 > " + int32To,
			args: args{
				obj: int16(77),
			},
			wantOut: int32(77),
		},
		{
			name: "int32 > " + int32To,
			args: args{
				obj: int32(77),
			},
			wantOut: int32(77),
		},
		{
			name: "int64 > " + int32To,
			args: args{
				obj: int64(77),
			},
			wantOut: int32(77),
		},
		{
			name: "uint > " + int32To,
			args: args{
				obj: uint(77),
			},
			wantOut: int32(77),
		},
		{
			name: "uint8 > " + int32To,
			args: args{
				obj: uint8(77),
			},
			wantOut: int32(77),
		},
		{
			name: "uint16 > " + int32To,
			args: args{
				obj: uint16(77),
			},
			wantOut: int32(77),
		},
		{
			name: "uint32 > " + int32To,
			args: args{
				obj: uint32(77),
			},
			wantOut: int32(77),
		},
		{
			name: "uint64 > " + int32To,
			args: args{
				obj: uint64(77),
			},
			wantOut: int32(77),
		},
		{
			name: "string > " + int32To,
			args: args{
				obj: "77",
			},
			wantOut: int32(77),
		},
	}
	for _, tt := range int32Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[int32](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const int64To = "int64"
	int64Tests := []testCase[int64]{
		{
			name: "complex64 > " + int64To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + int64To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + int64To,
			args: args{
				obj: float32(77),
			},
			wantOut: int64(77),
		},
		{
			name: "float64 > " + int64To,
			args: args{
				obj: float64(77),
			},
			wantOut: int64(77),
		},
		{
			name: "float32 dot0 > " + int64To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: int64(77),
		},
		{
			name: "float64 dot0 > " + int64To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: int64(77),
		},
		{
			name: "float32 dot1 > " + int64To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + int64To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + int64To,
			args: args{
				obj: 77,
			},
			wantOut: int64(77),
		},
		{
			name: "int8 > " + int64To,
			args: args{
				obj: int8(77),
			},
			wantOut: int64(77),
		},
		{
			name: "int16 > " + int64To,
			args: args{
				obj: int16(77),
			},
			wantOut: int64(77),
		},
		{
			name: "int32 > " + int64To,
			args: args{
				obj: int32(77),
			},
			wantOut: int64(77),
		},
		{
			name: "int64 > " + int64To,
			args: args{
				obj: int64(77),
			},
			wantOut: int64(77),
		},
		{
			name: "uint > " + int64To,
			args: args{
				obj: uint(77),
			},
			wantOut: int64(77),
		},
		{
			name: "uint8 > " + int64To,
			args: args{
				obj: uint8(77),
			},
			wantOut: int64(77),
		},
		{
			name: "uint16 > " + int64To,
			args: args{
				obj: uint16(77),
			},
			wantOut: int64(77),
		},
		{
			name: "uint32 > " + int64To,
			args: args{
				obj: uint32(77),
			},
			wantOut: int64(77),
		},
		{
			name: "uint64 > " + int64To,
			args: args{
				obj: uint64(77),
			},
			wantOut: int64(77),
		},
		{
			name: "string > " + int64To,
			args: args{
				obj: "77",
			},
			wantOut: int64(77),
		},
	}
	for _, tt := range int64Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[int64](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const uintTo = "uint"
	uintTests := []testCase[uint]{
		{
			name: "complex64 > " + uintTo,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + uintTo,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + uintTo,
			args: args{
				obj: float32(77),
			},
			wantOut: uint(77),
		},
		{
			name: "float64 > " + uintTo,
			args: args{
				obj: float64(77),
			},
			wantOut: uint(77),
		},
		{
			name: "float32 dot0 > " + uintTo,
			args: args{
				obj: float32(77.0),
			},
			wantOut: uint(77),
		},
		{
			name: "float64 dot0 > " + uintTo,
			args: args{
				obj: float64(77.0),
			},
			wantOut: uint(77),
		},
		{
			name: "float32 dot1 > " + uintTo,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + uintTo,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + uintTo,
			args: args{
				obj: 77,
			},
			wantOut: uint(77),
		},
		{
			name: "int8 > " + uintTo,
			args: args{
				obj: int8(77),
			},
			wantOut: uint(77),
		},
		{
			name: "int16 > " + uintTo,
			args: args{
				obj: int16(77),
			},
			wantOut: uint(77),
		},
		{
			name: "int32 > " + uintTo,
			args: args{
				obj: int32(77),
			},
			wantOut: uint(77),
		},
		{
			name: "int64 > " + uintTo,
			args: args{
				obj: int64(77),
			},
			wantOut: uint(77),
		},
		{
			name: "uint > " + uintTo,
			args: args{
				obj: uint(77),
			},
			wantOut: uint(77),
		},
		{
			name: "uint8 > " + uintTo,
			args: args{
				obj: uint8(77),
			},
			wantOut: uint(77),
		},
		{
			name: "uint16 > " + uintTo,
			args: args{
				obj: uint16(77),
			},
			wantOut: uint(77),
		},
		{
			name: "uint32 > " + uintTo,
			args: args{
				obj: uint32(77),
			},
			wantOut: uint(77),
		},
		{
			name: "uint64 > " + uintTo,
			args: args{
				obj: uint64(77),
			},
			wantOut: uint(77),
		},
		{
			name: "string > " + uintTo,
			args: args{
				obj: "77",
			},
			wantOut: uint(77),
		},
	}
	for _, tt := range uintTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uint](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const uint8To = "uint8"
	uint8Tests := []testCase[uint8]{
		{
			name: "complex64 > " + uint8To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + uint8To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + uint8To,
			args: args{
				obj: float32(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "float64 > " + uint8To,
			args: args{
				obj: float64(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "float32 dot0 > " + uint8To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: uint8(77),
		},
		{
			name: "float64 dot0 > " + uint8To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: uint8(77),
		},
		{
			name: "float32 dot1 > " + uint8To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + uint8To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + uint8To,
			args: args{
				obj: 77,
			},
			wantOut: uint8(77),
		},
		{
			name: "int8 > " + uint8To,
			args: args{
				obj: int8(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "int16 > " + uint8To,
			args: args{
				obj: int16(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "int32 > " + uint8To,
			args: args{
				obj: int32(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "int64 > " + uint8To,
			args: args{
				obj: int64(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "uint > " + uint8To,
			args: args{
				obj: uint(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "uint8 > " + uint8To,
			args: args{
				obj: uint8(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "uint16 > " + uint8To,
			args: args{
				obj: uint16(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "uint32 > " + uint8To,
			args: args{
				obj: uint32(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "uint64 > " + uint8To,
			args: args{
				obj: uint64(77),
			},
			wantOut: uint8(77),
		},
		{
			name: "string > " + uint8To,
			args: args{
				obj: "77",
			},
			wantOut: uint8(77),
		},
	}
	for _, tt := range uint8Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uint8](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const uint16To = "uint16"
	uint16Tests := []testCase[uint16]{
		{
			name: "complex64 > " + uint16To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + uint16To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + uint16To,
			args: args{
				obj: float32(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "float64 > " + uint16To,
			args: args{
				obj: float64(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "float32 dot0 > " + uint16To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: uint16(77),
		},
		{
			name: "float64 dot0 > " + uint16To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: uint16(77),
		},
		{
			name: "float32 dot1 > " + uint16To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + uint16To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + uint16To,
			args: args{
				obj: 77,
			},
			wantOut: uint16(77),
		},
		{
			name: "int8 > " + uint16To,
			args: args{
				obj: int8(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "int16 > " + uint16To,
			args: args{
				obj: int16(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "int32 > " + uint16To,
			args: args{
				obj: int32(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "int64 > " + uint16To,
			args: args{
				obj: int64(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "uint > " + uint16To,
			args: args{
				obj: uint(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "uint8 > " + uint16To,
			args: args{
				obj: uint8(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "uint16 > " + uint16To,
			args: args{
				obj: uint16(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "uint32 > " + uint16To,
			args: args{
				obj: uint32(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "uint64 > " + uint16To,
			args: args{
				obj: uint64(77),
			},
			wantOut: uint16(77),
		},
		{
			name: "string > " + uint16To,
			args: args{
				obj: "77",
			},
			wantOut: uint16(77),
		},
	}
	for _, tt := range uint16Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uint16](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const uint32To = "uint32"
	uint32Tests := []testCase[uint32]{
		{
			name: "complex64 > " + uint32To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + uint32To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + uint32To,
			args: args{
				obj: float32(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "float64 > " + uint32To,
			args: args{
				obj: float64(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "float32 dot0 > " + uint32To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: uint32(77),
		},
		{
			name: "float64 dot0 > " + uint32To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: uint32(77),
		},
		{
			name: "float32 dot1 > " + uint32To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + uint32To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + uint32To,
			args: args{
				obj: 77,
			},
			wantOut: uint32(77),
		},
		{
			name: "int8 > " + uint32To,
			args: args{
				obj: int8(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "int16 > " + uint32To,
			args: args{
				obj: int16(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "int32 > " + uint32To,
			args: args{
				obj: int32(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "int64 > " + uint32To,
			args: args{
				obj: int64(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "uint > " + uint32To,
			args: args{
				obj: uint(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "uint8 > " + uint32To,
			args: args{
				obj: uint8(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "uint16 > " + uint32To,
			args: args{
				obj: uint16(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "uint32 > " + uint32To,
			args: args{
				obj: uint32(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "uint64 > " + uint32To,
			args: args{
				obj: uint64(77),
			},
			wantOut: uint32(77),
		},
		{
			name: "string > " + uint32To,
			args: args{
				obj: "77",
			},
			wantOut: uint32(77),
		},
	}
	for _, tt := range uint32Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uint32](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const uint64To = "uint64"
	uint64Tests := []testCase[uint64]{
		{
			name: "complex64 > " + uint64To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantErr: true,
		},
		{
			name: "complex128 > " + uint64To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantErr: true,
		},
		{
			name: "float32 > " + uint64To,
			args: args{
				obj: float32(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "float64 > " + uint64To,
			args: args{
				obj: float64(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "float32 dot0 > " + uint64To,
			args: args{
				obj: float32(77.0),
			},
			wantOut: uint64(77),
		},
		{
			name: "float64 dot0 > " + uint64To,
			args: args{
				obj: float64(77.0),
			},
			wantOut: uint64(77),
		},
		{
			name: "float32 dot1 > " + uint64To,
			args: args{
				obj: float32(77.1),
			},
			wantErr: true,
		},
		{
			name: "float64 dot1 > " + uint64To,
			args: args{
				obj: float64(77.1),
			},
			wantErr: true,
		},
		{
			name: "int > " + uint64To,
			args: args{
				obj: 77,
			},
			wantOut: uint64(77),
		},
		{
			name: "int8 > " + uint64To,
			args: args{
				obj: int8(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "int16 > " + uint64To,
			args: args{
				obj: int16(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "int32 > " + uint64To,
			args: args{
				obj: int32(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "int64 > " + uint64To,
			args: args{
				obj: int64(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "uint > " + uint64To,
			args: args{
				obj: uint(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "uint8 > " + uint64To,
			args: args{
				obj: uint8(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "uint16 > " + uint64To,
			args: args{
				obj: uint16(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "uint32 > " + uint64To,
			args: args{
				obj: uint32(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "uint64 > " + uint64To,
			args: args{
				obj: uint64(77),
			},
			wantOut: uint64(77),
		},
		{
			name: "string > " + uint64To,
			args: args{
				obj: "77",
			},
			wantOut: uint64(77),
		},
	}
	for _, tt := range uint64Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uint64](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
	const stringTo = "string"
	stringTests := []testCase[string]{
		{
			name: "nil > " + stringTo,
			args: args{
				obj: nil,
			},
			wantOut: "<nil>",
		},
		{
			name: "complex64 > " + stringTo,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantOut: "(77+0i)",
		},
		{
			name: "complex128 > " + stringTo,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantOut: "(77+0i)",
		},
		{
			name: "float32 > " + stringTo,
			args: args{
				obj: float32(77),
			},
			wantOut: "77",
		},
		{
			name: "float64 > " + stringTo,
			args: args{
				obj: float64(77),
			},
			wantOut: "77",
		},
		{
			name: "float32 dot0 > " + stringTo,
			args: args{
				obj: float32(77.0),
			},
			wantOut: "77",
		},
		{
			name: "float64 dot0 > " + stringTo,
			args: args{
				obj: float64(77.0),
			},
			wantOut: "77",
		},
		{
			name: "float32 dot1 > " + stringTo,
			args: args{
				obj: float32(77.1),
			},
			wantOut: "77.1",
		},
		{
			name: "float64 dot1 > " + stringTo,
			args: args{
				obj: float64(77.1),
			},
			wantOut: "77.1",
		},
		{
			name: "int > " + stringTo,
			args: args{
				obj: 77,
			},
			wantOut: "77",
		},
		{
			name: "int8 > " + stringTo,
			args: args{
				obj: int8(77),
			},
			wantOut: "77",
		},
		{
			name: "int16 > " + stringTo,
			args: args{
				obj: int16(77),
			},
			wantOut: "77",
		},
		{
			name: "int32 > " + stringTo,
			args: args{
				obj: int32(77),
			},
			wantOut: "77",
		},
		{
			name: "int64 > " + stringTo,
			args: args{
				obj: int64(77),
			},
			wantOut: "77",
		},
		{
			name: "uint > " + stringTo,
			args: args{
				obj: uint(77),
			},
			wantOut: "77",
		},
		{
			name: "uint8 > " + stringTo,
			args: args{
				obj: uint8(77),
			},
			wantOut: "77",
		},
		{
			name: "uint16 > " + stringTo,
			args: args{
				obj: uint16(77),
			},
			wantOut: "77",
		},
		{
			name: "uint32 > " + stringTo,
			args: args{
				obj: uint32(77),
			},
			wantOut: "77",
		},
		{
			name: "uint64 > " + stringTo,
			args: args{
				obj: uint64(77),
			},
			wantOut: "77",
		},
		{
			name: "string > " + stringTo,
			args: args{
				obj: "77",
			},
			wantOut: "77",
		},
		{
			name: "error > " + stringTo,
			args: args{
				obj: errors.New("my error message"),
			},
			wantOut: "my error message",
		},
		{
			name: "array > " + stringTo,
			args: args{
				obj: [5]any{777, "777", true, 7.7, nil},
			},
			wantOut: `[777,"777",true,7.7,null]`,
		},
		{
			name: "map > " + stringTo,
			args: args{
				obj: jsonFullCheckMapObj,
			},
			wantOut: jsonFullCheckMapStr,
		},
		{
			name: "slice > " + stringTo,
			args: args{
				obj: []any{777, "777", true, 7.7, nil},
			},
			wantOut: `[777,"777",true,7.7,null]`,
		},
		{
			name: "func > " + stringTo,
			args: args{
				obj: func1,
			},
			wantOut: `func func1() {
	fmt.Sprintln("func1 body")
}`,
		},
		{
			name: "func interface > " + stringTo,
			args: args{
				obj: Shape.area,
			},
			wantOut: `func(conv.Shape) float64`,
		},
		{
			name: "func interface pointer > " + stringTo,
			args: args{
				obj: Shape.area,
			},
			wantOut: `func(conv.Shape) float64`,
		},
		{
			name: "func struct > " + stringTo,
			args: args{
				obj: Rect.area,
			},
			wantOut: `func (r Rect) area() float64 {
	return r.w * r.h
}`,
		},
		{
			name: "func struct pointer > " + stringTo,
			args: args{
				obj: (*Rect).area2,
			},
			wantOut: `func (r *Rect) area2() float64 {
	return r.w * r.h
}`,
		},
		{
			name: "func struct pointer outside interface > " + stringTo,
			args: args{
				obj: (*Rect).area3,
			},
			wantOut: `func (r *Rect) area3() float64 {
	return r.w * r.h
}`,
		},
		{
			name: "func struct instance > " + stringTo,
			args: args{
				obj: Rect{}.area,
			},
			wantOut: `func() float64`,
		},
		{
			name: "func struct instance pointer > " + stringTo,
			args: args{
				obj: Rect{}.area,
			},
			wantOut: `func() float64`,
		},
		{
			name: "func struct instance pointer outside interface > " + stringTo,
			args: args{
				obj: Rect{}.area,
			},
			wantOut: `func() float64`,
		},
		{
			name: "anonymous function > " + stringTo,
			args: args{
				obj: func(s string) error { return nil },
			},
			wantOut: `func(string) error`,
		},
		{
			name: "anonymous function2 > " + stringTo,
			args: args{
				obj: func(string) error { return nil },
			},
			wantOut: `func(string) error`,
		},
		{
			name: "anonymous function3 > " + stringTo,
			args: args{
				obj: func(x string, y ...string) error { return nil },
			},
			wantOut: `func(string, ...string) error`,
		},
		{
			name: "anonymous function3 > " + stringTo,
			args: args{
				obj: func(x string, y ...string) error { return nil },
			},
			wantOut: `func(string, ...string) error`,
		},
		{
			name: "struct instance > " + stringTo,
			args: args{
				obj: Rect{X: 777, Y: 888},
			},
			wantOut: `{"X":777,"Y":888}`,
		},
		{
			name: "struct instance pointer > " + stringTo,
			args: args{
				obj: &[]Rect{{X: 777, Y: 888}}[0],
			},
			wantOut: `{"X":777,"Y":888}`,
		},
		{
			name: "struct instance slice > " + stringTo,
			args: args{
				obj: []Rect{{X: 777, Y: 888}, {X: 111, Y: 222}},
			},
			wantOut: `[{"X":777,"Y":888},{"X":111,"Y":222}]`,
		},
		{
			name: "struct instance arr pointer > " + stringTo,
			args: args{
				obj: &[][]Rect{{{X: 777, Y: 888}, {X: 111, Y: 222}}}[0],
			},
			wantOut: `[{"X":777,"Y":888},{"X":111,"Y":222}]`,
		},
		{
			name: "channel pointer > " + stringTo,
			args: args{
				obj: make(chan int),
			},
			wantReg: regexp.MustCompile("^0x[0-9a-f]+$"),
		},
		{
			name: "byte slice > " + stringTo, //Golang json marshal will base64, we handle this as number array
			args: args{
				obj: []byte("byte to string"),
			},
			wantOut: `[98,121,116,101,32,116,111,32,115,116,114,105,110,103]`,
		},
		{
			name: "byte array > " + stringTo,
			args: args{
				obj: [14]byte{98, 121, 116, 101, 32, 116, 111, 32, 115, 116, 114, 105, 110, 103}, //Golang json marshal will convert as number array
			},
			wantOut: `[98,121,116,101,32,116,111,32,115,116,114,105,110,103]`,
		},
		{
			name: "string slice > " + stringTo, //Golang json marshal will base64, we handle this as number array
			args: args{
				obj: []string{"a", "b", "c"},
			},
			wantOut: `["a","b","c"]`,
		},
		{
			name: "string array > " + stringTo,
			args: args{
				obj: [3]string{"a", "b", "c"},
			},
			wantOut: `["a","b","c"]`,
		},
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[string](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantReg != nil {
				if !tt.wantReg.Match([]byte(gotOut)) {
					t.Errorf("To() gotOut = %v %v, wantReg %v", gotOut, reflect.TypeOf(gotOut), tt.wantReg)
				}
			} else if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const otherTo = "other"
	otherTests := []testCase[other]{
		{
			name: "uint64 > " + otherTo,
			args: args{
				obj: uint64(77),
			},
			wantOut: "",
			wantErr: true,
		},
		{
			name: "string > " + otherTo,
			args: args{
				obj: "77",
			},
			wantOut: "",
			wantErr: true,
		},
	}
	for _, tt := range otherTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[other](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const errorTo = "error"
	errorTests := []testCase[error]{
		{
			name: "string > " + errorTo,
			args: args{
				obj: "error message",
			},
			wantOut: errors.New("error message"),
		},
		{
			name: "complex64 > " + errorTo,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantOut: errors.New("(77+0i)"),
		},
	}
	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[error](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut.Error() != tt.wantOut.Error() {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const uintptrTo = "uintptr"
	uintptrTests := []testCase[uintptr]{
		{
			name: "int > " + uintptrTo,
			args: args{
				obj: 77,
			},
			wantOut: 0,
			wantErr: true,
		},
	}
	for _, tt := range uintptrTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[uintptr](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const structTo = "struct"
	structTests := []testCase[Rect]{
		{
			name: "string > " + structTo,
			args: args{
				obj: `{"X":777,"Y":888}`,
			},
			wantOut: Rect{X: 777, Y: 888},
		},
		{
			name: "string null > " + structTo,
			args: args{
				obj: `null`,
			},
			wantOut: Rect{},
		},
		{
			name: "string nil > " + structTo,
			args: args{
				obj: `<nil>`,
			},
			wantOut: Rect{},
		},
		{
			name: "string empty > " + structTo,
			args: args{
				obj: ``,
			},
			wantOut: Rect{},
		},
		{
			name: "string empty obj > " + structTo,
			args: args{
				obj: `{}`,
			},
			wantOut: Rect{},
		},
	}
	for _, tt := range structTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[Rect](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() struct = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const structPointerTo = "structPointer"
	structPointerTests := []testCase[*Rect]{
		{
			name: "string > " + structPointerTo,
			args: args{
				obj: `{"X":777,"Y":888}`,
			},
			wantOut: &Rect{X: 777, Y: 888},
		},
		{
			name: "string null > " + structPointerTo,
			args: args{
				obj: `null`,
			},
			wantOut: nil,
		},
		{
			name: "string nil > " + structPointerTo,
			args: args{
				obj: `<nil>`,
			},
			wantOut: nil,
		},
		{
			name: "string empty > " + structPointerTo,
			args: args{
				obj: ``,
			},
			wantOut: nil,
		},
		{
			name: "string empty obj > " + structPointerTo,
			args: args{
				obj: `{}`,
			},
			wantOut: &Rect{},
		},
	}
	for _, tt := range structPointerTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[*Rect](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() structPointer = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const mapTo = "map"
	mapTests := []testCase[map[string]any]{
		{
			name: "string > " + mapTo,
			args: args{
				obj: jsonFullCheckMapStr,
			},
			wantOut: jsonFullCheckMapObjOnlyFloat,
		},
		{
			name: "string null > " + mapTo,
			args: args{
				obj: `null`,
			},
			wantOut: nil,
		},
		{
			name: "string nil > " + mapTo,
			args: args{
				obj: `<nil>`,
			},
			wantOut: nil,
		},
		{
			name: "string empty > " + mapTo,
			args: args{
				obj: ``,
			},
			wantOut: nil,
		},
		{
			name: "string empty obj > " + mapTo,
			args: args{
				obj: `{}`,
			},
			wantOut: map[string]any{},
		},
	}
	for _, tt := range mapTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[map[string]any](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() map[string]any = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const anyTo = "map"
	anyTests := []testCase[any]{
		{
			name: "string > " + anyTo,
			args: args{
				obj: jsonFullCheckMapStr,
			},
			wantOut: jsonFullCheckMapStr,
		},
		{
			name: "string null > " + anyTo,
			args: args{
				obj: `null`,
			},
			wantOut: `null`,
		},
		{
			name: "string nil > " + anyTo,
			args: args{
				obj: `<nil>`,
			},
			wantOut: `<nil>`,
		},
		{
			name: "string empty > " + anyTo,
			args: args{
				obj: ``,
			},
			wantOut: ``,
		},
		{
			name: "string empty obj > " + anyTo,
			args: args{
				obj: `{}`,
			},
			wantOut: `{}`,
		},
		{
			name: "nil obj > " + anyTo,
			args: args{
				obj: nil,
			},
			wantOut: nil,
		},
	}
	for _, tt := range anyTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[any](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() any = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
}

func TestToForce(t *testing.T) {
	t.Parallel()
	type args struct {
		obj any
	}
	type testCase[T any] struct {
		name    string
		args    args
		wantOut T
	}

	const float32To = "float32"
	float32Tests := []testCase[float32]{
		{
			name: "complex64 > " + float32To,
			args: args{
				obj: complex(float32(77), float32(0)),
			},
			wantOut: 0,
		},
		{
			name: "complex128 > " + float32To,
			args: args{
				obj: complex(float64(77), float64(0)),
			},
			wantOut: 0,
		},
		{
			name: "float32 > " + float32To,
			args: args{
				obj: float32(77),
			},
			wantOut: 77,
		},
	}
	for _, tt := range float32Tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := ToForce[float32](tt.args.obj)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}

	const otherTo = "other"
	otherTests := []testCase[other]{
		{
			name: "uint64 > " + otherTo,
			args: args{
				obj: uint64(77),
			},
			wantOut: "",
		},
		{
			name: "string > " + otherTo,
			args: args{
				obj: "77",
			},
			wantOut: "",
		},
	}
	for _, tt := range otherTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := ToForce[other](tt.args.obj)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("To() gotOut = %v %v, want %v %v", gotOut, reflect.TypeOf(gotOut), tt.wantOut, reflect.TypeOf(tt.wantOut))
			}
		})
	}
}

func Test_getFuncAST(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	type args struct {
		funcname string
		filename string
	}
	tests := []struct {
		name  string
		args  args
		want  *ast.FuncDecl
		want1 *token.FileSet
	}{
		{
			name: "error test",
			args: args{
				funcname: "funcname",
				filename: "xyz",
			},
			want:  nil,
			want1: nil,
		},
		{
			name: "match func but not match struct",
			args: args{
				funcname: "Triangle.area2",
				filename: filename,
			},
			want:  nil,
			want1: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getFuncAST(tt.args.funcname, tt.args.filename)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFuncAST() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getFuncAST() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getFuncBodyString(t *testing.T) {
	t.Parallel()
	_, filename, _, _ := runtime.Caller(0)
	_, fs := getFuncAST("Test_getFuncBodyString", filename)
	type args struct {
		f  any
		fs *token.FileSet
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "error test",
			args: args{
				f:  "improper node type",
				fs: fs,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got1 := getFuncBodyString(tt.args.f, tt.args.fs); !reflect.DeepEqual(got1, tt.want) {
				t.Errorf("getFuncBodyString() got = %v, want %v", got1, tt.want)
			}
		})
	}
}

func toMap(str string) (out map[string]any) {
	json.Unmarshal([]byte(str), &out)
	return
}

func TestGetItemsNeedSortAsMapKeyUndeterministic(t *testing.T) {
	type args struct {
		obj  any
		keys string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "#k",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a","k2":2,"k3":3},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#.#k",
			},
			want: []string{"k2", "k3", "value", "value", "value", "value", "value", "valuex"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetItems(tt.args.obj, tt.args.keys)
			ss := []string{}
			for _, v := range got {
				ss = append(ss, v.(string))
			}
			sort.Strings(ss)
			if !reflect.DeepEqual(ss, tt.want) {

				t.Errorf("GetItems() = %v, want %v", ss, tt.want)
			}
		})
	}
}

func TestGetItems(t *testing.T) {
	var va any
	var vc64 complex64
	var vc128 complex128
	var vf32 float32 = 32.32
	var vf64 float64 = -64.64
	var vi1 int = 1
	var vi8 int8 = -8
	var vi16 int16 = 16
	var vi32 int32 = -32
	var vi64 int64 = 64
	var vui1 uint = 1
	var vui8 uint8 = 8
	var vui16 uint16 = 16
	var vui32 uint32 = 32
	var vui64 uint64 = 64
	var vs0 string
	var vs string = "string"
	var vuip uintptr

	type args struct {
		obj  any
		keys string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "#k slice",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a","k2":2,"k3":3},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#k",
			},
			want: []any{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "string&slice of string, unmarshal always use float not int",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#.value",
			},
			want: []any{"a", "b", nil, nil, nil, "777", 7.0},
		},
		{
			name: "string&slice of string, unmarshal always use float not int. alternative slice #v",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#v.value",
			},
			want: []any{"a", "b", nil, nil, nil, "777", 7.0},
		},
		{
			name: "string&slice of string, unmarshal always use float not int. map #",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#.#",
			},
			want: []any{"a", "b", nil, "b", "777", 7.0},
		},
		{
			name: "string/slice of string, unmarshal always use float not int. alternative map #v",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#.#v",
			},
			want: []any{"a", "b", nil, "b", "777", 7.0},
		},
		{
			name: "string&slice of string. map # order",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":"x"},{"valuex":"ax","valueb":"ab","valuea":"aa","value":"a"},{},{"value":"777"},{"value":"7"}]}}}`),
				keys: "path1.path2.items.#.#",
			},
			want: []any{"a", "b", "x", "a", "aa", "ab", "ax", "777", "7"},
		},
		{
			name: "string&slice of string. regexp",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":"x"},{"valuex":"ax","valueb":"ab","valuea":"aa","zz":"zz","value":"a"},{},{"value":"777"},{"value":"7"}]}}}`),
				keys: "path1.path2.items.#.^va",
			},
			want: []any{"a", "b", "x", "a", "aa", "ab", "ax", "777", "7"},
		},
		{
			name: "string&slice of string. regexp2",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":"x"},{"valuex":"ax","valueb":"ab","valuea":"aa","zz":"zz","value":"a"},{},{"value":"777"},{"value":"7"}]}}}`),
				keys: "path1.path2.items.#.^" + DotAlternative + "*x$",
			},
			want: []any{"ax"},
		},
		{
			name: "string&slice of string. regexp err",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":"x"},{"valuex":"ax","valueb":"ab","valuea":"aa","zz":"zz","value":"a"},{},{"value":"777"},{"value":"7"}]}}}`),
				keys: "path1.path2.items.#.^(xx" + DotAlternative + "*x$",
			},
			want: []any{},
		},
		{
			name: "map of all value types by any",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": va,
								},
								map[string]any{
									"value": nil,
								},
								map[string]any{},
								map[string]any{
									"value": true,
								},
								map[string]any{
									"value": false,
								},
								map[string]any{
									"value": vc64,
								},
								map[string]any{
									"value": vc128,
								},
								map[string]any{
									"value": vf32,
								},
								map[string]any{
									"value": vf64,
								},
								map[string]any{
									"value": vi1,
								},
								map[string]any{
									"value": vi8,
								},
								map[string]any{
									"value": vi16,
								},
								map[string]any{
									"value": vi32,
								},
								map[string]any{
									"value": vi64,
								},
								map[string]any{
									"value": vui1,
								},
								map[string]any{
									"value": vui8,
								},
								map[string]any{
									"value": vui16,
								},
								map[string]any{
									"value": vui32,
								},
								map[string]any{
									"value": vui64,
								},
								map[string]any{
									"value": vs0,
								},
								map[string]any{
									"value": vs,
								},
								map[string]any{
									"value": vuip,
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value",
			},
			want: []any{va, nil, nil, true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		},
		{
			name: "slice of all value types",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": []any{va, va},
								},
								map[string]any{
									"value": []any{nil, nil},
								},
								map[string]any{
									"value": []any{},
								},
								map[string]any{
									"value": []bool{true, true},
								},
								map[string]any{
									"value": []bool{false, false},
								},
								map[string]any{
									"value": []complex64{vc64, vc64},
								},
								map[string]any{
									"value": []complex128{vc128, vc128},
								},
								map[string]any{
									"value": []float32{vf32, vf32},
								},
								map[string]any{
									"value": []float64{vf64, vf64},
								},
								map[string]any{
									"value": []int{vi1, vi1},
								},
								map[string]any{
									"value": []int8{vi8, vi8},
								},
								map[string]any{
									"value": []int16{vi16, vi16},
								},
								map[string]any{
									"value": []int32{vi32, vi32},
								},
								map[string]any{
									"value": []int64{vi64, vi64},
								},
								map[string]any{
									"value": []uint{vui1, vui1},
								},
								map[string]any{
									"value": []uint8{vui8, vui8},
								},
								map[string]any{
									"value": []uint16{vui16, vui16},
								},
								map[string]any{
									"value": []uint32{vui32, vui32},
								},
								map[string]any{
									"value": []uint64{vui64, vui64},
								},
								map[string]any{
									"value": []string{vs0, vs0},
								},
								map[string]any{
									"value": []string{vs, vs},
								},
								map[string]any{
									"value": []uintptr{vuip, vuip},
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value.1",
			},
			want: []any{va, nil, nil, true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		},
		{
			name: "map of all value types by native",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": va,
								},
								map[string]any{
									"value": nil,
								},
								map[string]int{},
								map[string]bool{
									"value": true,
								},
								map[string]bool{
									"value": false,
								},
								map[string]complex64{
									"value": vc64,
								},
								map[string]complex128{
									"value": vc128,
								},
								map[string]float32{
									"value": vf32,
								},
								map[string]float64{
									"value": vf64,
								},
								map[string]int{
									"value": vi1,
								},
								map[string]int8{
									"value": vi8,
								},
								map[string]int16{
									"value": vi16,
								},
								map[string]int32{
									"value": vi32,
								},
								map[string]int64{
									"value": vi64,
								},
								map[string]uint{
									"value": vui1,
								},
								map[string]uint8{
									"value": vui8,
								},
								map[string]uint16{
									"value": vui16,
								},
								map[string]uint32{
									"value": vui32,
								},
								map[string]uint64{
									"value": vui64,
								},
								map[string]string{
									"value": vs0,
								},
								map[string]string{
									"value": vs,
								},
								map[string]uintptr{
									"value": vuip,
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value",
			},
			want: []any{va, nil, nil, true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetItems(tt.args.obj, tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetItemsWithOptOmitNoValue(t *testing.T) {
	var va any
	var vc64 complex64
	var vc128 complex128
	var vf32 float32 = 32.32
	var vf64 float64 = -64.64
	var vi1 int = 1
	var vi8 int8 = -8
	var vi16 int16 = 16
	var vi32 int32 = -32
	var vi64 int64 = 64
	var vui1 uint = 1
	var vui8 uint8 = 8
	var vui16 uint16 = 16
	var vui32 uint32 = 32
	var vui64 uint64 = 64
	var vs0 string
	var vs string = "string"
	var vuip uintptr

	type args struct {
		obj  any
		keys string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "string/slice of string, unmarshal always use float not int",
			args: args{
				obj:  toMap(`{"path1":{"path2":{"items":[{"value":"a"},{"value":"b"},{"value":null},{"valuex":"b"},{},{"value":"777"},{"value":7}]}}}`),
				keys: "path1.path2.items.#.value",
			},
			want: []any{"a", "b", "777", 7.0},
		}, {
			name: "map of all value types by any",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": va,
								},
								map[string]any{
									"value": nil,
								},
								map[string]any{},
								map[string]any{
									"value": true,
								},
								map[string]any{
									"value": false,
								},
								map[string]any{
									"value": vc64,
								},
								map[string]any{
									"value": vc128,
								},
								map[string]any{
									"value": vf32,
								},
								map[string]any{
									"value": vf64,
								},
								map[string]any{
									"value": vi1,
								},
								map[string]any{
									"value": vi8,
								},
								map[string]any{
									"value": vi16,
								},
								map[string]any{
									"value": vi32,
								},
								map[string]any{
									"value": vi64,
								},
								map[string]any{
									"value": vui1,
								},
								map[string]any{
									"value": vui8,
								},
								map[string]any{
									"value": vui16,
								},
								map[string]any{
									"value": vui32,
								},
								map[string]any{
									"value": vui64,
								},
								map[string]any{
									"value": vs0,
								},
								map[string]any{
									"value": vs,
								},
								map[string]any{
									"value": vuip,
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value",
			},
			want: []any{true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		}, {
			name: "slice of all value types",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": []any{va, va},
								},
								map[string]any{
									"value": []any{nil, nil},
								},
								map[string]any{
									"value": []any{},
								},
								map[string]any{
									"value": []bool{true, true},
								},
								map[string]any{
									"value": []bool{false, false},
								},
								map[string]any{
									"value": []complex64{vc64, vc64},
								},
								map[string]any{
									"value": []complex128{vc128, vc128},
								},
								map[string]any{
									"value": []float32{vf32, vf32},
								},
								map[string]any{
									"value": []float64{vf64, vf64},
								},
								map[string]any{
									"value": []int{vi1, vi1},
								},
								map[string]any{
									"value": []int8{vi8, vi8},
								},
								map[string]any{
									"value": []int16{vi16, vi16},
								},
								map[string]any{
									"value": []int32{vi32, vi32},
								},
								map[string]any{
									"value": []int64{vi64, vi64},
								},
								map[string]any{
									"value": []uint{vui1, vui1},
								},
								map[string]any{
									"value": []uint8{vui8, vui8},
								},
								map[string]any{
									"value": []uint16{vui16, vui16},
								},
								map[string]any{
									"value": []uint32{vui32, vui32},
								},
								map[string]any{
									"value": []uint64{vui64, vui64},
								},
								map[string]any{
									"value": []string{vs0, vs0},
								},
								map[string]any{
									"value": []string{vs, vs},
								},
								map[string]any{
									"value": []uintptr{vuip, vuip},
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value.1",
			},
			want: []any{true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		}, {
			name: "map of all value types by native",
			args: args{
				obj: map[string]any{
					"path1": map[string]any{
						"path2": map[string]any{
							"items": []any{
								map[string]any{
									"value": va,
								},
								map[string]any{
									"value": nil,
								},
								map[string]int{},
								map[string]bool{
									"value": true,
								},
								map[string]bool{
									"value": false,
								},
								map[string]complex64{
									"value": vc64,
								},
								map[string]complex128{
									"value": vc128,
								},
								map[string]float32{
									"value": vf32,
								},
								map[string]float64{
									"value": vf64,
								},
								map[string]int{
									"value": vi1,
								},
								map[string]int8{
									"value": vi8,
								},
								map[string]int16{
									"value": vi16,
								},
								map[string]int32{
									"value": vi32,
								},
								map[string]int64{
									"value": vi64,
								},
								map[string]uint{
									"value": vui1,
								},
								map[string]uint8{
									"value": vui8,
								},
								map[string]uint16{
									"value": vui16,
								},
								map[string]uint32{
									"value": vui32,
								},
								map[string]uint64{
									"value": vui64,
								},
								map[string]string{
									"value": vs0,
								},
								map[string]string{
									"value": vs,
								},
								map[string]uintptr{
									"value": vuip,
								},
							},
						},
					},
				},
				keys: "path1.path2.items.#.value",
			},
			want: []any{true, false, vc64, vc128, vf32, vf64, vi1, vi8, vi16, vi32, vi64, vui1, vui8, vui16, vui32, vui64, vs0, vs, vuip},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetItems(tt.args.obj, tt.args.keys, OptOmitNoValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetItems() = %v, want %v", got, tt.want)
			}
		})
	}
}
