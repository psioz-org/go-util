package zz

import (
	"reflect"
	"testing"
)

type other string

func TestGetItemTestGeneric(t *testing.T) {
	type args struct {
		obj  interface{}
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{
						"y": map[string]interface{}{
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{
						"y": []interface{}{
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{
						"y": map[string]interface{}{
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
	type args struct {
		obj  interface{}
		keys string
		sep  string
	}
	type result struct{}
	result1 := result{}
	result2 := result{}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "test result success",
			args: args{
				obj: map[string]interface{}{
					"x": map[string]interface{}{
						"y": map[string]interface{}{
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{
						"y": []interface{}{
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{},
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
				obj: map[string]interface{}{
					"x": map[string]interface{}{},
				},
				keys: "x/y/0",
				sep:  "/",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetItem[interface{}](tt.args.obj, tt.args.keys, tt.args.sep)
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

func TestGetItem1(t *testing.T) {
	//TODO somehow wantErr is true
	//type args struct {
	//	obj  interface{}
	//	keys string
	//	sep  string
	//}
	//type testCase[T any] struct {
	//	name    string
	//	args    args
	//	want    T
	//	wantErr bool
	//}
	//
	//tests := []testCase[int]{
	//	{
	//		name: "case 1",
	//		args: args{
	//			obj:  ToMap(`{"x":{"y":[false,true,{"z":777}]}}`),
	//			keys: "x.y.2.z",
	//			sep:  ".",
	//		},
	//		want: 777,
	//	},
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		got, err := GetItem[int](tt.args.obj, tt.args.keys, tt.args.sep)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("GetItem() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("GetItem() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}

func TestTernary(t *testing.T) {
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
	type args struct {
		obj interface{}
	}
	type testCase[T any] struct {
		name    string
		args    args
		wantOut T
		wantErr bool
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
	}
	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := To[string](tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("To() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
}

func TestToForce(t *testing.T) {
	type args struct {
		obj interface{}
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
